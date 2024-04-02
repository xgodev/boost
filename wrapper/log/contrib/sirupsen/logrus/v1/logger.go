package logrus

import (
	"context"
	"github.com/xgodev/boost/wrapper/log/contrib/sirupsen/logrus/v1/formatter/text"
	"io"
	"io/ioutil"
	"os"
	"reflect"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/xgodev/boost/wrapper/log"
	"gopkg.in/natefinch/lumberjack.v2"
)

type ctxKey string

const (
	key                   ctxKey = "ctxfields"
	defaultConsoleEnabled        = true
	defaultConsoleLevel          = "INFO"
	defaultFileEnabled           = false
	defaultFileLevel             = "INFO"
	defaultFilePath              = "/tmp"
	defaultFileName              = "application.log"
	defaultFileMaxSize           = 100
	defaultFileCompress          = true
	defaultFileMaxAge            = 28
	defaultTimeFormat            = "2006/01/02 15:04:05.000"
	defaultErrorFieldName        = "err"
)

// NewLogger constructs a new Logger from provided variadic Option.
func NewLogger(option ...Option) log.Logger {
	options := options(option)
	return NewLoggerWithOptions(options)
}

// NewLoggerWithOptions constructs a new Logger from provided Options.
func NewLoggerWithOptions(options *Options) log.Logger {

	lLogger := new(logrus.Logger)

	for _, hook := range options.Hooks {
		// init level hooks
		lLogger.Hooks = logrus.LevelHooks{}
		lLogger.AddHook(hook)
	}

	var fileHandler *lumberjack.Logger

	lLogger.SetOutput(ioutil.Discard)

	if options.File.Enabled {

		s := []string{options.File.Path, "/", options.File.Name}
		fileLocation := strings.Join(s, "")

		fileHandler = &lumberjack.Logger{
			Filename: fileLocation,
			MaxSize:  options.File.MaxSize,
			Compress: options.File.Compress,
			MaxAge:   options.File.MaxAge,
		}

	}

	if options.Console.Enabled && options.File.Enabled {
		lLogger.SetOutput(io.MultiWriter(os.Stdout, fileHandler))
	} else if options.File.Enabled {
		lLogger.SetOutput(fileHandler)
	} else if options.Console.Enabled {
		lLogger.SetOutput(os.Stdout)
	}

	level := logLevel(options.Console.Level)
	lLogger.SetLevel(level)

	lLogger.SetFormatter(options.Formatter)

	// Default options are only applied if this is called via NewLogger
	// If called direct, the options passed to this function may be empty.
	// Hence the default is reinforced here.
	errorField := options.ErrorFieldName
	if errorField == "" {
		errorField = defaultErrorFieldName
	}

	logger := &logger{
		logger:         lLogger,
		fields:         log.Fields{},
		errorFieldName: errorField,
	}

	log.SetGlobalLogger(logger)
	return logger
}

func defaultOptions() *Options {
	return &Options{
		Formatter:      text.New(),
		ErrorFieldName: defaultErrorFieldName,
		Time: struct {
			Format string
		}{
			Format: defaultTimeFormat,
		},
		Console: struct {
			Enabled bool
			Level   string
		}{
			Enabled: defaultConsoleEnabled,
			Level:   defaultConsoleLevel,
		},
		File: struct {
			Enabled  bool
			Level    string
			Path     string
			Name     string
			MaxSize  int
			Compress bool
			MaxAge   int
		}{
			Enabled:  defaultFileEnabled,
			Level:    defaultFileLevel,
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

func logLevel(level string) logrus.Level {

	switch level {

	case "DEBUG":
		return logrus.DebugLevel
	case "WARN":
		return logrus.WarnLevel
	case "FATAL":
		return logrus.FatalLevel
	case "ERROR":
		return logrus.ErrorLevel
	case "TRACE":
		return logrus.TraceLevel
	case "PANIC":
		return logrus.PanicLevel
	default:
		return logrus.InfoLevel
	}

}

type logger struct {
	logger         *logrus.Logger
	fields         log.Fields
	errorFieldName string
}

func (l *logger) Trace(args ...interface{}) {
	l.logger.Trace(args...)
}

func (l *logger) Debug(args ...interface{}) {
	l.logger.Debug(args...)
}

func (l *logger) Info(args ...interface{}) {
	l.logger.Info(args...)
}

func (l *logger) Warn(args ...interface{}) {
	l.logger.Warn(args...)
}

func (l *logger) Error(args ...interface{}) {
	l.logger.Error(args...)
}

func (l *logger) Panic(args ...interface{}) {
	l.logger.Panic(args...)
}

func (l *logger) Fatal(args ...interface{}) {
	l.logger.Fatal(args...)
}

func (l *logger) Printf(format string, args ...interface{}) {
	l.logger.Printf(format, args...)
}

func (l *logger) Tracef(format string, args ...interface{}) {
	l.logger.Tracef(format, args...)
}

func (l *logger) Debugf(format string, args ...interface{}) {
	l.logger.Debugf(format, args...)
}

func (l *logger) Infof(format string, args ...interface{}) {
	l.logger.Infof(format, args...)
}

func (l *logger) Warnf(format string, args ...interface{}) {
	l.logger.Warnf(format, args...)
}

func (l *logger) Errorf(format string, args ...interface{}) {
	l.logger.Errorf(format, args...)
}

func (l *logger) Panicf(format string, args ...interface{}) {
	l.logger.Panicf(format, args...)
}

func (l *logger) Fatalf(format string, args ...interface{}) {
	l.logger.Fatalf(format, args...)
}

func (l *logger) WithField(key string, value interface{}) log.Logger {

	entry := l.logger.WithField(key, value)

	return &logEntry{
		entry:  entry,
		fields: convertToFields(entry.Data),
	}
}

func (l *logger) WithFields(fields map[string]interface{}) log.Logger {
	return &logEntry{
		entry:  l.logger.WithFields(convertToLogrusFields(fields)),
		fields: fields,
	}
}

func (l *logger) WithTypeOf(obj interface{}) log.Logger {

	t := reflect.TypeOf(obj)

	return l.WithFields(log.Fields{
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
	return l.logger.Out
}

func (l *logger) ToContext(ctx context.Context) context.Context {
	return toContext(ctx, l.fields)
}

func (l *logger) FromContext(ctx context.Context) log.Logger {
	fields := fieldsFromContext(ctx)
	return l.WithFields(fields)
}

type logEntry struct {
	entry          *logrus.Entry
	fields         map[string]interface{}
	errorFieldName string
}

func (l *logEntry) Trace(args ...interface{}) {
	l.entry.Trace(args...)
}

func (l *logEntry) Debug(args ...interface{}) {
	l.entry.Debug(args...)
}

func (l *logEntry) Info(args ...interface{}) {
	l.entry.Info(args...)
}

func (l *logEntry) Warn(args ...interface{}) {
	l.entry.Warn(args...)
}

func (l *logEntry) Error(args ...interface{}) {
	l.entry.Error(args...)
}

func (l *logEntry) Fatal(args ...interface{}) {
	l.entry.Fatal(args...)
}

func (l *logEntry) Panic(args ...interface{}) {
	l.entry.Panic(args...)
}

func (l *logEntry) WithField(key string, value interface{}) log.Logger {

	entry := l.entry.WithField(key, value)

	return &logEntry{
		entry:          entry,
		fields:         convertToFields(entry.Data),
		errorFieldName: l.errorFieldName,
	}
}

func (l *logEntry) Fields() log.Fields {
	return l.fields
}

func (l *logEntry) Output() io.Writer {
	return l.entry.Logger.Out
}

func (l *logEntry) Printf(format string, args ...interface{}) {
	l.entry.Printf(format, args...)
}

func (l *logEntry) Tracef(format string, args ...interface{}) {
	l.entry.Tracef(format, args...)
}

func (l *logEntry) Debugf(format string, args ...interface{}) {
	l.entry.Debugf(format, args...)
}

func (l *logEntry) Infof(format string, args ...interface{}) {
	l.entry.Infof(format, args...)
}

func (l *logEntry) Warnf(format string, args ...interface{}) {
	l.entry.Warnf(format, args...)
}

func (l *logEntry) Errorf(format string, args ...interface{}) {
	l.entry.Errorf(format, args...)
}

func (l *logEntry) Panicf(format string, args ...interface{}) {
	l.entry.Panicf(format, args...)
}

func (l *logEntry) Fatalf(format string, args ...interface{}) {
	l.entry.Fatalf(format, args...)
}

func (l *logEntry) WithFields(fields map[string]interface{}) log.Logger {
	entry := l.entry.WithFields(convertToLogrusFields(fields))
	return &logEntry{
		entry:          entry,
		fields:         convertToFields(entry.Data),
		errorFieldName: l.errorFieldName,
	}
}

func (l *logEntry) WithTypeOf(obj interface{}) log.Logger {

	t := reflect.TypeOf(obj)

	return l.WithFields(log.Fields{
		"reflect.type.name":    t.Name(),
		"reflect.type.package": t.PkgPath(),
	})
}

func (l *logEntry) WithError(err error) log.Logger {
	return l.WithField(l.errorFieldName, err.Error())
}

func (l *logEntry) ToContext(ctx context.Context) context.Context {
	return toContext(ctx, l.fields)
}

func (l *logEntry) FromContext(ctx context.Context) log.Logger {
	fields := fieldsFromContext(ctx)
	return l.WithFields(fields)
}

func toContext(ctx context.Context, fields log.Fields) context.Context {
	ctxFields := fieldsFromContext(ctx)

	for k, v := range fields {
		ctxFields[k] = v
	}

	return context.WithValue(ctx, key, ctxFields)
}

func fieldsFromContext(ctx context.Context) log.Fields {
	fields := make(log.Fields)

	if ctx == nil {
		return fields
	}

	if f, ok := ctx.Value(key).(log.Fields); ok && f != nil {
		for k, v := range f {
			fields[k] = v
		}
	}

	return fields
}

func convertToLogrusFields(fields log.Fields) logrus.Fields {
	logrusFields := logrus.Fields{}
	for index, val := range fields {
		logrusFields[index] = val
	}
	return logrusFields
}

func convertToFields(logrusFields logrus.Fields) log.Fields {
	fields := make(map[string]interface{})
	for index, val := range logrusFields {
		fields[index] = val
	}
	return fields
}
