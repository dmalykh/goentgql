package gqlgen

import (
	"github.com/99designs/gqlgen/api"
	gqlconfig "github.com/99designs/gqlgen/codegen/config"
	"github.com/dmalykh/goentgql/generator/config"
)

type GQLGen struct {
	config *config.GraphQLConfig
}

func New(config *config.GraphQLConfig) *GQLGen {
	return &GQLGen{
		config: config,
	}
}

func (g *GQLGen) Generate() error {
	if err := api.Generate((*gqlconfig.Config)(g.config)); err != nil {
		return err
	}

	return nil
}
