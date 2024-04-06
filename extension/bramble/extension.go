package bramble

import (
	"bytes"
	"context"
	gqlgen "github.com/99designs/gqlgen/codegen/config"
	"github.com/99designs/gqlgen/graphql"
	"github.com/dmalykh/entcontrib/entgql"
	"github.com/dmalykh/goentgql"
	"github.com/dmalykh/goentgql/generator/config"
	"github.com/urfave/cli/v2"
	"github.com/vektah/gqlparser/v2/ast"
	"github.com/vektah/gqlparser/v2/formatter"
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

	cfg.GraphQLConfig().Directives[`boundary`] = gqlgen.DirectiveConfig{
		SkipRuntime: false,
	}

	cfg.GraphQLConfig().Directives[`namespace`] = gqlgen.DirectiveConfig{
		SkipRuntime: false,
	}

	cfg.EntConfig().Extensions[`withRelaySpec`] = entgql.WithRelaySpec(true, map[string][]entgql.Directive{
		entgql.RelayCursor: {
			NamespaceDirective(),
		},
		entgql.RelayPageInfo: {
			NamespaceDirective(),
		},
		entgql.RelayNode: {
			NamespaceDirective(),
		},
		entgql.OrderDirectionEnum: {
			NamespaceDirective(),
		},
	})

	return nil
}

func (b *bramble) Runner(c *cli.Context, svc goentgql.Service) error {
	// Add directive
	svc.AddDirective(`Boundary`, func(ctx context.Context, obj interface{}, next graphql.Resolver) (res interface{}, err error) {
		return next(ctx)
	})

	// Add data to a service
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
