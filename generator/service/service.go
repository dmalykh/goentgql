package service

import (
	"fmt"
	"github.com/dmalykh/goentgql/generator/config"
	"html/template"
	"io"
)

type ServiceGenerator struct {
	config *config.ServiceConfig
	wr     io.Writer
}

func New(config *config.ServiceConfig) *ServiceGenerator {
	return &ServiceGenerator{
		config: config,
	}
}

func (g *ServiceGenerator) Generate() error {
	tpl := template.Must(template.New("service-template").Parse(serviceTemplate))
	if err := tpl.Execute(g.wr, g.config); err != nil {
		return fmt.Errorf(`error generate service: %w`, err)
	}
	return nil
}
