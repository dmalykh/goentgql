package goentgql

import (
	"context"
	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/schema"
	"github.com/99designs/gqlgen/graphql"
)

type Option func(s *app)

type Migrator interface {
	Create(ctx context.Context, opts ...schema.MigrateOption) error
}

type dummyMigrator struct{}

func (*dummyMigrator) Create(ctx context.Context, opts ...schema.MigrateOption) error {
	return nil
}

type ExecutableSchema func(drv *sql.Driver) graphql.ExecutableSchema
type MigrationRunner func(drv *sql.Driver) Migrator

func SchemaDir(path string) Option {
	return func(s *app) {
		s.schemaDir = path
	}
}
func Package(name string) Option {
	return func(s *app) {
		s.packageName = name
	}
}

func RunMigrations(m MigrationRunner) Option {
	return func(s *app) {
		s.migrator = m
	}
}

func GqlSchema(schema ExecutableSchema) Option {
	return func(s *app) {
		s.executableSchema = schema
	}
}

func WithArgs(args []string) Option {
	return func(s *app) {
	}
}

func Args(args []string) Option {
	return func(s *app) {
	}
}

func Header(text string) Option {
	return func(s *app) {

	}
}
