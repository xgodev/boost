package zap

type Options struct {
	Console struct {
		Enabled   bool   // enable/disable console logging
		Level     string // console log level
		Formatter string // console formatter TEXT/JSON
	}
	File struct {
		Enabled   bool   // enable/disable file logging
		Level     string // file log level
		Path      string // file log path
		Name      string // file log filename
		MaxSize   int    // log file max size (MB)
		Compress  bool   // enabled/disable file compress
		MaxAge    int    // file max age
		Formatter string // file formatter TEXT/JSON
	}

	ErrorFieldName string // define field name for error logging
}

type Option func(options *Options)

func WithErrorFieldName(value string) Option {
	return func(options *Options) {
		options.ErrorFieldName = value
	}
}

func WithConsoleEnabled(value bool) Option {
	return func(options *Options) {
		options.Console.Enabled = value
	}
}

func WithConsoleLevel(value string) Option {
	return func(options *Options) {
		options.Console.Level = value
	}
}

func WithConsoleFormatter(value string) Option {
	return func(options *Options) {
		options.Console.Formatter = value
	}
}

func WithFileEnabled(value bool) Option {
	return func(options *Options) {
		options.File.Enabled = value
	}
}

func WithFileLevel(value string) Option {
	return func(options *Options) {
		options.File.Level = value
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

func WithFileFormatter(value string) Option {
	return func(options *Options) {
		options.File.Formatter = value
	}
}
