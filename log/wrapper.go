package log

// A global variable so that l functions can be directly accessed.
var (
	l Logger = Noop{}
)

// NewLogger returns an instance of logger.
// Deprecated: prefer SetGlobalLogger
func NewLogger(logger Logger) {
	l = logger
}

func SetGlobalLogger(logger Logger) {
	l = logger
}

// Printf logs a templated message.
//
// For templating details see implementation doc.
func Printf(format string, args ...interface{}) {
	l.Printf(format, args...)
}

// Tracef logs a templated message at trace level.
//
// For templating details see implementation doc.
func Tracef(format string, args ...interface{}) {
	l.Tracef(format, args...)
}

// Trace logs a message at trace level.
func Trace(args ...interface{}) {
	l.Trace(args...)
}

// Debugf logs a templated message at debug level.
//
// For templating details see implementation doc.
func Debugf(format string, args ...interface{}) {
	l.Debugf(format, args...)
}

// Debug logs a message at debug level.
func Debug(args ...interface{}) {
	l.Debug(args...)
}

// Infof logs a templated message at info level.
//
// For templating details see implementation doc.
func Infof(format string, args ...interface{}) {
	l.Infof(format, args...)
}

// Info logs a message at info level.
func Info(args ...interface{}) {
	l.Info(args...)
}

// Warnf logs a templated message at warn level.
//
// For templating details see implementation doc.
func Warnf(format string, args ...interface{}) {
	l.Warnf(format, args...)
}

// Warn logs a message at warn level.
func Warn(args ...interface{}) {
	l.Warn(args...)
}

// Errorf logs a templated message at error level.
//
// For templating details see implementation doc.
func Errorf(format string, args ...interface{}) {
	l.Errorf(format, args...)
}

// Error logs a message at error level.
func Error(args ...interface{}) {
	l.Error(args...)
}

// Panicf is equivalent to Printf() followed by a call to panic().
func Panicf(format string, args ...interface{}) {
	l.Panicf(format, args...)
}

// Panic is equivalent to Print() followed by a call to panic().
func Panic(args ...interface{}) {
	l.Panic(args...)
}

// Fatal is equivalent to Print() followed by a call to os.Exit(1).
func Fatal(args ...interface{}) {
	l.Fatal(args...)
}

// Fatalf is equivalent to Printf() followed by a call to os.Exit(1).
func Fatalf(format string, args ...interface{}) {
	l.Fatalf(format, args...)
}

// WithField adds a key and value to logger.
func WithField(key string, value interface{}) Logger {
	return l.WithField(key, value)
}

// WithError adds an error as a field to logger
func WithError(err error) Logger {
	return l.WithError(err)
}

// WithFields adds fields to logger.
func WithFields(keyValues map[string]interface{}) Logger {
	return l.WithFields(keyValues)
}

// WithTypeOf adds type information to logger.
func WithTypeOf(obj interface{}) Logger {
	return l.WithTypeOf(obj)
}

// GetLogger returns instance of Logger.
func GetLogger() Logger {
	return l
}
