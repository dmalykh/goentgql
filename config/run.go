package config

type Client struct {
	Addr           string `name:"port" default:"8081" help:"Address to listen on"`
	PrometheusAddr string `name:"promp" default:"2112" help:"Address to listen on prometheus metrics"`
	Verbose        bool   `name:"verbose" help:"Enable verbose mode"`
	Driver         string `name:"driver" help:"SQL driver name"`
	DSN            string `name:"dsn" help:"SQL dsn"`
}

//kong.Parse(&cli)
