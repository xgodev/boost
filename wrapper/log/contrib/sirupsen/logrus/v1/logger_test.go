package logrus

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/syslog"
	"os"
	"os/exec"
	"reflect"
	"strings"
	"testing"

	"github.com/sirupsen/logrus"
	logrus_syslog "github.com/sirupsen/logrus/hooks/syslog"
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
			name: "New Logger with hook",
			want: func() log.Logger {
				opts := defaultOptions()
				opts.Hooks = []logrus.Hook{
					&logrus_syslog.SyslogHook{
						Writer: &syslog.Writer{},
					},
				}
				opts.Console.Enabled = false
				opts.File.Enabled = true
				return NewLoggerWithOptions(opts)
			},
			opts: []Option{
				WithHook(&logrus_syslog.SyslogHook{
					Writer: &syslog.Writer{},
				}),
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
			s.Assert().True(reflect.DeepEqual(got, t.want()), "got  %v\nwant %v", got, t.want)
		})
	}
}

func (s *LoggerSuite) Test_logLevel() {
	tt := []struct {
		name  string
		level string
		want  logrus.Level
	}{
		{
			name:  "log level TRACE",
			level: "TRACE",
			want:  logrus.TraceLevel,
		},
		{
			name:  "log level DEBUG",
			level: "DEBUG",
			want:  logrus.DebugLevel,
		},
		{
			name:  "log level INFO",
			level: "INFO",
			want:  logrus.InfoLevel,
		},
		{
			name:  "log level ERROR",
			level: "ERROR",
			want:  logrus.ErrorLevel,
		},
		{
			name:  "log level WARN",
			level: "WARN",
			want:  logrus.WarnLevel,
		},
		{
			name:  "log level PANIC",
			level: "PANIC",
			want:  logrus.PanicLevel,
		},
		{
			name:  "log level FATAL",
			level: "FATAL",
			want:  logrus.FatalLevel,
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
			want:   "level=info msg=Blah",
		},
		{
			name:   "logger Trace method",
			method: "Trace",
			want:   "level=trace msg=Blah",
		},
		{
			name:   "logger Tracef method",
			method: "Tracef",
			want:   "level=trace msg=Blah",
		},
		{
			name:   "logger Debug method",
			method: "Debug",
			want:   "level=debug msg=Blah",
		},
		{
			name:   "logger Debugf method",
			method: "Debugf",
			want:   "level=debug msg=Blah",
		},
		{
			name:   "logger Info method",
			method: "Info",
			want:   "level=info msg=Blah",
		},
		{
			name:   "logger Infof method",
			method: "Infof",
			want:   "level=info msg=Blah",
		},
		{
			name:   "logger Warn method",
			method: "Warn",
			want:   "level=warning msg=Blah",
		},
		{
			name:   "logger Warnf method",
			method: "Warnf",
			want:   "level=warning msg=Blah",
		},
		{
			name:   "logger Error method",
			method: "Error",
			want:   "level=error msg=Blah",
		},
		{
			name:   "logger Errorf method",
			method: "Errorf",
			want:   "level=error msg=Blah",
		},
		{
			name:   "logger Panic method",
			method: "Panic",
			want:   "level=panic msg=Blah",
		},
		{
			name:   "logger Panicf method",
			method: "Panicf",
			want:   "level=panic msg=Blah",
		},
	}
	for _, t := range tt {
		s.Run(t.name, func() {
			defer func() {
				recover() //for panic case
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
			want:   "level=fatal msg=Blah",
		},
		{
			name:   "logger Fatalf method",
			method: "Fatalf",
			want:   "level=fatal msg=Blah",
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
	fmt.Println("Blah")
	if method == "" {
		return
	}
	logger := NewLogger(WithConsoleLevel("TRACE"))
	m := reflect.ValueOf(logger).MethodByName(method)
	m.Call([]reflect.Value{reflect.ValueOf("Blah")})
}

func (s *LoggerSuite) TestLoggerWithMethods() {
	l := &logger{
		logger:         logrus.New(),
		fields:         log.Fields{},
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
				l2 := &logEntry{
					entry: l.logger.WithField("ID", "1"),
					fields: log.Fields{
						"ID": "1",
					}}
				return l2
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
				l2 := &logEntry{
					entry: l.logger.WithFields(logrus.Fields{
						"ID":   "12",
						"Name": "Stockton",
					}),
					fields: log.Fields{
						"ID":   "12",
						"Name": "Stockton",
					}}
				return l2
			},
		},
		{
			name: "logger WithTypeOf",
			method: func() log.Logger {
				return l.WithTypeOf(l)
			},
			want: func() log.Logger {
				t := reflect.TypeOf(l)
				l2 := &logEntry{
					entry: l.logger.WithFields(logrus.Fields{
						"reflect.type.name":    t.Name(),
						"reflect.type.package": t.PkgPath(),
					}),
					fields: log.Fields{
						"reflect.type.name":    t.Name(),
						"reflect.type.package": t.PkgPath(),
					}}
				return l2
			},
		},
		{
			name: "logger WithError",
			method: func() log.Logger {
				return l.WithError(errors.New("something bad"))
			},
			want: func() log.Logger {
				l2 := &logEntry{
					entry: l.logger.WithFields(logrus.Fields{
						"err": "something bad",
					}),
					fields: log.Fields{
						"err": "something bad",
					}}
				return l2
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
		want *os.File
	}{
		{
			name: "get logger output",
			want: os.Stdout,
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
	logger := NewLogger()
	tt := []struct {
		name string
		want log.Fields
	}{
		{
			name: "set logger to context",
			want: log.Fields{},
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
			s.Assert().True(reflect.DeepEqual(got, t.want), "got  %v\nwant %v", got, t.want)
		})
	}
}

func (s *LoggerSuite) TestLoggerEntry() {
	logger, w, r := initLogCapture()
	logger = logger.WithField("ID", "1")
	tt := []struct {
		name   string
		method string
		want   string
	}{
		{
			name:   "logger entry Printf method",
			method: "Printf",
			want:   "level=info msg=Blah ID=1",
		},
		{
			name:   "logger entry Trace method",
			method: "Trace",
			want:   "level=trace msg=Blah ID=1",
		},
		{
			name:   "logger entry Tracef method",
			method: "Tracef",
			want:   "level=trace msg=Blah ID=1",
		},
		{
			name:   "logger entry Debug method",
			method: "Debug",
			want:   "level=debug msg=Blah ID=1",
		},
		{
			name:   "logger entry Debugf method",
			method: "Debugf",
			want:   "level=debug msg=Blah ID=1",
		},
		{
			name:   "logger entry Info method",
			method: "Info",
			want:   "level=info msg=Blah ID=1",
		},
		{
			name:   "logger entry Infof method",
			method: "Infof",
			want:   "level=info msg=Blah ID=1",
		},
		{
			name:   "logger entry Warn method",
			method: "Warn",
			want:   "level=warning msg=Blah ID=1",
		},
		{
			name:   "logger entry Warnf method",
			method: "Warnf",
			want:   "level=warning msg=Blah ID=1",
		},
		{
			name:   "logger entry Error method",
			method: "Error",
			want:   "level=error msg=Blah ID=1",
		},
		{
			name:   "logger entry Errorf method",
			method: "Errorf",
			want:   "level=error msg=Blah ID=1",
		},
		{
			name:   "logger entry Panic method",
			method: "Panic",
			want:   "level=panic msg=Blah ID=1",
		},
		{
			name:   "logger entry Panicf method",
			method: "Panicf",
			want:   "level=panic msg=Blah ID=1",
		},
	}
	for _, t := range tt {
		s.Run(t.name, func() {
			defer func() {
				recover() //for panic case
				got := captureLog(w, r)
				s.Assert().True(strings.Contains(got, t.want), "got %v\nmust contain %v", got, t.want)
			}()
			m := reflect.ValueOf(logger).MethodByName(t.method)
			m.Call([]reflect.Value{reflect.ValueOf("Blah")})
		})
	}
}

func (s *LoggerSuite) TestLoggerEntryFatal() {
	tt := []struct {
		name   string
		method string
		want   string
	}{
		{
			name:   "logger entry Fatal method",
			method: "Fatal",
			want:   "level=fatal msg=Blah ID=1",
		},
		{
			name:   "logger entry Fatalf method",
			method: "Fatalf",
			want:   "level=fatal msg=Blah ID=1",
		},
	}
	for _, t := range tt {
		s.Run(t.name, func() {
			cmd := exec.Command(os.Args[0], "-test.run=TestLoggerEntryMethod")
			cmd.Env = append(os.Environ(), fmt.Sprintf("LOGGER_ENTRY_TEST_METHOD=%s", t.method))
			out, e := cmd.CombinedOutput()
			s.Assert().False(false, "got an unexpected error = %v", e)
			got := string(out)
			s.Assert().True(strings.Contains(got, t.want), "got %v\nmust contain %v", got, t.want)
		})
	}
}

func TestLoggerEntryMethod(t *testing.T) {
	method := os.Getenv("LOGGER_ENTRY_TEST_METHOD")
	if method == "" {
		return
	}
	logger := NewLogger(WithConsoleLevel("TRACE")).WithField("ID", "1")
	m := reflect.ValueOf(logger).MethodByName(method)
	m.Call([]reflect.Value{reflect.ValueOf("Blah")})
}

func (s *LoggerSuite) TestLoggerEntryWithMethods() {
	l := &logger{
		logger: logrus.New(),
		fields: log.Fields{},
	}
	le := l.WithField("ID", "12")
	tt := []struct {
		name   string
		method func() log.Logger
		want   func() log.Logger
	}{
		{
			name: "logger WithField",
			method: func() log.Logger {
				return le.WithField("Name", "Stockton")
			},
			want: func() log.Logger {
				l2 := &logEntry{
					entry: l.logger.WithFields(logrus.Fields{
						"ID":   "12",
						"Name": "Stockton",
					}),
					fields: log.Fields{
						"ID":   "12",
						"Name": "Stockton",
					}}
				return l2
			},
		},
		{
			name: "logger WithFields",
			method: func() log.Logger {
				return le.WithFields(log.Fields{
					"Name":     "Stockton",
					"Position": "Point Guard",
				})
			},
			want: func() log.Logger {
				l2 := &logEntry{
					entry: l.logger.WithFields(logrus.Fields{
						"ID":       "12",
						"Name":     "Stockton",
						"Position": "Point Guard",
					}),
					fields: log.Fields{
						"ID":       "12",
						"Name":     "Stockton",
						"Position": "Point Guard",
					}}
				return l2
			},
		},
		{
			name: "logger WithTypeOf",
			method: func() log.Logger {
				return le.WithTypeOf(l)
			},
			want: func() log.Logger {
				t := reflect.TypeOf(l)
				l2 := &logEntry{
					entry: l.logger.WithFields(logrus.Fields{
						"ID":                   "12",
						"reflect.type.name":    t.Name(),
						"reflect.type.package": t.PkgPath(),
					}),
					fields: log.Fields{
						"ID":                   "12",
						"reflect.type.name":    t.Name(),
						"reflect.type.package": t.PkgPath(),
					}}
				return l2
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

func (s *LoggerSuite) TestLoggerEntryFields() {
	fields := log.Fields{
		"ID":   "12",
		"Name": "Stockton",
	}
	logger := NewLogger().WithFields(fields)
	tt := []struct {
		name string
		want log.Fields
	}{
		{
			name: "get logger entry fields",
			want: fields,
		},
	}
	for _, t := range tt {
		s.Run(t.name, func() {
			got := logger.Fields()
			s.Assert().True(reflect.DeepEqual(got, t.want), "got  %v\nwant %v", got, t.want)
		})
	}
}

func (s *LoggerSuite) TestLoggerEntryOutput() {
	logger := NewLogger().WithFields(log.Fields{})
	tt := []struct {
		name string
		want *os.File
	}{
		{
			name: "get logger entry output",
			want: os.Stdout,
		},
	}
	for _, t := range tt {
		s.Run(t.name, func() {
			got := logger.Output()
			s.Assert().True(reflect.DeepEqual(got, t.want), "got  %v\nwant %v", got, t.want)
		})
	}
}

func (s *LoggerSuite) TestLoggerEntryToContext() {
	logger := NewLogger().WithFields(log.Fields{})
	tt := []struct {
		name string
		want log.Fields
	}{
		{
			name: "set logger to context",
			want: log.Fields{},
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

func (s *LoggerSuite) TestLoggerEntryFromContext() {
	l := NewLogger().WithFields(log.Fields{})
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
			s.Assert().True(reflect.DeepEqual(got, t.want), "got  %v\nwant %v", got, t.want)
		})
	}
}

func (s *LoggerSuite) Test_toContext() {
	fields := log.Fields{
		"ID":   "12",
		"Name": "Stockton",
	}
	tt := []struct {
		name string
		in   context.Context
		want func() context.Context
	}{
		{
			name: "when fields are not previously present on context",
			in:   context.Background(),
			want: func() context.Context {
				return context.WithValue(context.Background(), key, fields)
			},
		},
		{
			name: "when fields are previously present on context",
			in: context.WithValue(context.Background(), key, log.Fields{
				"Position": "Point guard",
			}),
			want: func() context.Context {
				previousFields := log.Fields{
					"Position": "Point guard",
				}
				ctx := context.WithValue(context.Background(), key, previousFields)
				newFields := log.Fields{}
				for k, v := range previousFields {
					newFields[k] = v
				}
				for k, v := range fields {
					newFields[k] = v
				}
				return context.WithValue(ctx, key, newFields)
			},
		},
	}
	for _, t := range tt {
		s.Run(t.name, func() {
			got := toContext(t.in, fields)
			want := t.want()
			s.Assert().True(reflect.DeepEqual(got, want), "got  %v\nwant %v", got, want)
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

func (s *LoggerSuite) Test_convertToLogrusFields() {

	tt := []struct {
		name string
		in   log.Fields
		want logrus.Fields
	}{
		{
			name: "success converting to logrus fields",
			in: log.Fields{
				"ID":   "12",
				"Name": "Stockton",
			},
			want: logrus.Fields{
				"ID":   "12",
				"Name": "Stockton",
			},
		},
	}
	for _, t := range tt {
		s.Run(t.name, func() {
			got := convertToLogrusFields(t.in)
			s.Assert().True(reflect.DeepEqual(got, t.want), "got  %v\nwant %v", got, t.want)
		})
	}
}

func (s *LoggerSuite) Test_convertToFields() {

	tt := []struct {
		name string
		in   logrus.Fields
		want log.Fields
	}{
		{
			name: "success converting to fields",
			in: logrus.Fields{
				"ID":   "12",
				"Name": "Stockton",
			},
			want: log.Fields{
				"ID":   "12",
				"Name": "Stockton",
			},
		},
	}
	for _, t := range tt {
		s.Run(t.name, func() {
			got := convertToFields(t.in)
			s.Assert().True(reflect.DeepEqual(got, t.want), "got  %v\nwant %v", got, t.want)
		})
	}
}
