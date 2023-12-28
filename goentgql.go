package goentgql

import (
	"context"
	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/schema"
	"fmt"
	"github.com/dmalykh/goentgql/generator/config"
	"github.com/dmalykh/goentgql/generator/entgen"
	"github.com/dmalykh/goentgql/generator/gqlgen"
	"github.com/dmalykh/goentgql/generator/gqlgen/middleware"
	"github.com/dmalykh/goentgql/generator/service"
	"github.com/dmalykh/goentgql/runner"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
	_ "github.com/xiaoqidun/entps"
	"os"
)

type GoEntGQL interface {
	Execute(ctx context.Context) error
}

func New(options ...Option) GoEntGQL {
	var s = app{
		schemaDir: `/schema`,
	}

	for _, option := range options {
		option(&s)
	}

	return &s
}

type app struct {
	schemaDir   string
	packageName string
	service     ServiceRunner
}

func (s *app) Execute(ctx context.Context) error {
	app := &cli.App{
		Commands: []*cli.Command{
			s.generateCmd(),
			s.runCmd(),
		},
	}

	if err := app.Run(os.Args); err != nil {
		return fmt.Errorf(`execution error: %w`, err)
	}

	return nil
}

func (s *app) generateCmd() *cli.Command {
	return &cli.Command{
		Name:  "generate",
		Usage: "generates ent&gqlgen code",
		Action: func(c *cli.Context) error {
			// Create Configuration
			cfg, err := config.NewGenerator(&config.GeneratorConfig{
				BasePath:  c.Args().Get(0),
				Package:   s.packageName,
				OutputDir: `/generated`,
				EntConfig: &config.EntConfig{
					SchemaPath: s.schemaDir,
				},
			})

			if err != nil {
				return fmt.Errorf(`generate: config error: %w`, err)
			}

			{
				// Generate ent
				log.Info().Interface(`end config`, cfg.EntConfig()).
					Interface(`gqlgen config`, cfg.GraphQLConfig()).Msg(`Generating ent files...`)

				gen, err := entgen.New(cfg.EntConfig(), cfg.GraphQLConfig())
				if err != nil {
					return fmt.Errorf(`init entgen error: %w`, err)
				}

				if err := gen.Generate(); err != nil {
					return fmt.Errorf(`error generate ent: %w`, err)
				}
			}

			{
				// Generate gqlgen
				log.Info().Interface(`gqlgen config`, cfg.GraphQLConfig()).Msg(`Generating gqlgen files...`)

				if err := gqlgen.New(cfg.GraphQLConfig()).Generate(); err != nil {
					return fmt.Errorf(`error generate gqlgen: %w`, err)
				}
			}

			{
				// Generate service
				log.Info().Interface(`service config`, cfg.ServiceConfig()).Msg(`Generating service files...`)

				if err := service.New(cfg.ServiceConfig()).Generate(); err != nil {
					return fmt.Errorf(`error generate service: %w`, err)
				}
			}
			return nil
		},
	}
}

func (s *app) runCmd() *cli.Command {
	return &cli.Command{
		Name:  "run",
		Usage: "run the service",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    `verbose`,
				EnvVars: []string{`VERBOSE`},
			},
			&cli.StringFlag{
				Name:    `addr`,
				Value:   `:8080`,
				EnvVars: []string{`ADDR`},
			},
			&cli.StringFlag{
				Name:    `driver`,
				Usage:   `driver for ent`,
				EnvVars: []string{`DRIVER`},
			},
			&cli.StringFlag{
				Name:    `dsn`,
				Usage:   `dsn for ent`,
				EnvVars: []string{`DSN`},
			},
			&cli.BoolFlag{
				Name:    `skip-migrations`,
				Usage:   `skip migrations`,
				EnvVars: []string{`SKIP_MIGRATIONS`},
			},
		},
		Action: func(c *cli.Context) error {

			drv, err := sql.Open(c.String(`driver`), c.String(`dsn`))
			if err != nil {
				return fmt.Errorf(`can't connect to database: %w`, err)
			}

			svc := s.service(drv)

			if !c.Bool(`skip-migrations`) {
				if err := svc.MigrateSchema(c.Context, schema.WithGlobalUniqueID(true)); err != nil {
					return fmt.Errorf(`can't migrate database: %w`, err)
				}
			}

			fmt.Printf("Listen %s...\n", c.String(`addr`))
			return runner.Run(svc.ExecutionSchema(), &runner.Config{
				Addr:    c.String(`addr`),
				Name:    s.packageName,
				Verbose: c.Bool(`verbose`),
			}, middleware.Metrics())
		},
	}
}

//
//func Run() {
//
//
//	client, err := ent.Open(
//		cli.Driver, // i.e. sqlite3
//		cli.DSN,    //i.e. file:ent?mode=memory&cache=shared&_fk=1
//		ent.Log(func(i ...interface{}) {
//			log.Debug().Fields(i)
//		}),
//	)
//
//	// Enable verbose mode
//	if cli.Verbose {
//		client = client.Debug()
//	}
//
//	// Create schema
//	if !cli.NoMigration {

//	}
//
//	// Metrics
//	var totalRequestMetric = promauto.NewCounter(prometheus.CounterOpts{
//		Name: "request_total",
//		Help: "The total number of request",
//	})
//	var operationDurationMetric = promauto.NewHistogram(prometheus.HistogramOpts{
//		Name:    "http_server_request_duration_seconds",
//		Help:    "Histogram of response time for handler in seconds",
//		Buckets: []float64{.000001, .00001, .0001, .001, .01, .025, .05, .1, .25, .5, 1},
//	})
//
//	var prometheusMiddleware = func(ctx context.Context, next gqlgen.OperationHandler) gqlgen.ResponseHandler {
//		totalRequestMetric.Inc()
//		var start = time.Now()
//		var resp = next(ctx)
//		var t = time.Since(start).Seconds()
//		operationDurationMetric.Observe(t)
//		return resp
//	}
//
//	var errChan = make(chan error)
//
//	// Listener for metrics
//	go func() {
//		if err := http.ListenAndServe(fmt.Sprintf(`:%s`, cli.PrometheusAddr), promhttp.Handler()); err != nil {
//			errChan <- err
//		}
//	}()
//
//	// Server
//	go func() {
//		if err := graphql.Run(client, fmt.Sprintf(`:%s`, cli.Addr), cli.Verbose, prometheusMiddleware); err != nil {
//			errChan <- err
//		}
//	}()
//
//}
