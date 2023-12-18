package gqlgen

import (
	"github.com/dmalykh/entcontrib/entgql"
	"github.com/vektah/gqlparser/v2/ast"
)

// Validation create `@validation` directive to apply on the field/type
func Validation(rule string) entgql.Directive {
	var args []*ast.Argument
	if rule != "" {
		args = append(args, &ast.Argument{
			Name: "",
			Value: &ast.Value{
				Raw:  `rule`,
				Kind: ast.StringValue,
			},
		})
	}
	return entgql.NewDirective("validation", args...)
}
