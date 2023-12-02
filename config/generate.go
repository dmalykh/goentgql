package config

import (
	"dario.cat/mergo"
	"fmt"
	"github.com/99designs/gqlgen/codegen/config"
	"log"
	"path/filepath"
)

type ConfigGenerate struct {
	BasePath string
	Package  string
	Dir      string

	EntConfig     *EntConfig
	GraphQLConfig *GraphQLConfig
}

type ConfiguratorGenerate struct {
	config *ConfigGenerate
}

func NewGenerate(conf *ConfigGenerate) (*ConfiguratorGenerate, error) {
	var gen = new(ConfiguratorGenerate)
	log.Printf(`%+v`, conf)
	log.Println(`BasePath:`, conf.BasePath)
	log.Println(`Package:`, conf.Package)
	if conf.BasePath == `` || conf.Package == `` {
		return nil, fmt.Errorf(`path & package are required`)
	}

	// Initialize config with default values
	var cfg = new(ConfigGenerate)
	cfg.EntConfig = EntDefaultConfig()
	cfg.GraphQLConfig = GQLGenDefaultConfig()

	// Merge default config with config
	if err := mergo.Merge(cfg, *conf, mergo.WithOverride); err != nil {
		return nil, fmt.Errorf(`merge engen configs error: %w`, err)
	}

	cfg.GraphQLConfig.AutoBind = append(cfg.GraphQLConfig.AutoBind,
		filepath.Join(cfg.Package, cfg.Dir, cfg.EntConfig.Target))
	cfg.EntConfig.Package = filepath.Join(cfg.Package, cfg.Dir, cfg.EntConfig.Package)

	// Make paths for generated files absolute
	cfg.EntConfig.SchemaPath = filepath.Join(cfg.BasePath, cfg.EntConfig.SchemaPath)
	cfg.EntConfig.Target = filepath.Join(cfg.BasePath, cfg.Dir, cfg.EntConfig.Target)

	cfg.EntConfig.GraphQLSchemaOutputDir = filepath.Join(cfg.BasePath, cfg.Dir, cfg.EntConfig.GraphQLSchemaOutputDir)
	cfg.GraphQLConfig.SchemaFilename = append(
		config.StringList{filepath.Join(cfg.EntConfig.GraphQLSchemaOutputDir, cfg.EntConfig.GraphQLSchemaOutputFilename)},
		cfg.GraphQLConfig.SchemaFilename...,
	)

	// Completr gqlconfig
	gqlconfig, err := GQLConfig(cfg.GraphQLConfig)
	if err != nil {
		return nil, fmt.Errorf(`GQLConfig error: %w`, err)
	}

	cfg.GraphQLConfig = gqlconfig

	gen.config = cfg

	return gen, nil
}

func (c *ConfiguratorGenerate) EntConfig() *EntConfig {
	return c.config.EntConfig
}

func (c *ConfiguratorGenerate) GraphQLConfig() *GraphQLConfig {
	return c.config.GraphQLConfig
}
