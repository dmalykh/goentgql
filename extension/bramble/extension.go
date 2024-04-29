package bramble

import (
	"bytes"
	"context"
	"entgo.io/ent/entc/gen"
	"fmt"
	gqlgen "github.com/99designs/gqlgen/codegen/config"
	"github.com/99designs/gqlgen/graphql"
	"github.com/dmalykh/entcontrib/entgql"
	"github.com/dmalykh/goentgql"
	"github.com/dmalykh/goentgql/generator/config"
	"github.com/urfave/cli/v2"
	"github.com/vektah/gqlparser/v2/ast"
	"github.com/vektah/gqlparser/v2/formatter"
)

var _ goentgql.GeneratorExtension = (*bramble)(nil)
var _ goentgql.RunnerExtension = (*bramble)(nil)

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

	cfg.EntConfig().GraphQLExtensions[`brambleExtension`] = entgql.WithSchemaHook(func(g *gen.Graph, s *ast.Schema) error {
		for _, node := range g.Nodes {
			if _, exist := node.Annotations[ConnectionAnnotationName]; !exist {
				continue
			}

			for _, typ := range s.Types {
				switch typ.Name {
				case entgql.RelayCursor:
					typ.Directives = append(typ.Directives, &ast.Directive{
						Name: `namespace`,
					})
				case entgql.RelayPageInfo:
					typ.Directives = append(typ.Directives, &ast.Directive{
						Name: `namespace`,
					})
				case entgql.RelayNode:
					typ.Directives = append(typ.Directives, &ast.Directive{
						Name: `namespace`,
					})
				case entgql.OrderDirectionEnum:
					typ.Directives = append(typ.Directives, &ast.Directive{
						Name: `namespace`,
					})
				case fmt.Sprintf(`%sConnection`, node.Name):
					typ.Directives = append(typ.Directives, &ast.Directive{
						Name: `boundary`,
					}, &ast.Directive{
						Name: `goModel`,
						Arguments: []*ast.Argument{
							{
								Name: `model`,
								Value: &ast.Value{
									Raw:  fmt.Sprintf(`%s/%sBrambleConnection`, cfg.EntConfig().Package, node.Name),
									Kind: ast.StringValue,
								},
							},
						},
					})
					typ.Fields = append(typ.Fields, &ast.FieldDefinition{
						Name:        `ID`,
						Type:        ast.NamedType(`ID`, nil),
						Description: "ID for connection",
					})
				}
			}
		}

		return nil
	})

	cfg.EntConfig().EntExtensions[`brambleConnectionIDExtension`] = &entExtension{}

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
