package echo

import (
	"context"
	"log/slog"

	"github.com/xgodev/boost/wrapper/log"
)

// WrapLogger returns a slog logger, the logger echo v5 uses, that writes
// through the boost logger. Attributes become boost fields; group names
// prefix the keys of the attributes inside them ("http.address").
func WrapLogger(logger log.Logger) *slog.Logger {
	return slog.New(&handler{logger: logger, fields: map[string]interface{}{}})
}

// handler is the slog.Handler behind WrapLogger.
type handler struct {
	logger log.Logger
	fields map[string]interface{}
	prefix string
}

// Enabled reports true for every level: the boost logger filters levels itself.
func (h *handler) Enabled(context.Context, slog.Level) bool {
	return true
}

// Handle writes the record through the boost logger at the matching level.
func (h *handler) Handle(_ context.Context, r slog.Record) error {
	fields := h.copyFields()
	r.Attrs(func(a slog.Attr) bool {
		addAttr(fields, h.prefix, a)
		return true
	})

	logger := h.logger.WithFields(fields)
	switch {
	case r.Level >= slog.LevelError:
		logger.Errorf("%s", r.Message)
	case r.Level >= slog.LevelWarn:
		logger.Warnf("%s", r.Message)
	case r.Level >= slog.LevelInfo:
		logger.Infof("%s", r.Message)
	default:
		logger.Debugf("%s", r.Message)
	}
	return nil
}

// WithAttrs returns a handler that adds attrs to every record.
func (h *handler) WithAttrs(attrs []slog.Attr) slog.Handler {
	fields := h.copyFields()
	for _, a := range attrs {
		addAttr(fields, h.prefix, a)
	}
	return &handler{logger: h.logger, fields: fields, prefix: h.prefix}
}

// WithGroup returns a handler that prefixes later attribute keys with name.
func (h *handler) WithGroup(name string) slog.Handler {
	if name == "" {
		return h
	}
	return &handler{logger: h.logger, fields: h.fields, prefix: h.prefix + name + "."}
}

func (h *handler) copyFields() map[string]interface{} {
	fields := make(map[string]interface{}, len(h.fields))
	for k, v := range h.fields {
		fields[k] = v
	}
	return fields
}

// addAttr flattens a into fields under prefix: groups add their key to the
// prefix and errors are kept as their message.
func addAttr(fields map[string]interface{}, prefix string, a slog.Attr) {
	a.Value = a.Value.Resolve()
	if a.Equal(slog.Attr{}) {
		return
	}
	if a.Value.Kind() == slog.KindGroup {
		if a.Key != "" {
			prefix += a.Key + "."
		}
		for _, ga := range a.Value.Group() {
			addAttr(fields, prefix, ga)
		}
		return
	}
	v := a.Value.Any()
	if err, ok := v.(error); ok {
		v = err.Error()
	}
	fields[prefix+a.Key] = v
}
