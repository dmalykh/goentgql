package service

var serviceTemplate = `
package {{ .PackageName }}

import (
	"context"
	"entgo.io/ent/dialect/sql/schema"
	"github.com/99designs/gqlgen/graphql"
	"github.com/dmalykh/goentgql"
	"github.com/dmalykh/goentgql/generator/gqlgen/directive"
	"reflect"

	gent "{{ .EntModulePath }}"
	graphqlgen "{{ .GraphQLModulePath }}"
)

func NewService(client *gent.Client, resolver graphqlgen.ResolverRoot) *Service {
	return &Service{
		resolver: resolver,
		migrator: client.Schema,
	}
}

type Service struct {
	resolver graphqlgen.ResolverRoot
	migrator goentgql.Migrator
}

func (s *Service) MigrateSchema(ctx context.Context, opts ...schema.MigrateOption) error {
	return s.migrator.Create(ctx, opts...)
}

func (s *Service) ExecutionSchema() graphql.ExecutableSchema {
	var directiveRoot = graphqlgen.DirectiveRoot{}
	if f := reflect.ValueOf(&directiveRoot).Elem().FieldByName("Validation"); f.CanSet() {
		f.Set(reflect.ValueOf(directive.NewValidator().Validation))
	}

	return graphqlgen.NewExecutableSchema(
		graphqlgen.Config{
			Resolvers:  s.resolver,
			Directives: directiveRoot,
		},
	)
}

`
