package goentgql

import (
	"context"
	"database/sql"
	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/schema"
	"fmt"
	gqlgenconfig "github.com/99designs/gqlgen/codegen/config"
	"github.com/dmalykh/goentgql/generator/config"
	"github.com/dmalykh/goentgql/generator/entgen"
	"github.com/dmalykh/goentgql/generator/gqlgen"
	"github.com/dmalykh/goentgql/generator/gqlgen/middleware"
	"github.com/dmalykh/goentgql/generator/service"
	"github.com/dmalykh/goentgql/runner"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
	_ "github.com/xiaoqidun/entps"
	"os"
)

type GoEntGQL interface {
	Execute(ctx context.Context) error
}

func New(options ...Option) GoEntGQL {
	var s = App{
		schemaDir: `/schema`,
	}

	for _, option := range options {
		option(&s)
	}

	return &s
}

type App struct {
	schemaDir     string
	packageName   string
	service       ServiceRunner
	extensions    []Extension
	configOptions struct {
		IDType             []string
		WithGlobalUniqueID bool
	}
}

func (s *App) Execute(ctx context.Context) error {
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

func (s *App) generateCmd() *cli.Command {
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
				GraphQLConfig: &config.GraphQLConfig{
					Models: gqlgenconfig.TypeMap{
						`ID`: gqlgenconfig.TypeMapEntry{
							Model: s.configOptions.IDType,
						},
					},
				},
			})

			if err != nil {
				return fmt.Errorf(`generate: config error: %w`, err)
			}

			for _, ext := range s.extensions {
				if genext, ok := ext.(GeneratorExtension); ok {
					if err := genext.Generator(c, cfg); err != nil {
						return fmt.Errorf(`can't execute extension: %w`, err)
					}
				}
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

			log.Info().Msg(`Done`)

			return nil
		},
	}
}

func (s *App) runCmd() *cli.Command {
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
				Usage:   `driver for a database`,
				EnvVars: []string{`DRIVER`},
				Value:   `sqlite3`,
			},
			&cli.StringFlag{
				Name:    `dialect`,
				Usage:   `dialect for ent`,
				EnvVars: []string{`DIALECT`},
				Value:   dialect.SQLite,
			},
			&cli.StringFlag{
				Name:    `dsn`,
				Usage:   `dsn for ent`,
				EnvVars: []string{`DSN`},
				Value:   `file:ent?mode=memory&cache=shared&_fk=1`,
			},
			&cli.BoolFlag{
				Name:    `skip-migrations`,
				Usage:   `skip migrations`,
				EnvVars: []string{`SKIP_MIGRATIONS`},
			},
		},
		Action: func(c *cli.Context) error {

			db, err := sql.Open(c.String(`driver`), c.String(`dsn`))
			if err != nil {
				return fmt.Errorf(`can't connect to database: %w`, err)
			}

			svc := s.service(entsql.OpenDB(c.String(`dialect`), db))

			for _, ext := range s.extensions {
				if runext, ok := ext.(RunnerExtension); ok {
					if err := runext.Runner(c, svc); err != nil {
						return fmt.Errorf(`can't execute extension: %w`, err)
					}
				}
			}

			if !c.Bool(`skip-migrations`) {
				if err := svc.MigrateSchema(c.Context, schema.WithGlobalUniqueID(s.configOptions.WithGlobalUniqueID)); err != nil {
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
