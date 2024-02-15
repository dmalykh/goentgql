package bramble

import (
	"bytes"
	"github.com/dmalykh/goentgql"
	"github.com/dmalykh/goentgql/generator/config"
	"github.com/urfave/cli/v2"
	"github.com/vektah/gqlparser/v2/ast"
	"github.com/vektah/gqlparser/v2/formatter"

	gqlgen "github.com/99designs/gqlgen/codegen/config"
)

//var _ goentgql.GeneratorExtension = nil.(*bramble)

type bramble struct {
	service *Service
}

func (b *bramble) Generator(c *cli.Context, cfg *config.ConfiguratorGenerate) error {
	cfg.GraphQLConfig().Sources = append([]*ast.Source{
		{Name: `bramble`, Input: `
			type Service {
			  name: String! # unique name for the service
			  version: String! # any string
			  schema: String! # the full schema for the service
			}
			
			extend type Query {
			  service: Service!
			}`},
	}, cfg.GraphQLConfig().Sources...)

	cfg.GraphQLConfig().Models[`Service`] = gqlgen.TypeMapEntry{
		Model: gqlgen.StringList{
			`github.com/dmalykh/goentgql/extension/bramble.Service`,
		},
		ForceGenerate: false,
	}

	return nil
}

func (b *bramble) Runner(c *cli.Context, svc goentgql.Service) error {
	var buf bytes.Buffer
	formatter.NewFormatter(&buf).FormatSchema(svc.ExecutionSchema().Schema())

	service.Name = b.service.Name
	service.Version = b.service.Version
	service.Schema = buf.String()

	return nil
}

func (b *bramble) Name() string {
	return "bramble"
}

func Bramble(name, version string) goentgql.Option {
	return goentgql.AddExtension(&bramble{
		service: &Service{
			Name:    name,
			Version: version,
		},
	})
	//return goentgql.AddExtension(func(c *cli.Context, svc goentgql.Service) error {

	//
	//	return nil
	//})
}
