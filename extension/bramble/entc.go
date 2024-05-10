package bramble

import (
	"entgo.io/ent/entc"
	"entgo.io/ent/entc/gen"
)

type GreetExtension struct {
	entc.DefaultExtension
}

func (*GreetExtension) Templates() []*gen.Template {
	return []*gen.Template{
		gen.MustParse(
			gen.NewTemplate("greet").Parse(`
{{ define "bramble_connection" }}
{{ $pkg := base $.Config.Package }}
{{ template "header" $ }}
	{{ range $n := $.Nodes }}
		{{ $receiver := $n.Receiver }}
		func ({{ $receiver }} *{{ $n.Name }}) Greet() string {
			return "Greetings, {{ $n.Name }}"
		}
	{{ end }}
{{ end }}`),
		),
	}
}
