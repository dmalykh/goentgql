package bramble

import (
	"github.com/dmalykh/entcontrib/entgql"
)

// BoundaryDirective create `@boundary` directive to apply on the field/type
func BoundaryDirective() entgql.Directive {
	return entgql.NewDirective(`boundary`)
}

// NamespaceDirective create `@namespace` directive to apply on the field/type
func NamespaceDirective() entgql.Directive {
	return entgql.NewDirective(`namespace`)
}
