package runner

import (
	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/debug"
	"github.com/99designs/gqlgen/graphql/playground"
	"net/http"
)

type Config struct {
	Addr    string
	Name    string
	Verbose bool
}

func Run(executableSchema graphql.ExecutableSchema, cfg *Config, middleware ...graphql.OperationMiddleware) error {
	var srv = handler.NewDefaultServer(executableSchema)
	if cfg.Verbose {
		srv.Use(&debug.Tracer{})
	}

	// Handle playground
	http.Handle(`/`, playground.Handler(cfg.Name, `/query`))

	{
		for _, m := range middleware {
			srv.AroundOperations(m)
		}
		http.Handle(`/query`, srv)
	}

	if err := http.ListenAndServe(cfg.Addr, nil); err != nil {
		return err
	}

	return nil
}
