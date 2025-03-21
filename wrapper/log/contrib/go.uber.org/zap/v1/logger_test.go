package zap

import (
	// "github.com/stretchr/testify/mock"

	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"reflect"
	"strings"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/xgodev/boost/wrapper/log"
	"go.uber.org/zap/zapcore"
)

type LoggerSuite struct {
	suite.Suite
}

func TestLoggerSuite(t *testing.T) {
	suite.Run(t, new(LoggerSuite))
}

func (s *LoggerSuite) TestNewLogger() {

	tt := []struct {
		name string
		want func() log.Logger
		opts []Option
	}{
		{
			name: "New Logger with default options",
			want: func() log.Logger {
				return NewLoggerWithOptions(defaultOptions())
			},
			opts: []Option{},
		},
		{
			name: "New Logger with console enabled",
			want: func() log.Logger {
				opts := defaultOptions()
				opts.Console.Enabled = true
				return NewLoggerWithOptions(opts)
			},
			opts: []Option{
				WithConsoleEnabled(true),
			},
		},
		{
			name: "New Logger with console and file enabled",
			want: func() log.Logger {
				opts := defaultOptions()
				opts.Console.Enabled = true
				opts.File.Enabled = true
				return NewLoggerWithOptions(opts)
			},
			opts: []Option{
				WithConsoleEnabled(true),
				WithFileEnabled(true),
			},
		},
		{
			name: "New Logger with console disabled and file enabled",
			want: func() log.Logger {
				opts := defaultOptions()
				opts.Console.Enabled = false
				opts.File.Enabled = true
				return NewLoggerWithOptions(opts)
			},
			opts: []Option{
				WithConsoleEnabled(false),
				WithFileEnabled(true),
			},
		},
		{
			name: "New Logger with custom error field",
			want: func() log.Logger {
				opts := defaultOptions()
				opts.ErrorFieldName = "error"
				return NewLoggerWithOptions(opts)
			},
			opts: []Option{
				WithErrorFieldName("error"),
			},
		},
	}

	for _, t := range tt {
		s.Run(t.name, func() {
			got := NewLogger(t.opts...)
			want := t.want()

			s.Assert().True(loggersAreEqual(got, want), "got  %v\nwant %v", got, want)
		})
	}
}

func loggersAreEqual(got log.Logger, want log.Logger) bool {
	g := got.(*zapLogger)
	w := want.(*zapLogger)
	return reflect.DeepEqual(g.fields, w.fields) && g.errorFieldName == w.errorFieldName && reflect.DeepEqual(g.writers, w.writers)
}

func (s *LoggerSuite) Test_getZapLevel() {
	tt := []struct {
		name  string
		level string
		want  zapcore.Level
	}{
		{
			name:  "log level TRACE (zap -> DEBUG)",
			level: "TRACE",
			want:  zapcore.DebugLevel,
		},
		{
			name:  "log level DEBUG",
			level: "DEBUG",
			want:  zapcore.DebugLevel,
		},
		{
			name:  "log level INFO",
			level: "INFO",
			want:  zapcore.InfoLevel,
		},
		{
			name:  "log level ERROR",
			level: "ERROR",
			want:  zapcore.ErrorLevel,
		},
		{
			name:  "log level WARN",
			level: "WARN",
			want:  zapcore.WarnLevel,
		},
		{
			name:  "log level PANIC",
			level: "PANIC",
			want:  zapcore.PanicLevel,
		},
		{
			name:  "log level FATAL",
			level: "FATAL",
			want:  zapcore.FatalLevel,
		},
	}

	for _, t := range tt {
		s.Run(t.name, func() {
			got := logLevel(t.level)
			s.Assert().True(reflect.DeepEqual(got, t.want), "got  %v\nwant %v", got, t.want)
		})
	}
}

func initLogCapture() (logger log.Logger, w *os.File, r *os.File) {
	original := os.Stdout
	defer func() { os.Stdout = original }()
	r, w, _ = os.Pipe()
	os.Stdout = w
	logger = NewLogger(WithConsoleLevel("TRACE"))
	return logger, w, r
}

func captureLog(w *os.File, r *os.File) string {
	w.Close()
	b, _ := io.ReadAll(r)
	r2, w2, _ := os.Pipe()
	*r = *r2
	*w = *w2
	return string(b)
}

func (s *LoggerSuite) TestLogger() {
	logger, w, r := initLogCapture()
	tt := []struct {
		name   string
		method string
		want   string
	}{
		{
			name:   "logger Printf method",
			method: "Printf",
			want:   "info\treflect/value.go:368\tBlah",
		},
		{
			name:   "logger Trace method",
			method: "Trace",
			want:   "debug\treflect/value.go:368\tBlah",
		},
		{
			name:   "logger Tracef method",
			method: "Tracef",
			want:   "debug\treflect/value.go:368\tBlah",
		},
		{
			name:   "logger Debug method",
			method: "Debug",
			want:   "debug\treflect/value.go:368\tBlah",
		},
		{
			name:   "logger Debugf method",
			method: "Debugf",
			want:   "debug\treflect/value.go:368\tBlah",
		},
		{
			name:   "logger Info method",
			method: "Info",
			want:   "info\treflect/value.go:368\tBlah",
		},
		{
			name:   "logger Infof method",
			method: "Infof",
			want:   "info\treflect/value.go:368\tBlah",
		},
		{
			name:   "logger Warn method",
			method: "Warn",
			want:   "warn\treflect/value.go:368\tBlah",
		},
		{
			name:   "logger Warnf method",
			method: "Warnf",
			want:   "warn\treflect/value.go:368\tBlah",
		},
		{
			name:   "logger Error method",
			method: "Error",
			want:   "error\treflect/value.go:368\tBlah",
		},
		{
			name:   "logger Errorf method",
			method: "Errorf",
			want:   "error\treflect/value.go:368\tBlah",
		},
		{
			name:   "logger Panic method",
			method: "Panic",
			want:   "panic\treflect/value.go:368\tBlah",
		},
		{
			name:   "logger Panicf method",
			method: "Panicf",
			want:   "panic\treflect/value.go:368\tBlah",
		},
	}
	for _, t := range tt {
		s.Run(t.name, func() {
			defer func() {
				recover() // for panic case
				got := captureLog(w, r)
				s.Assert().True(strings.Contains(got, t.want), "got %v\nmust contain %v", got, t.want)
			}()
			m := reflect.ValueOf(logger).MethodByName(t.method)
			m.Call([]reflect.Value{reflect.ValueOf("Blah")})
		})
	}
}

func (s *LoggerSuite) TestLoggerFatal() {
	tt := []struct {
		name   string
		method string
		want   string
	}{
		{
			name:   "logger Fatal method",
			method: "Fatal",
			want:   "fatal\treflect/value.go:368\tBlah",
		},
		{
			name:   "logger Fatalf method",
			method: "Fatalf",
			want:   "fatal\treflect/value.go:368\tBlah",
		},
	}
	for _, t := range tt {
		s.Run(t.name, func() {
			cmd := exec.Command(os.Args[0], "-test.run=TestLoggerMethod")
			cmd.Env = append(os.Environ(), fmt.Sprintf("LOGGER_TEST_METHOD=%s", t.method))
			out, e := cmd.CombinedOutput()
			s.Assert().False(false, "got an unexpected error = %v", e)
			got := string(out)
			s.Assert().True(strings.Contains(got, t.want), "got %v\nmust contain %v", got, t.want)
		})
	}
}

func TestLoggerMethod(t *testing.T) {
	method := os.Getenv("LOGGER_TEST_METHOD")
	if method == "" {
		return
	}
	logger := NewLogger(WithConsoleLevel("TRACE"))
	m := reflect.ValueOf(logger).MethodByName(method)
	m.Call([]reflect.Value{reflect.ValueOf("Blah")})
}

func buildLogger() *zapLogger {
	level := logLevel("TRACE")
	writer := zapcore.Lock(os.Stdout)
	coreconsole := zapcore.NewCore(getEncoder("TEXT"), writer, level)

	core := zapcore.NewTee(coreconsole)
	zaplogger := newSugaredLogger(core)
	return &zapLogger{
		fields:         log.Fields{},
		sugaredLogger:  zaplogger,
		writers:        []io.Writer{writer},
		core:           core,
		errorFieldName: "err",
	}
}

func (s *LoggerSuite) TestLoggerWithMethods() {
	l := buildLogger()
	tt := []struct {
		name   string
		method func() log.Logger
		want   func() log.Logger
	}{
		{
			name: "logger WithField",
			method: func() log.Logger {
				return l.WithField("ID", "1")
			},
			want: func() log.Logger {
				return &zapLogger{l.sugaredLogger.With("ID", "1"), log.Fields{"ID": "1"}, l.writers, l.core, l.errorFieldName}
			},
		},
		{
			name: "logger WithFields",
			method: func() log.Logger {
				return l.WithFields(log.Fields{
					"ID":   "12",
					"Name": "Stockton",
				})
			},
			want: func() log.Logger {
				return &zapLogger{l.sugaredLogger.With("ID", "12", "Name", "Stockton"), log.Fields{
					"ID":   "12",
					"Name": "Stockton",
				}, l.writers, l.core, l.errorFieldName}
			},
		},
		{
			name: "logger WithTypeOf",
			method: func() log.Logger {
				return l.WithTypeOf(l)
			},
			want: func() log.Logger {
				t := reflect.TypeOf(l)
				l2 := &zapLogger{
					l.sugaredLogger.With(
						"reflect.type.name", t.Name(),
						"reflect.type.package", t.PkgPath(),
					),
					log.Fields{
						"reflect.type.name":    t.Name(),
						"reflect.type.package": t.PkgPath(),
					},
					l.writers,
					l.core,
					l.errorFieldName,
				}
				return l2
			},
		},
		{
			name: "logger WithError",
			method: func() log.Logger {
				return l.WithError(errors.New("something bad"))
			},
			want: func() log.Logger {
				return &zapLogger{l.sugaredLogger.With("err", "something bad"), log.Fields{
					"err": "something bad",
				}, l.writers, l.core, l.errorFieldName}
			},
		},
	}

	for _, t := range tt {
		s.Run(t.name, func() {
			got := t.method()
			want := t.want()
			s.Assert().True(loggersAreEqual(got, want), "got  %v\nwant %v", got, want)
		})
	}
}

func (s *LoggerSuite) TestLoggerFields() {
	logger := NewLogger()
	tt := []struct {
		name string
		want log.Fields
	}{
		{
			name: "get logger fields",
			want: log.Fields{},
		},
	}
	for _, t := range tt {
		s.Run(t.name, func() {
			got := logger.Fields()
			s.Assert().True(reflect.DeepEqual(got, t.want), "got  %v\nwant %v", got, t.want)
		})
	}
}

func (s *LoggerSuite) TestLoggerOutput() {
	logger := NewLogger()
	tt := []struct {
		name string
		want io.Writer
	}{
		{
			name: "get logger output",
			want: io.MultiWriter(zapcore.Lock(os.Stdout)),
		},
	}
	for _, t := range tt {
		s.Run(t.name, func() {
			got := logger.Output()
			s.Assert().True(reflect.DeepEqual(got, t.want), "got  %v\nwant %v", got, t.want)
		})
	}
}

func (s *LoggerSuite) TestLoggerToContext() {
	logger := NewLogger().WithField("ID", "1")
	tt := []struct {
		name string
		want log.Fields
	}{
		{
			name: "set logger to context",
			want: log.Fields{"ID": "1"},
		},
	}
	for _, t := range tt {
		s.Run(t.name, func() {
			ctx := logger.ToContext(context.Background())
			got := ctx.Value(key)
			s.Assert().True(reflect.DeepEqual(got, t.want), "got  %v\nwant %v", got, t.want)
		})
	}
}

func (s *LoggerSuite) TestLoggerFromContext() {
	l := NewLogger()
	tt := []struct {
		name string
		want log.Logger
	}{
		{
			name: "get logger from context",
			want: l.WithFields(log.Fields{}),
		},
	}
	for _, t := range tt {
		s.Run(t.name, func() {
			got := l.FromContext(context.Background())
			s.Assert().True(loggersAreEqual(got, t.want), "got  %v\nwant %v", got, t.want)
		})
	}
}

func (s *LoggerSuite) Test_fieldsFromContext() {
	fields := log.Fields{
		"ID":   "12",
		"Name": "Stockton",
	}
	tt := []struct {
		name string
		in   context.Context
		want log.Fields
	}{
		{
			name: "when fields are not previously present on context",
			in:   context.Background(),
			want: log.Fields{},
		},
		{
			name: "when fields are present on context",
			in:   context.WithValue(context.Background(), key, fields),
			want: fields,
		},
		{
			name: "when context is nil",
			in:   nil,
			want: log.Fields{},
		},
	}
	for _, t := range tt {
		s.Run(t.name, func() {
			got := fieldsFromContext(t.in)
			s.Assert().True(reflect.DeepEqual(got, t.want), "got  %v\nwant %v", got, t.want)
		})
	}
}

func (s *LoggerSuite) Test_getEncoder() {

	tt := []struct {
		name string
		in   string
		want string
	}{
		{
			name: "when JSON",
			in:   "JSON",
			want: "*zapcore.jsonEncoder",
		},
		{
			name: "when default",
			in:   "TEXT",
			want: "zapcore.consoleEncoder",
		},
	}
	for _, t := range tt {
		s.Run(t.name, func() {
			got := reflect.TypeOf(getEncoder(t.in)).String()
			s.Assert().True(got == t.want, "got  %v\nwant %v", got, t.want)
		})
	}
}
