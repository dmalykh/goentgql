package directive

import (
	"context"
	"fmt"
	"github.com/99designs/gqlgen/graphql"
	"github.com/go-playground/validator/v10"
	"strings"
)

func NewValidator() *Directive {
	var validate = validator.New()
	{
		validate.RegisterValidation(`alphanumwith`, func(fl validator.FieldLevel) bool {
			var symbols = strings.Split(fl.Param(), "")
			var oldnew = make([]string, 0, len(symbols)*2)
			for i := 0; i < len(symbols); i++ {
				oldnew = append(oldnew, symbols[i], "")
			}

			return validate.Var(strings.NewReplacer(oldnew...).Replace(fl.Field().String()), "alphanum") == nil
		})
		validate.RegisterValidation(`alphanumunicodewith`, func(fl validator.FieldLevel) bool {
			var symbols = strings.Split(fl.Param(), "")
			var oldnew = make([]string, 0, len(symbols)*2)
			for i := 0; i < len(symbols); i++ {
				oldnew = append(oldnew, symbols[i], "")
			}

			return validate.Var(strings.NewReplacer(oldnew...).Replace(fl.Field().String()), "alphanumunicode") == nil
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
