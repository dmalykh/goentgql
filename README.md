### GoEntGQL

A swiss knife to develop services uses ent+gqlgen.

```go
package main

import (
	"context"
	"entgo.io/ent/dialect/sql"
	"github.com/dmalykh/goentgql"
	"github.com/dmalykh/goentgql/extension/bramble"
	"log"
)

//go:generate go run main.go generate $PWD

func main() {
	var ctx = context.Background()
	var service = goentgql.New(
		goentgql.SchemaDir(`./schema`),
		goentgql.Package(`github.com/username/repository`),
		//goentgql.RunService(func(drv *sql.Driver) goentgql.Service {
		//	var client = gent.NewClient(gent.Driver(drv))
		//
		//	return generated.NewService(client, resolver.New(client))
		//}),
	)

	if err := service.Execute(ctx); err != nil {
		log.Fatal(err)
	}
}

```


TODO:

[] Refactoring: the current codebase was written as experiment and looks ugly

[] Test coverage

[] Make everything configurable

[] Make the run command configurable using https://github.com/knadh/koanf