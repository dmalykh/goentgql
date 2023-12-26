package config

type EntConfig struct {
	SchemaPath                  string
	Header                      string
	Target                      string
	GraphQLSchemaOutputDir      string
	GraphQLSchemaOutputFilename string
	ModuleName                  string
}

func EntDefaultConfig() *EntConfig {
	return &EntConfig{
		SchemaPath:                  `/schema`,
		Target:                      `/gent`,
		GraphQLSchemaOutputDir:      `/schema/`,
		GraphQLSchemaOutputFilename: `ent.generated.graphql`,
		ModuleName:                  `gent`,
	}
}
