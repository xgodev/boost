package fx

import (
	"context"
	"errors"
	"fmt"
	"io"
	"strings"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xgodev/boost"
	"github.com/xgodev/boost/wrapper/log"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
)

// entry is one line a recorder received: its level and its text.
type entry struct {
	level string
	text  string
}

// recorder is a log.Logger that keeps every line in memory.
type recorder struct {
	mu      sync.Mutex
	entries []entry
}

func (r *recorder) add(level, msg string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.entries = append(r.entries, entry{level: level, text: msg})
}

func (r *recorder) lines() []entry {
	r.mu.Lock()
	defer r.mu.Unlock()
	return append([]entry(nil), r.entries...)
}

func (r *recorder) Printf(format string, args ...interface{}) {
	r.add("print", fmt.Sprintf(format, args...))
}
func (r *recorder) Tracef(format string, args ...interface{}) {
	r.add("trace", fmt.Sprintf(format, args...))
}
func (r *recorder) Trace(args ...interface{}) { r.add("trace", fmt.Sprint(args...)) }
func (r *recorder) Debugf(format string, args ...interface{}) {
	r.add("debug", fmt.Sprintf(format, args...))
}
func (r *recorder) Debug(args ...interface{}) { r.add("debug", fmt.Sprint(args...)) }
func (r *recorder) Infof(format string, args ...interface{}) {
	r.add("info", fmt.Sprintf(format, args...))
}
func (r *recorder) Info(args ...interface{}) { r.add("info", fmt.Sprint(args...)) }
func (r *recorder) Warnf(format string, args ...interface{}) {
	r.add("warn", fmt.Sprintf(format, args...))
}
func (r *recorder) Warn(args ...interface{}) { r.add("warn", fmt.Sprint(args...)) }
func (r *recorder) Errorf(format string, args ...interface{}) {
	r.add("error", fmt.Sprintf(format, args...))
}
func (r *recorder) Error(args ...interface{}) { r.add("error", fmt.Sprint(args...)) }
func (r *recorder) Fatalf(format string, args ...interface{}) {
	r.add("fatal", fmt.Sprintf(format, args...))
}
func (r *recorder) Fatal(args ...interface{}) { r.add("fatal", fmt.Sprint(args...)) }
func (r *recorder) Panicf(format string, args ...interface{}) {
	r.add("panic", fmt.Sprintf(format, args...))
}
func (r *recorder) Panic(args ...interface{}) { r.add("panic", fmt.Sprint(args...)) }

func (r *recorder) WithFields(map[string]interface{}) log.Logger  { return r }
func (r *recorder) WithField(string, interface{}) log.Logger      { return r }
func (r *recorder) WithError(error) log.Logger                    { return r }
func (r *recorder) WithTypeOf(interface{}) log.Logger             { return r }
func (r *recorder) ToContext(ctx context.Context) context.Context { return ctx }
func (r *recorder) FromContext(context.Context) log.Logger        { return r }
func (r *recorder) Output() io.Writer                             { return io.Discard }
func (r *recorder) Fields() log.Fields                            { return log.Fields{} }

func record(t *testing.T) *recorder {
	t.Helper()
	rec := &recorder{}
	previous := log.GetLogger()
	log.Set(rec)
	t.Cleanup(func() { log.Set(previous) })
	return rec
}

// TestNewApp_ABootErrorIsLoggedAtErrorWhateverTheFxLevel boots an app whose
// invoke fails, with the fx log level at its default: the reason the app
// does not start reaches the error level, the rest of fx's lines stay at the
// configured level.
func TestNewApp_ABootErrorIsLoggedAtErrorWhateverTheFxLevel(t *testing.T) {
	boost.Start()
	rec := record(t)

	app := NewApp(
		fx.Provide(func() string { return "provided" }),
		fx.Invoke(func(string) error { return errors.New("boom-206") }),
	)
	require.Error(t, app.Start(context.Background()))

	var reason bool
	for _, l := range rec.lines() {
		if strings.Contains(l.text, "boom-206") && l.level == "error" {
			reason = true
		}
		if strings.Contains(l.text, "PROVIDE") {
			assert.Equal(t, strings.ToLower(LogLevel()), l.level, "an fx line that is not an error stays at the fx level: %q", l.text)
		}
	}
	assert.True(t, reason, "no error-level line with the boot error: %v", rec.lines())
}

// TestEventLogger_EveryEventWithAnErrorIsLoggedAtError covers each fx event
// that can carry an error, with and without it.
func TestEventLogger_EveryEventWithAnErrorIsLoggedAtError(t *testing.T) {
	boost.Start()
	boom := errors.New("boom")
	for _, c := range []struct {
		name   string
		failed fxevent.Event
		ok     fxevent.Event
	}{
		{"OnStartExecuted", &fxevent.OnStartExecuted{FunctionName: "start", Err: boom}, &fxevent.OnStartExecuted{FunctionName: "start"}},
		{"OnStopExecuted", &fxevent.OnStopExecuted{FunctionName: "stop", Err: boom}, &fxevent.OnStopExecuted{FunctionName: "stop"}},
		{"Supplied", &fxevent.Supplied{TypeName: "string", Err: boom}, &fxevent.Supplied{TypeName: "string"}},
		{"Provided", &fxevent.Provided{OutputTypeNames: []string{"string"}, Err: boom}, &fxevent.Provided{OutputTypeNames: []string{"string"}}},
		{"Replaced", &fxevent.Replaced{OutputTypeNames: []string{"string"}, Err: boom}, &fxevent.Replaced{OutputTypeNames: []string{"string"}}},
		{"Decorated", &fxevent.Decorated{OutputTypeNames: []string{"string"}, Err: boom}, &fxevent.Decorated{OutputTypeNames: []string{"string"}}},
		{"Run", &fxevent.Run{Name: "run", Kind: "provide", Err: boom}, &fxevent.Run{Name: "run", Kind: "provide"}},
		{"Invoked", &fxevent.Invoked{FunctionName: "invoke", Err: boom}, &fxevent.Invoked{FunctionName: "invoke"}},
		{"Started", &fxevent.Started{Err: boom}, &fxevent.Started{}},
		{"Stopped", &fxevent.Stopped{Err: boom}, &fxevent.Stopped{}},
		{"RollingBack", &fxevent.RollingBack{StartErr: boom}, nil},
		{"RolledBack", &fxevent.RolledBack{Err: boom}, &fxevent.RolledBack{}},
		{"LoggerInitialized", &fxevent.LoggerInitialized{ConstructorName: "logger", Err: boom}, &fxevent.LoggerInitialized{ConstructorName: "logger"}},
	} {
		t.Run(c.name, func(t *testing.T) {
			rec := record(t)
			NewEventLogger().LogEvent(c.failed)
			require.NotEmpty(t, rec.lines())
			for _, l := range rec.lines() {
				assert.Equal(t, "error", l.level, "%q", l.text)
			}
			if c.ok == nil {
				return
			}
			rec = record(t)
			NewEventLogger().LogEvent(c.ok)
			for _, l := range rec.lines() {
				assert.Equal(t, strings.ToLower(LogLevel()), l.level, "%q", l.text)
			}
		})
	}
}

// TestEventLogger_AnEventWithNothingToSayWritesNoLine: fx renders no line for
// a provide without outputs, and the logger writes none.
func TestEventLogger_AnEventWithNothingToSayWritesNoLine(t *testing.T) {
	boost.Start()
	rec := record(t)
	NewEventLogger().LogEvent(&fxevent.Provided{})
	NewEventLogger().LogEvent(&fxevent.Invoking{FunctionName: "invoke"})
	assert.Len(t, rec.lines(), 1, "only the invoke has a line: %v", rec.lines())
}

// TestLogger_PrintsAtItsLevel: each fx log level writes at its boost level,
// DEBUG for any other value.
func TestLogger_PrintsAtItsLevel(t *testing.T) {
	for level, want := range map[string]string{"INFO": "info", "TRACE": "trace", "DEBUG": "debug", "other": "debug"} {
		rec := record(t)
		(&Logger{level: level}).Printf("line %d", 1)
		assert.Equal(t, []entry{{level: want, text: "line 1"}}, rec.lines(), level)
	}
}
