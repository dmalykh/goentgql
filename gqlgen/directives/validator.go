package directives

import (
	"context"
	"fmt"
	"github.com/99designs/gqlgen/graphql"
	"github.com/go-playground/validator/v10"
	"regexp"
)

func NewValidator() *Directive {
	var validate = validator.New()
	{
		validate.RegisterValidation(`slug`, func(fl validator.FieldLevel) bool {
			var r = regexp.MustCompile(`^[a-zA-Z0-9-]+$`)
			return r.MatchString(fl.Field().String())
		})
		validate.RegisterValidation(`title`, func(fl validator.FieldLevel) bool {
			var r = regexp.MustCompile(`^[\\p{L} \\p{N}]+$`)
			return r.MatchString(fl.Field().String())
		})
	}
	return &Directive{
		validate: validate,
	}
}

type Directive struct {
	validate *validator.Validate
}

func (d *Directive) Validation(ctx context.Context, obj interface{}, next graphql.Resolver, constraint string) (interface{}, error) {
	val, err := next(ctx)
	if err != nil {
		return nil, err
	}

	if err := d.validate.Var(val, constraint); err != nil {
		_, ok := err.(validator.ValidationErrors)
		if !ok {
			return nil, err
		}
		return nil, fmt.Errorf(`validation error in field "%s"`, *graphql.GetPathContext(ctx).Field)
	}

	return val, nil
}
