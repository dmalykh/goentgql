package goentgql

type Option func(s *App)

func SchemaDir(path string) Option {
	return func(s *App) {
		s.schemaDir = path
	}
}
func Package(name string) Option {
	return func(s *App) {
		s.packageName = name
	}
}

func RunService(f ServiceRunner) Option {
	return func(s *App) {
		s.service = f
	}
}

func IDType(idType string) Option {
	return func(s *App) {
		s.configOptions.IDType = append(s.configOptions.IDType, idType)
	}
}

func WithGlobalUniqueID(yes bool) Option {
	return func(s *App) {
		s.configOptions.WithGlobalUniqueID = yes
	}
}

func AddExtension(ext Extension) Option {
	return func(s *App) {
		s.extensions = append(s.extensions, ext)
	}
}

func WithArgs(args []string) Option {
	return func(s *App) {
	}
}

func Args(args []string) Option {
	return func(s *App) {
	}
}
