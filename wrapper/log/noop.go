package log

import (
	"context"
	"io"
)

// Noop is a dummy implementation of Logger
type Noop struct{}

// NewNoop create a Noop Logger
func NewNoop() Noop {
	i := Noop{}
	SetGlobalLogger(i)
	return i
}

func (n Noop) Printf(format string, args ...interface{}) {}

func (n Noop) Tracef(format string, args ...interface{}) {}

func (n Noop) Trace(args ...interface{}) {}

func (n Noop) Debugf(format string, args ...interface{}) {}

func (n Noop) Debug(args ...interface{}) {}

func (n Noop) Infof(format string, args ...interface{}) {}

func (n Noop) Info(args ...interface{}) {}

func (n Noop) Warnf(format string, args ...interface{}) {}

func (n Noop) Warn(args ...interface{}) {}

func (n Noop) Errorf(format string, args ...interface{}) {}

func (n Noop) Error(args ...interface{}) {}

func (n Noop) Fatalf(format string, args ...interface{}) {}

func (n Noop) Fatal(args ...interface{}) {}

func (n Noop) Panicf(format string, args ...interface{}) {}

func (n Noop) Panic(args ...interface{}) {}

func (n Noop) WithFields(keyValues map[string]interface{}) Logger { return n }

func (n Noop) WithField(key string, value interface{}) Logger { return n }

func (n Noop) WithError(err error) Logger { return n }

func (n Noop) WithTypeOf(obj interface{}) Logger { return n }

func (n Noop) ToContext(ctx context.Context) context.Context { return ctx }

func (n Noop) FromContext(ctx context.Context) Logger { return n }

func (n Noop) Output() io.Writer { return io.Discard }

func (n Noop) Fields() Fields { return Fields{} }
