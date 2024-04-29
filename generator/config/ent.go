package config

import (
	"entgo.io/ent/entc"
	"entgo.io/ent/entc/gen"
	"entgo.io/ent/schema/field"
	"github.com/dmalykh/entcontrib/entgql"
)

type EntConfig struct {
	*gen.Config

	SchemaPath                  string
	GraphQLSchemaOutputDir      string
	GraphQLSchemaOutputFilename string

	// Extension as a map to allow manage them in a more flexible way using extensions
	GraphQLExtensions map[string]entgql.ExtensionOption
	EntExtensions     map[string]entc.Extension
}

func EntDefaultConfig() *EntConfig {
	return &EntConfig{
		Config: &gen.Config{
			Target:  `/gent`,
			Package: `gent`,
			IDType:  &field.TypeInfo{Type: field.TypeInt64},
		},
		SchemaPath:                  `/schema`,
		GraphQLSchemaOutputDir:      `/schema/`,
		GraphQLSchemaOutputFilename: `ent.generated.graphql`,
		GraphQLExtensions: map[string]entgql.ExtensionOption{
			`withWhereInputs`:         entgql.WithWhereInputs(false),
			`withRelaySpec`:           entgql.WithRelaySpec(true, nil),
			`withSchemaGenerator`:     entgql.WithSchemaGenerator(),
			`withEmptyQueryGenerator`: entgql.WithEmptyQueryGenerator(),
		},
		EntExtensions: map[string]entc.Extension{},
	}
}
