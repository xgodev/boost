package logrus

import "github.com/sirupsen/logrus"

type Options struct {
	Formatter      logrus.Formatter // formatter TEXT/JSON/CLOUDWATCH
	ErrorFieldName string           // define field name for error logging
	Time           struct {
		Format string // date and time formats
	}
	Console struct {
		Enabled bool   // enable/disable console logging
		Level   string // console log level
	}
	Hooks []logrus.Hook
	File  struct {
		Enabled  bool   // enable/disable file logging
		Level    string // file log level
		Path     string // file log path
		Name     string // log filename
		MaxSize  int    // log file max size (MB)
		Compress bool   // enabled/disable file compress
		MaxAge   int    // file max age
	}
}

type Option func(options *Options)

func WithErrorFieldName(value string) Option {
	return func(options *Options) {
		options.ErrorFieldName = value
	}
}

func WithFormatter(value logrus.Formatter) Option {
	return func(options *Options) {
		options.Formatter = value
	}
}

func WithTimeFormat(value string) Option {
	return func(options *Options) {
		options.Time.Format = value
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

func WithHook(value logrus.Hook) Option {
	return func(options *Options) {
		options.Hooks = append(options.Hooks, value)
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
