package config

import (
	"github.com/dmalykh/entcontrib/entgql"
)

type EntConfig struct {
	SchemaPath                  string
	Header                      string
	Target                      string
	GraphQLSchemaOutputDir      string
	GraphQLSchemaOutputFilename string
	ModuleName                  string
	Extensions                  map[string]entgql.ExtensionOption
}

func EntDefaultConfig() *EntConfig {
	return &EntConfig{
		SchemaPath:                  `/schema`,
		Target:                      `/gent`,
		GraphQLSchemaOutputDir:      `/schema/`,
		GraphQLSchemaOutputFilename: `ent.generated.graphql`,
		ModuleName:                  `gent`,
		Extensions: map[string]entgql.ExtensionOption{
			`withWhereInputs`:         entgql.WithWhereInputs(false),
			`withRelaySpec`:           entgql.WithRelaySpec(true, nil),
			`withSchemaGenerator`:     entgql.WithSchemaGenerator(),
			`withEmptyQueryGenerator`: entgql.WithEmptyQueryGenerator(),
		},
	}
}
