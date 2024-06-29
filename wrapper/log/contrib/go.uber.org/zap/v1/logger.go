package zap

import (
	"context"
	"io"
	"os"
	"reflect"
	"strings"

	"github.com/xgodev/boost/wrapper/log"
	"gopkg.in/natefinch/lumberjack.v2"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type ctxKey string

const (
	key                     ctxKey = "ctxfields"
	defaultConsoleFormatter        = "TEXT"
	defaultConsoleEnabled          = true
	defaultConsoleLevel            = "INFO"
	defaultFileEnabled             = false
	defaultFileLevel               = "INFO"
	defaultFilePath                = "/tmp"
	defaultFileName                = "application.log"
	defaultFileMaxSize             = 100
	defaultFileCompress            = true
	defaultFileMaxAge              = 28
	defaultFileFormatter           = "TEXT"
	defaultErrorFieldName          = "err"
)

// NewLogger constructs a new Logger from provided variadic Option.
func NewLogger(option ...Option) log.Logger {
	options := options(option)
	return NewLoggerWithOptions(options)
}

// NewLoggerWithOptions constructs a new Logger from provided Options.
func NewLoggerWithOptions(options *Options) log.Logger {

	cores := []zapcore.Core{}
	var writers []io.Writer

	if options.Console.Enabled {
		level := logLevel(options.Console.Level)
		writer := zapcore.Lock(os.Stdout)
		coreconsole := zapcore.NewCore(getEncoder(options.Console.Formatter), writer, level)
		cores = append(cores, coreconsole)
		writers = append(writers, writer)
	}

	if options.File.Enabled {
		s := []string{options.File.Path, "/", options.File.Name}
		fileLocation := strings.Join(s, "")

		lumber := &lumberjack.Logger{
			Filename: fileLocation,
			MaxSize:  options.File.MaxSize,
			Compress: options.File.Compress,
			MaxAge:   options.File.MaxAge,
		}

		level := logLevel(options.File.Level)
		writer := zapcore.AddSync(lumber)
		corefile := zapcore.NewCore(getEncoder(options.File.Formatter), writer, level)
		cores = append(cores, corefile)
		writers = append(writers, lumber)
	}

	combinedCore := zapcore.NewTee(cores...)

	// AddCallerSkip skips 2 number of callers, this is important else the file that gets
	// logged will always be the wrapped file. In our case zap.go
	zaplogger := newSugaredLogger(combinedCore)

	// Default options are only applied if this is called via NewLogger
	// If called direct, the options passed to this function may be empty.
	// Hence the default is reinforced here.
	errorField := options.ErrorFieldName
	if errorField == "" {
		errorField = defaultErrorFieldName
	}

	newlogger := &zapLogger{
		fields:         log.Fields{},
		sugaredLogger:  zaplogger,
		writers:        writers,
		core:           combinedCore,
		errorFieldName: errorField,
	}

	log.Set(newlogger)

	return newlogger
}

func defaultOptions() *Options {
	return &Options{
		ErrorFieldName: defaultErrorFieldName,

		Console: struct {
			Enabled   bool
			Level     string
			Formatter string
		}{
			Enabled:   defaultConsoleEnabled,
			Level:     defaultConsoleLevel,
			Formatter: defaultConsoleFormatter,
		},
		File: struct {
			Enabled   bool
			Level     string
			Path      string
			Name      string
			MaxSize   int
			Compress  bool
			MaxAge    int
			Formatter string
		}{
			Enabled:   defaultFileEnabled,
			Level:     defaultFileLevel,
			Path:      defaultFilePath,
			Name:      defaultFileName,
			MaxSize:   defaultFileMaxSize,
			Compress:  defaultFileCompress,
			MaxAge:    defaultFileMaxAge,
			Formatter: defaultFileFormatter,
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

func newSugaredLogger(core zapcore.Core) *zap.SugaredLogger {
	return zap.New(core,
		zap.AddCallerSkip(2),
		zap.AddCaller(),
	).Sugar()
}

func getEncoder(format string) zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	switch format {
	case "JSON":
		return zapcore.NewJSONEncoder(encoderConfig)
	default:
		return zapcore.NewConsoleEncoder(encoderConfig)
	}
}

func logLevel(level string) zapcore.Level {
	switch level {
	case "TRACE":
		return zapcore.DebugLevel
	case "WARN":
		return zapcore.WarnLevel
	case "DEBUG":
		return zapcore.DebugLevel
	case "ERROR":
		return zapcore.ErrorLevel
	case "PANIC":
		return zapcore.PanicLevel
	case "FATAL":
		return zapcore.FatalLevel
	default:
		return zapcore.InfoLevel
	}
}

type zapLogger struct {
	sugaredLogger  *zap.SugaredLogger
	fields         log.Fields
	writers        []io.Writer
	core           zapcore.Core
	errorFieldName string
}

// Printf uses (*zap.SugaredLogger).Infof to log a templated message.
func (l *zapLogger) Printf(format string, args ...interface{}) {
	l.sugaredLogger.Infof(format, args...)
}

// Tracef uses (*zap.SugaredLogger).Debugf to log a templated message.
func (l *zapLogger) Tracef(format string, args ...interface{}) {
	l.sugaredLogger.Debugf(format, args...)
}

// Trace uses (*zap.SugaredLogger).Debug to log a message.
func (l *zapLogger) Trace(args ...interface{}) {
	l.sugaredLogger.Debug(args...)
}

// Debug uses (*zap.SugaredLogger).Debug to log a message.
func (l *zapLogger) Debug(args ...interface{}) {
	l.sugaredLogger.Debug(args...)
}

// Info uses (*zap.SugaredLogger).Info to log a message.
func (l *zapLogger) Info(args ...interface{}) {
	l.sugaredLogger.Info(args...)
}

// Warn uses (*zap.SugaredLogger).Warn to log a message.
func (l *zapLogger) Warn(args ...interface{}) {
	l.sugaredLogger.Warn(args...)
}

// Error uses (*zap.SugaredLogger).Error to log a message.
func (l *zapLogger) Error(args ...interface{}) {
	l.sugaredLogger.Error(args...)
}

// Fatal uses (*zap.SugaredLogger).Fatal to log a message and call os.Exit(1).
func (l *zapLogger) Fatal(args ...interface{}) {
	l.sugaredLogger.Fatal(args...)
}

// Panic uses (*zap.SugaredLogger).Panic to log a message and panic.
func (l *zapLogger) Panic(args ...interface{}) {
	l.sugaredLogger.Panic(args...)
}

// WithField constructs a new Logger with l.fields and provided key and value field.
func (l *zapLogger) WithField(key string, value interface{}) log.Logger {
	newFields := log.Fields{}
	for k, v := range l.fields {
		newFields[k] = v
	}

	newFields[key] = value

	f := mapToSlice(newFields)
	newLogger := newSugaredLogger(l.core).With(f...)
	return &zapLogger{newLogger, newFields, l.writers, l.core, l.errorFieldName}
}

// Output returns a Writer that represents the zap writers.
func (l *zapLogger) Output() io.Writer {
	return io.MultiWriter(l.writers...)
}

// Debugf uses (*zap.SugaredLogger).Debugf to log a templated message.
func (l *zapLogger) Debugf(format string, args ...interface{}) {
	l.sugaredLogger.Debugf(format, args...)
}

// Infof uses (*zap.SugaredLogger).Infof to log a templated message.
func (l *zapLogger) Infof(format string, args ...interface{}) {
	l.sugaredLogger.Infof(format, args...)
}

// Warnf uses (*zap.SugaredLogger).Warnf to log a templated message.
func (l *zapLogger) Warnf(format string, args ...interface{}) {
	l.sugaredLogger.Warnf(format, args...)
}

// Errorf uses (*zap.SugaredLogger).Errorf to log a templated message.
func (l *zapLogger) Errorf(format string, args ...interface{}) {
	l.sugaredLogger.Errorf(format, args...)
}

// Fatalf uses (*zap.SugaredLogger).Fatalf to log a templated message and call os.Exit(1).
func (l *zapLogger) Fatalf(format string, args ...interface{}) {
	l.sugaredLogger.Fatalf(format, args...)
}

// Panicf uses (*zap.SugaredLogger).Panif to log a templated message and panic.
func (l *zapLogger) Panicf(format string, args ...interface{}) {
	l.sugaredLogger.Panicf(format, args...)
}

// WithFields constructs a new Logger with l.fields and the provided fields.
func (l *zapLogger) WithFields(fields map[string]interface{}) log.Logger {
	newFields := log.Fields{}

	for k, v := range l.fields {
		newFields[k] = v
	}

	for k, v := range fields {
		newFields[k] = v
	}

	f := mapToSlice(newFields)
	newLogger := newSugaredLogger(l.core).With(f...)
	return &zapLogger{newLogger, newFields, l.writers, l.core, l.errorFieldName}
}

// WithTypeOf adds type and package information fields.
func (l *zapLogger) WithTypeOf(obj interface{}) log.Logger {

	t := reflect.TypeOf(obj)

	return l.WithFields(log.Fields{
		"reflect.type.name":    t.Name(),
		"reflect.type.package": t.PkgPath(),
	})
}

func (l *zapLogger) WithError(err error) log.Logger {
	return l.WithField(l.errorFieldName, err.Error())
}

func (l *zapLogger) Fields() log.Fields {
	return l.fields
}

// ToContext returns a copy of ctx in which its fields are added to those of l.
func (l *zapLogger) ToContext(ctx context.Context) context.Context {
	fields := l.Fields()

	ctxFields := fieldsFromContext(ctx)

	for k, v := range fields {
		ctxFields[k] = v
	}

	return context.WithValue(ctx, key, ctxFields)
}

// FromContext returns a Logger from ctx.
func (l *zapLogger) FromContext(ctx context.Context) log.Logger {
	fields := fieldsFromContext(ctx)
	return l.WithFields(fields)
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

func mapToSlice(m log.Fields) []interface{} {
	f := make([]interface{}, 2*len(m))
	i := 0
	for k, v := range m {
		f[i] = k
		f[i+1] = v
		i = i + 2
	}

	return f
}
