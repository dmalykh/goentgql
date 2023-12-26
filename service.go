package goentgql

import (
	"context"
	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/schema"
	"github.com/99designs/gqlgen/graphql"
)

type ServiceRunner func(drv *sql.Driver) Service

type Service interface {
	MigrationSchema(drv *sql.Driver) Migrator
	ExecutionSchema(drv *sql.Driver) graphql.ExecutableSchema
}

type Migrator interface {
	Create(ctx context.Context, opts ...schema.MigrateOption) error
}
