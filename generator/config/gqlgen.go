package config

import (
	"dario.cat/mergo"
	"fmt"
	"github.com/99designs/gqlgen/codegen/config"
	"github.com/vektah/gqlparser/v2/ast"
)

type GraphQLConfig config.Config

func GQLConfig(conf *GraphQLConfig) (*GraphQLConfig, error) {
	cfg := GQLGenDefaultConfig()

	if err := mergo.Merge(cfg, conf,
		mergo.WithOverride, mergo.WithAppendSlice, mergo.WithoutDereference); err != nil {
		return nil, fmt.Errorf(`merge gqlgen configs error: %w`, err)
	}

	if err := config.CompleteConfig((*config.Config)(conf)); err != nil {
		return nil, fmt.Errorf(`complete gqlgen configs error: %w`, err)
	}

	return conf, nil
}

func GQLGenDefaultConfig() *GraphQLConfig {
	return &GraphQLConfig{
		StructFieldsAlwaysPointers:    true,
		ReturnPointersInUmarshalInput: false,
		ResolversAlwaysReturnPointers: true,
		NullableInputOmittable:        false,

		Sources: []*ast.Source{
			{
				Name:  `directive_validation`,
				Input: `directive @validation(rule: String!) on INPUT_FIELD_DEFINITION | FIELD_DEFINITION`,
			},
		},
		SchemaFilename: config.StringList{
			`schema/*.graphql`,
		},
		Exec: config.ExecConfig{
			Filename: `generated/graphqlgen/generated.go`,
			Package:  `graphqlgen`,
		},
		Model: config.PackageConfig{
			Filename: `generated/graphqlgen/genmodel/models.go`,
			Package:  `genmodel`,
			Version:  2,
		},
		Resolver: config.ResolverConfig{
			Layout:  config.LayoutFollowSchema,
			DirName: `./generated/graphqlgen`,
			Package: `graphqlgen`,
		},
		Directives: map[string]config.DirectiveConfig{
			`validation`: {SkipRuntime: false},
		},
		AutoBind: []string{},
		Models: config.TypeMap{
			`ID`: config.TypeMapEntry{
				Model: config.StringList{
					`github.com/99designs/gqlgen/graphql.String`,
				},
			},
			`Int`: config.TypeMapEntry{
				Model: config.StringList{
					`github.com/99designs/gqlgen/graphql.Int`,
					`github.com/99designs/gqlgen/graphql.Int64`,
				},
			},
			`Node`: config.TypeMapEntry{
				Model: config.StringList{
					`gent.Noder`,
				},
			},
		},
	}
}
