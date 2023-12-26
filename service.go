package goentgql

import (
	"context"
	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/schema"
	"github.com/99designs/gqlgen/graphql"
)

type ServiceRunner func(drv *sql.Driver) Service

type Service interface {
	MigrateSchema(ctx context.Context, opts ...schema.MigrateOption) error
	ExecutionSchema() graphql.ExecutableSchema
}
