package goentgql

type Option func(s *app)

func SchemaDir(path string) Option {
	return func(s *app) {
		s.schemaDir = path
	}
}
func Package(name string) Option {
	return func(s *app) {
		s.packageName = name
	}
}

func RunService(f ServiceRunner) Option {
	return func(s *app) {
		s.service = f
	}
}

func WithArgs(args []string) Option {
	return func(s *app) {
	}
}

func Args(args []string) Option {
	return func(s *app) {
	}
}
