package fx

import (
	"bytes"
	"strings"

	"github.com/xgodev/boost/wrapper/log"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
)

// EventLogger is the fx event logger NewApp installs. It writes fx's lines
// through the boost logger at the fx log level, and every event that carries
// an error (a constructor, invoke or hook that failed, an app that did not
// start or stop) at the error level whatever the fx log level, so the reason
// an app exits always reaches the output.
type EventLogger struct {
	printer fx.Printer
}

// NewEventLogger returns the fx event logger, at the fx log level.
func NewEventLogger() fxevent.Logger {
	return &EventLogger{printer: NewLogger()}
}

// LogEvent writes the line fx's console logger renders for event.
func (l *EventLogger) LogEvent(event fxevent.Event) {
	var line bytes.Buffer
	(&fxevent.ConsoleLogger{W: &line}).LogEvent(event)
	text := strings.TrimSuffix(line.String(), "\n")
	if text == "" {
		return
	}
	if eventErr(event) != nil {
		log.Error(text)
		return
	}
	l.printer.Printf("%s", text)
}

// eventErr is the error event carries, nil when it carries none.
func eventErr(event fxevent.Event) error {
	switch e := event.(type) {
	case *fxevent.OnStartExecuted:
		return e.Err
	case *fxevent.OnStopExecuted:
		return e.Err
	case *fxevent.Supplied:
		return e.Err
	case *fxevent.Provided:
		return e.Err
	case *fxevent.Replaced:
		return e.Err
	case *fxevent.Decorated:
		return e.Err
	case *fxevent.Run:
		return e.Err
	case *fxevent.Invoked:
		return e.Err
	case *fxevent.Started:
		return e.Err
	case *fxevent.Stopped:
		return e.Err
	case *fxevent.RollingBack:
		return e.StartErr
	case *fxevent.RolledBack:
		return e.Err
	case *fxevent.LoggerInitialized:
		return e.Err
	default:
		return nil
	}
}
