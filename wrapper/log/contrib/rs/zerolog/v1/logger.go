package zerolog

import (
	"bytes"
	"context"
	"io"
	"os"
	"reflect"
	"strings"

	"github.com/rs/zerolog"
	"github.com/xgodev/boost/wrapper/log"
	"gopkg.in/natefinch/lumberjack.v2"
)

type ctxKey string

const (
	key                   ctxKey = "ctxfields"
	defaultFormatter             = "TEXT"
	defaultLevel                 = "INFO"
	defaultConsoleEnabled        = true
	defaultFileEnabled           = false
	defaultFilePath              = "/tmp"
	defaultFileName              = "application.log"
	defaultFileMaxSize           = 100
	defaultFileCompress          = true
	defaultFileMaxAge            = 28
	defaultErrorFieldName        = "err"
)

// NewLogger constructs a new Logger from provided variadic Option.
func NewLogger(option ...Option) log.Logger {
	options := options(option)
	return NewLoggerWithOptions(options)
}

// NewLoggerWithOptions constructs a new Logger from provided Options.
func NewLoggerWithOptions(options *Options) log.Logger {
	writer := getWriter(options)
	if writer == nil {
		zerologger := zerolog.Nop()
		logger := &logger{
			logger: zerologger,
		}

		log.Set(logger)
		return logger
	}

	zerolog.MessageFieldName = "log_message"
	zerolog.LevelFieldName = "log_level"

	zerologger := zerolog.New(writer).With().Timestamp().Logger()
	level := logLevel(options.Level)
	zerologger = zerologger.Level(level)

	// Default options are only applied if this is called via NewLogger
	// If called direct, the options passed to this function may be empty.
	// Hence the default is reinforced here.
	errorField := options.ErrorFieldName
	if errorField == "" {
		errorField = defaultErrorFieldName
	}

	logger := &logger{
		logger:         zerologger,
		writer:         writer,
		fields:         log.Fields{},
		errorFieldName: errorField,
	}

	log.Set(logger)
	return logger
}

func defaultOptions() *Options {
	return &Options{
		Formatter:      defaultFormatter,
		Level:          defaultLevel,
		ErrorFieldName: defaultErrorFieldName,

		Console: struct {
			Enabled bool
		}{
			Enabled: defaultConsoleEnabled,
		},
		File: struct {
			Enabled  bool
			Path     string
			Name     string
			MaxSize  int
			Compress bool
			MaxAge   int
		}{
			Enabled:  defaultFileEnabled,
			Path:     defaultFilePath,
			Name:     defaultFileName,
			MaxSize:  defaultFileMaxSize,
			Compress: defaultFileCompress,
			MaxAge:   defaultFileMaxAge,
		},
	}
}

func options(option []Option) *Options {
	options := defaultOptions()

	for _, o := range option {
		o(options)
	}
	return options
}

type logger struct {
	logger         zerolog.Logger
	writer         io.Writer
	fields         log.Fields
	errorFieldName string
}

func logLevel(level string) zerolog.Level {
	switch level {
	case "TRACE":
		return zerolog.TraceLevel
	case "DEBUG":
		return zerolog.DebugLevel
	case "WARN":
		return zerolog.WarnLevel
	case "ERROR":
		return zerolog.ErrorLevel
	case "PANIC":
		return zerolog.PanicLevel
	case "FATAL":
		return zerolog.FatalLevel
	default:
		return zerolog.InfoLevel
	}
}

func getWriter(options *Options) io.Writer {
	var writer io.Writer
	switch options.Formatter {
	case "TEXT":
		writer = zerolog.ConsoleWriter{Out: os.Stdout}
	default:
		writer = os.Stdout
	}

	if options.File.Enabled {
		s := []string{options.File.Path, "/", options.File.Name}
		fileLocation := strings.Join(s, "")

		fileHandler := &lumberjack.Logger{
			Filename: fileLocation,
			MaxSize:  options.File.MaxSize,
			Compress: options.File.Compress,
			MaxAge:   options.File.MaxAge,
		}

		if options.Console.Enabled {
			return io.MultiWriter(writer, fileHandler)
		}
		return fileHandler
	} else if options.Console.Enabled {
		return writer
	}

	return nil
}

func (l *logger) Printf(format string, args ...interface{}) {
	l.logger.Printf(format, args...)
}

func (l *logger) Tracef(format string, args ...interface{}) {
	l.logger.Trace().Msgf(format, args...)
}

func (l *logger) Trace(args ...interface{}) {
	format := bytes.NewBufferString("")
	for range args {
		format.WriteString("%v")
	}

	l.logger.Trace().Msgf(format.String(), args...)
}

func (l *logger) Debugf(format string, args ...interface{}) {
	l.logger.Debug().Msgf(format, args...)
}

func (l *logger) Debug(args ...interface{}) {
	format := bytes.NewBufferString("")
	for range args {
		format.WriteString("%v")
	}

	l.logger.Debug().Msgf(format.String(), args...)
}

func (l *logger) Infof(format string, args ...interface{}) {
	l.logger.Info().Msgf(format, args...)
}

func (l *logger) Info(args ...interface{}) {
	format := bytes.NewBufferString("")
	for range args {
		format.WriteString("%v")
	}

	l.logger.Info().Msgf(format.String(), args...)
}

func (l *logger) Warnf(format string, args ...interface{}) {
	l.logger.Warn().Msgf(format, args...)
}

func (l *logger) Warn(args ...interface{}) {
	format := bytes.NewBufferString("")
	for range args {
		format.WriteString("%v")
	}

	l.logger.Warn().Msgf(format.String(), args...)
}

func (l *logger) Errorf(format string, args ...interface{}) {
	l.logger.Error().Msgf(format, args...)
}

func (l *logger) Error(args ...interface{}) {
	format := bytes.NewBufferString("")
	for range args {
		format.WriteString("%v")
	}

	l.logger.Error().Msgf(format.String(), args...)
}

func (l *logger) Fatalf(format string, args ...interface{}) {
	l.logger.Fatal().Msgf(format, args...)
}

func (l *logger) Fatal(args ...interface{}) {
	format := bytes.NewBufferString("")
	for range args {
		format.WriteString("%v")
	}

	l.logger.Fatal().Msgf(format.String(), args...)
}

func (l *logger) Panicf(format string, args ...interface{}) {
	l.logger.Panic().Msgf(format, args...)
}

func (l *logger) Panic(args ...interface{}) {
	format := bytes.NewBufferString("")
	for range args {
		format.WriteString("%v")
	}

	l.logger.Panic().Msgf(format.String(), args...)
}

func (l *logger) WithField(key string, value interface{}) log.Logger {
	newField := make(map[string]interface{})
	newField[key] = value

	newLogger := l.logger.With().Fields(newField).Logger()
	return &logger{newLogger, l.writer, newField, l.errorFieldName}
}

func (l *logger) WithFields(fields map[string]interface{}) log.Logger {
	newLogger := l.logger.With().Fields(fields).Logger()
	return &logger{newLogger, l.writer, fields, l.errorFieldName}
}

func (l *logger) WithTypeOf(obj interface{}) log.Logger {

	t := reflect.TypeOf(obj)

	return l.WithFields(map[string]interface{}{
		"reflect.type.name":    t.Name(),
		"reflect.type.package": t.PkgPath(),
	})
}

func (l *logger) WithError(err error) log.Logger {
	return l.WithField(l.errorFieldName, err.Error())
}

func (l *logger) Fields() log.Fields {
	return l.fields
}

func (l *logger) Output() io.Writer {
	return l.writer
}

func (l *logger) ToContext(ctx context.Context) context.Context {
	return l.logger.WithContext(context.WithValue(ctx, key, l.fields))
}

func (l *logger) FromContext(ctx context.Context) log.Logger {
	zerologger := zerolog.Ctx(ctx)
	if zerologger.GetLevel() == zerolog.Disabled {
		return l
	}
	rawFields := ctx.Value(key)
	fields := log.Fields{}
	if rawFields != nil {
		switch v := rawFields.(type) {
		case log.Fields:
			fields = v
		}
	}
	return &logger{*zerologger, l.writer, fields, l.errorFieldName}
}
