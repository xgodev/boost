package zerolog

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

	"github.com/rs/zerolog"
	"github.com/xgodev/boost/wrapper/log"

	"github.com/stretchr/testify/suite"
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
			name: "New Logger with console disabled and file disabled",
			want: func() log.Logger {
				opts := defaultOptions()
				opts.Console.Enabled = false
				opts.File.Enabled = false
				return NewLoggerWithOptions(opts)
			},
			opts: []Option{
				WithConsoleEnabled(false),
				WithFileEnabled(false),
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
			s.Assert().True(reflect.DeepEqual(got, t.want()), "got  %v\nwant %v", got, t.want)
		})
	}
}

func (s *LoggerSuite) Test_logLevel() {
	tt := []struct {
		name  string
		level string
		want  zerolog.Level
	}{
		{
			name:  "log level TRACE",
			level: "TRACE",
			want:  zerolog.TraceLevel,
		},
		{
			name:  "log level DEBUG",
			level: "DEBUG",
			want:  zerolog.DebugLevel,
		},
		{
			name:  "log level INFO",
			level: "INFO",
			want:  zerolog.InfoLevel,
		},
		{
			name:  "log level ERROR",
			level: "ERROR",
			want:  zerolog.ErrorLevel,
		},
		{
			name:  "log level WARN",
			level: "WARN",
			want:  zerolog.WarnLevel,
		},
		{
			name:  "log level PANIC",
			level: "PANIC",
			want:  zerolog.PanicLevel,
		},
		{
			name:  "log level FATAL",
			level: "FATAL",
			want:  zerolog.FatalLevel,
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
	logger = NewLogger(WithLevel("TRACE"))
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
		want   []string
	}{
		{
			name:   "logger Printf method",
			method: "Printf",
			want:   []string{"DBG", "Blah"},
		},
		{
			name:   "logger Trace method",
			method: "Trace",
			want:   []string{"TRC", "Blah"},
		},
		{
			name:   "logger Tracef method",
			method: "Tracef",
			want:   []string{"TRC", "Blah"},
		},
		{
			name:   "logger Debug method",
			method: "Debug",
			want:   []string{"DBG", "Blah"},
		},
		{
			name:   "logger Debugf method",
			method: "Debugf",
			want:   []string{"DBG", "Blah"},
		},
		{
			name:   "logger Info method",
			method: "Info",
			want:   []string{"INF", "Blah"},
		},
		{
			name:   "logger Infof method",
			method: "Infof",
			want:   []string{"INF", "Blah"},
		},
		{
			name:   "logger Warn method",
			method: "Warn",
			want:   []string{"WRN", "Blah"},
		},
		{
			name:   "logger Warnf method",
			method: "Warnf",
			want:   []string{"WRN", "Blah"},
		},
		{
			name:   "logger Error method",
			method: "Error",
			want:   []string{"ERR", "Blah"},
		},
		{
			name:   "logger Errorf method",
			method: "Errorf",
			want:   []string{"ERR", "Blah"},
		},
		{
			name:   "logger Panic method",
			method: "Panic",
			want:   []string{"PNC", "Blah"},
		},
		{
			name:   "logger Panicf method",
			method: "Panicf",
			want:   []string{"PNC", "Blah"},
		},
	}
	for _, t := range tt {
		s.Run(t.name, func() {
			defer func() {
				recover() // for panic case
				got := captureLog(w, r)
				for _, w := range t.want {
					s.Assert().True(strings.Contains(got, w), "got %v\nmust contain %v", got, w)
				}
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
		want   []string
	}{
		{
			name:   "logger Fatal method",
			method: "Fatal",
			want:   []string{"FTL", "Blah"},
		},
		{
			name:   "logger Fatalf method",
			method: "Fatalf",
			want:   []string{"FTL", "Blah"},
		},
	}
	for _, t := range tt {
		s.Run(t.name, func() {
			cmd := exec.Command(os.Args[0], "-test.run=TestLoggerMethod")
			cmd.Env = append(os.Environ(), fmt.Sprintf("LOGGER_TEST_METHOD=%s", t.method))
			out, e := cmd.CombinedOutput()
			s.Assert().False(false, "got an unexpected error = %v", e)
			got := string(out)
			for _, w := range t.want {
				s.Assert().True(strings.Contains(got, w), "got %v\nmust contain %v", got, w)
			}
		})
	}
}

func TestLoggerMethod(t *testing.T) {
	method := os.Getenv("LOGGER_TEST_METHOD")
	if method == "" {
		return
	}
	logger := NewLogger(WithLevel("TRACE"))
	m := reflect.ValueOf(logger).MethodByName(method)
	m.Call([]reflect.Value{reflect.ValueOf("Blah")})
}

func (s *LoggerSuite) TestLoggerWithMethods() {
	l := &logger{
		logger:         zerolog.New(os.Stdout),
		fields:         log.Fields{},
		writer:         os.Stdout,
		errorFieldName: "err",
	}
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
				fields := map[string]interface{}{
					"ID": "1",
				}
				return &logger{
					logger:         zerolog.New(os.Stdout).With().Fields(fields).Logger(),
					fields:         fields,
					writer:         os.Stdout,
					errorFieldName: l.errorFieldName,
				}
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
				fields := map[string]interface{}{
					"ID":   "12",
					"Name": "Stockton",
				}
				return &logger{
					logger:         zerolog.New(os.Stdout).With().Fields(fields).Logger(),
					fields:         fields,
					writer:         os.Stdout,
					errorFieldName: l.errorFieldName,
				}
			},
		},
		{
			name: "logger WithTypeOf",
			method: func() log.Logger {
				return l.WithTypeOf(l)
			},
			want: func() log.Logger {
				t := reflect.TypeOf(l)
				fields := map[string]interface{}{
					"reflect.type.name":    t.Name(),
					"reflect.type.package": t.PkgPath(),
				}
				return &logger{
					logger:         zerolog.New(os.Stdout).With().Fields(fields).Logger(),
					fields:         fields,
					writer:         os.Stdout,
					errorFieldName: l.errorFieldName,
				}
			},
		},
		{
			name: "logger WithError",
			method: func() log.Logger {
				return l.WithError(errors.New("something bad"))
			},
			want: func() log.Logger {
				fields := map[string]interface{}{
					"err": "something bad",
				}
				return &logger{
					logger:         zerolog.New(os.Stdout).With().Fields(fields).Logger(),
					fields:         fields,
					writer:         os.Stdout,
					errorFieldName: l.errorFieldName,
				}
			},
		},
	}

	for _, t := range tt {
		s.Run(t.name, func() {
			got := t.method()
			want := t.want()
			s.Assert().True(reflect.DeepEqual(got, want), "got  %v\nwant %v", got, want)
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
			want: zerolog.ConsoleWriter{Out: os.Stdout},
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
		ctx  func() context.Context
		want log.Logger
	}{
		{
			name: "when logger is not in context",
			ctx:  context.Background,
			want: l.WithFields(log.Fields{}),
		},
		{
			name: "when logger is in context",
			ctx: func() context.Context {
				ctx := context.Background()
				return l.ToContext(ctx)
			},
			want: l.WithFields(log.Fields{}),
		},
	}
	for _, t := range tt {
		s.Run(t.name, func() {
			got := l.FromContext(t.ctx())
			s.Assert().True(reflect.DeepEqual(got, t.want), "got  %v\nwant %v", got, t.want)
		})
	}
}
