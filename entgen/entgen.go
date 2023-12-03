package entgen

import (
	"entgo.io/ent/entc"
	"entgo.io/ent/entc/gen"
	"entgo.io/ent/schema/field"
	"fmt"
	gqlconfig "github.com/99designs/gqlgen/codegen/config"
	"github.com/dmalykh/entcontrib/entgql"
	"github.com/dmalykh/goentgql/config"
	"os"
	"path/filepath"
)

type Generator struct {
	basePath string
	config   *config.EntConfig
	gqlgen   *config.GraphQLConfig
}

func New(conf *config.EntConfig, gqlgen *config.GraphQLConfig) (*Generator, error) {
	if conf.SchemaPath == `` || conf.Target == `` {
		return nil, fmt.Errorf(`schemapath & target are required`)
	}

	// Create absent directories for generated files
	if err := os.MkdirAll(filepath.Dir(conf.Target), os.ModePerm); err != nil {
		return nil, fmt.Errorf(`create dir %q error: %w`, filepath.Dir(conf.Target), err)
	}
	if err := os.MkdirAll(conf.GraphQLSchemaOutputDir, os.ModePerm); err != nil {
		return nil, fmt.Errorf(`create dir %q error: %w`, conf.GraphQLSchemaOutputDir, err)
	}

	return &Generator{
		config: conf,
		gqlgen: gqlgen,
	}, nil
}

func (g *Generator) Generate() error {
	var graphqlOutput = filepath.Join(g.config.GraphQLSchemaOutputDir, g.config.GraphQLSchemaOutputFilename)
	ex, err := entgql.NewExtension(
		entgql.WithCompletedConfig((*gqlconfig.Config)(g.gqlgen)),
		entgql.WithSchemaPath(graphqlOutput),
		entgql.WithWhereInputs(false),
		entgql.WithRelaySpec(true),
		entgql.WithSchemaGenerator(),
		entgql.WithEmptyQueryGenerator(),
	)
	if err != nil {
		return fmt.Errorf("creating entgql extension: %w", err)
	}

	if err := entc.Generate(g.config.SchemaPath, &gen.Config{
		Header:  g.config.Header,
		IDType:  &field.TypeInfo{Type: field.TypeInt64},
		Target:  g.config.Target,
		Package: g.config.Package,
	},
		entc.Extensions(ex),
	); err != nil {
		return fmt.Errorf("running ent codegen: %w", err)
	}

	return nil
}
