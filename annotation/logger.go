package annotation

import (
	"fmt"
)

var log = NewLogger()

func WithLogger(logger Logger) {
	log = logger
}

// Logger is our contract for the logger.
type Logger interface {
	Tracef(format string, args ...interface{})

	Trace(args ...interface{})

	Debugf(format string, args ...interface{})

	Debug(args ...interface{})

	Infof(format string, args ...interface{})

	Info(args ...interface{})

	Warnf(format string, args ...interface{})

	Warn(args ...interface{})

	Errorf(format string, args ...interface{})

	Error(args ...interface{})

	Fatalf(format string, args ...interface{})

	Fatal(args ...interface{})

	Panicf(format string, args ...interface{})

	Panic(args ...interface{})
}

type DefaultLogger struct{}

func (n DefaultLogger) Tracef(format string, args ...interface{}) {
	fmt.Println(fmt.Sprintf(format, args...))
}

func (n DefaultLogger) Trace(args ...interface{}) {
	fmt.Println(fmt.Sprintln(args...))
}

func (n DefaultLogger) Debugf(format string, args ...interface{}) {
	fmt.Println(fmt.Sprintf(format, args...))
}

func (n DefaultLogger) Debug(args ...interface{}) {
	fmt.Println(fmt.Sprintln(args...))
}

func (n DefaultLogger) Infof(format string, args ...interface{}) {
	fmt.Println(fmt.Sprintf(format, args...))
}

func (n DefaultLogger) Info(args ...interface{}) {
	fmt.Println(fmt.Sprintln(args...))
}

func (n DefaultLogger) Warnf(format string, args ...interface{}) {
	fmt.Println(fmt.Sprintf(format, args...))
}

func (n DefaultLogger) Warn(args ...interface{}) {
	fmt.Println(fmt.Sprintln(args...))
}

func (n DefaultLogger) Errorf(format string, args ...interface{}) {
	fmt.Println(fmt.Sprintf(format, args...))
}

func (n DefaultLogger) Error(args ...interface{}) {
	fmt.Println(fmt.Sprintln(args...))
}

func (n DefaultLogger) Fatalf(format string, args ...interface{}) {
	fmt.Println(fmt.Sprintf(format, args...))
}

func (n DefaultLogger) Fatal(args ...interface{}) {
	fmt.Println(fmt.Sprintln(args...))
}

func (n DefaultLogger) Panicf(format string, args ...interface{}) {
	fmt.Println(fmt.Sprintf(format, args...))
}

func (n DefaultLogger) Panic(args ...interface{}) {
	fmt.Println(fmt.Sprintln(args...))
}

func NewLogger() Logger {
	return &DefaultLogger{}
}
