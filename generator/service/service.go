package service

import (
	"fmt"
	"github.com/dmalykh/goentgql/generator/config"
	"html/template"
	"os"
)

type ServiceGenerator struct {
	config *config.ServiceConfig
}

func New(config *config.ServiceConfig) *ServiceGenerator {
	return &ServiceGenerator{
		config: config,
	}
}

func (g *ServiceGenerator) Generate() error {
	wr, err := os.OpenFile(g.config.OutputPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		return fmt.Errorf(`con't open the file %q: %w`, g.config.OutputPath, err)
	}

	defer func() {
		if err := wr.Close(); err != nil {
			panic(err)
		}
	}()

	tpl := template.Must(template.New("service-template").Parse(serviceTemplate))
	if err := tpl.Execute(wr, g.config); err != nil {
		return fmt.Errorf(`error generate service: %w`, err)
	}
	return nil
}
