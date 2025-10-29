package loader

type Option func(*Loader)

// WithBaseDirectory returns an Option that sets the BaseDirectory
func WithBaseDirectory(dir string) Option {
	return func(l *Loader) {
		l.BaseDirectory = dir
	}
}

// WithOverrideDirectory returns an Option that sets the OverrideDirectory
func WithOverrideDirectory(dir string) Option {
	return func(l *Loader) {
		l.OverrideDirectory = dir
	}
}
