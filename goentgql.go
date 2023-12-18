package goentgql

import (
	"context"
	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/schema"
	"fmt"
	"github.com/99designs/gqlgen/graphql"
	"github.com/dmalykh/goentgql/config"
	"github.com/dmalykh/goentgql/entgen"
	"github.com/dmalykh/goentgql/gqlgen"
	"github.com/dmalykh/goentgql/gqlgen/middleware"
	"github.com/dmalykh/goentgql/runner"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/urfave/cli/v2"
	_ "github.com/xiaoqidun/entps"
	"os"
)

type GoEntGQL interface {
	Execute(ctx context.Context) error
}

func New(options ...Option) GoEntGQL {
	var s = app{
		entSchemaPath: `/schema`,
		migrator: func(drv *sql.Driver) Migrator {
			return &dummyMigrator{}
		},
	}
	for _, option := range options {
		option(&s)
	}
	return &s
}

type app struct {
	entSchemaPath    string
	packageName      string
	executableSchema ExecutableSchema
	extensions       []graphql.HandlerExtension
	driver           *sql.Driver
	migrator         MigrationRunner
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
			cfg, err := config.NewGenerate(&config.ConfigGenerate{
				BasePath: c.Args().Get(0),
				Package:  s.packageName,
				Dir:      `/generated`,
				EntConfig: &config.EntConfig{
					SchemaPath: s.entSchemaPath,
				},
			})

			if err != nil {
				return fmt.Errorf(`generate: config error: %w`, err)
			}

			{
				// Generate ent
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
				if err := gqlgen.New(cfg.GraphQLConfig()).Generate(); err != nil {
					return fmt.Errorf(`error generate gqlgen: %w`, err)
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
		},
		Action: func(c *cli.Context) error {

			drv, err := sql.Open(c.String(`driver`), c.String(`dsn`))
			if err != nil {
				return fmt.Errorf(`can't connect to database: %w`, err)
			}

			if err := s.migrator(drv).Create(c.Context, schema.WithGlobalUniqueID(true)); err != nil {
				return fmt.Errorf(`can't migrate database: %w`, err)
			}

			fmt.Printf("Listen %s...\n", c.String(`addr`))
			return runner.Run(s.executableSchema(drv), &runner.Config{
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
