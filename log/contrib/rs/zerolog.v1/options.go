package zerolog

type Options struct {
	Formatter string // formatter TEXT/JSON
	Level     string // log level

	Console struct {
		Enabled bool // enable/disable console logging
	}
	File struct {
		Enabled  bool   // enable/disable file logging
		Path     string // file log path
		Name     string // file log filename
		MaxSize  int    // file log file max size (MB)
		Compress bool   // enabled/disable file compress
		MaxAge   int    // file max age
	}

	ErrorFieldName string // define field name for error logging
}

type Option func(options *Options)

func WithErrorFieldName(value string) Option {
	return func(options *Options) {
		options.ErrorFieldName = value
	}
}

func WithFormatter(value string) Option {
	return func(options *Options) {
		options.Formatter = value
	}
}

func WithConsoleEnabled(value bool) Option {
	return func(options *Options) {
		options.Console.Enabled = value
	}
}

func WithLevel(value string) Option {
	return func(options *Options) {
		options.Level = value
	}
}

func WithFileEnabled(value bool) Option {
	return func(options *Options) {
		options.File.Enabled = value
	}
}

func WithFilePath(value string) Option {
	return func(options *Options) {
		options.File.Path = value
	}
}

func WithFileName(value string) Option {
	return func(options *Options) {
		options.File.Name = value
	}
}

func WithFileMaxSize(value int) Option {
	return func(options *Options) {
		options.File.MaxSize = value
	}
}

func WithFileCompress(value bool) Option {
	return func(options *Options) {
		options.File.Compress = value
	}
}

func WithFileMaxAge(value int) Option {
	return func(options *Options) {
		options.File.MaxAge = value
	}
}
