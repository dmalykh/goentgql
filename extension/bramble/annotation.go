package bramble

import (
	"entgo.io/ent/schema"
	"github.com/dmalykh/entcontrib/entgql"
)

func BrambleBoundary() schema.Annotation {
	return entgql.Directives(BoundaryDirective(), NamespaceDirective())
}
