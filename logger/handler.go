package logger

import (
	"context"
	"fmt"
	"log/slog"
)

type SimpleHandler struct {
	minLevel slog.Level
}

func NewSimpleHanlder(minLevel slog.Level) *SimpleHandler {
	return &SimpleHandler{
		minLevel: minLevel,
	}
}

func (h *SimpleHandler) Handle(_ context.Context, r slog.Record) error {
	if r.Level < h.minLevel {
		return nil
	}

	msg := ""
	if r.Level == slog.LevelDebug {
		msg += "[Debug] "
	}

	msg += r.Message

	attrs := ""
	r.Attrs(func(a slog.Attr) bool {
		attrs += " " + a.Key + "=" + a.Value.String()
		return true
	})

	_, err := fmt.Printf("%s %s\n", msg, attrs)
	return err
}

// Enabled determines if a log level is enabled.
func (h *SimpleHandler) Enabled(_ context.Context, l slog.Level) bool {
	return l >= h.minLevel
}

// WithAttrs adds attributes to the handler.
func (h *SimpleHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return h
}

// WithGroup sets the group for the handler.
func (h *SimpleHandler) WithGroup(name string) slog.Handler {
	return h
}
