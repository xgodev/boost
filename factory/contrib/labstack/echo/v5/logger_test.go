package echo

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/xgodev/boost/wrapper/log"
)

// recorder is a boost logger that keeps what the slog bridge hands it. The
// embedded nil Logger panics on any method the bridge is not meant to call.
type recorder struct {
	log.Logger
	fields  map[string]interface{}
	entries *[]string
}

func newRecorder() recorder { return recorder{entries: &[]string{}} }

func (r recorder) add(level, format string, args ...interface{}) {
	*r.entries = append(*r.entries, fmt.Sprintf("%s %s %v", level, fmt.Sprintf(format, args...), r.fields))
}

func (r recorder) Debugf(format string, args ...interface{}) { r.add("DEBUG", format, args...) }
func (r recorder) Infof(format string, args ...interface{})  { r.add("INFO", format, args...) }
func (r recorder) Warnf(format string, args ...interface{})  { r.add("WARN", format, args...) }
func (r recorder) Errorf(format string, args ...interface{}) { r.add("ERROR", format, args...) }

func (r recorder) WithFields(kv map[string]interface{}) log.Logger {
	fields := map[string]interface{}{}
	for k, v := range r.fields {
		fields[k] = v
	}
	for k, v := range kv {
		fields[k] = v
	}
	return recorder{fields: fields, entries: r.entries}
}

func TestWrapLoggerMapsLevelsToTheBoostLogger(t *testing.T) {
	rec := newRecorder()
	logger := WrapLogger(rec)

	logger.Debug("d")
	logger.Info("i")
	logger.Warn("w")
	logger.Error("e")

	assert.Equal(t, []string{"DEBUG d map[]", "INFO i map[]", "WARN w map[]", "ERROR e map[]"}, *rec.entries)
}

func TestWrapLoggerCarriesAttributesAndGroupsAsFields(t *testing.T) {
	rec := newRecorder()
	logger := WrapLogger(rec).With("service", "orders").WithGroup("http")

	logger.Info("started", "address", ":8080", slog.Group("tls", "enabled", true))

	assert.Equal(t,
		[]string{"INFO started map[http.address::8080 http.tls.enabled:true service:orders]"},
		*rec.entries)
}

func TestWrapLoggerLogsErrorsAsStrings(t *testing.T) {
	rec := newRecorder()

	WrapLogger(rec).Error("failed", "error", errors.New("boom"))

	assert.Equal(t, []string{"ERROR failed map[error:boom]"}, *rec.entries)
}

func TestWrapLoggerIsEnabledForEveryLevel(t *testing.T) {
	logger := WrapLogger(newRecorder())

	assert.True(t, logger.Enabled(context.Background(), slog.LevelDebug),
		"the boost logger filters levels itself")
}
