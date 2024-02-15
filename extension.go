package goentgql

import (
	"github.com/dmalykh/goentgql/generator/config"
	"github.com/urfave/cli/v2"
)

type RunnerExtension interface {
	Runner(c *cli.Context, svc Service) error
}

type GeneratorExtension interface {
	Generator(c *cli.Context, cfg *config.ConfiguratorGenerate) error
}

type Extension interface {
	Name() string
}
