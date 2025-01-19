package log

import (
	"context"
	"io"
	"log/slog"
	"os"
	"runtime/debug"

	"github.com/indrasaputra/arjuna/pkg/sdk/trace"
)

const (
	traceIDKey = "trace_id"
)

// SlogJSONHandler is a wrapper for slog.JSONHandler.
type SlogJSONHandler struct {
	*slog.JSONHandler
}

// NewSlogJSONHandler creates an instance of SlogJSONHandler.
func NewSlogJSONHandler(w io.Writer, o *slog.HandlerOptions) *SlogJSONHandler {
	return &SlogJSONHandler{slog.NewJSONHandler(w, o)}
}

// Handle overrides the Handle method from slog.JSONHandler
// with the imbued trace_id and stacktrace.
func (s *SlogJSONHandler) Handle(ctx context.Context, r slog.Record) error {
	s.addTraceID(ctx, &r)
	s.printStackTrace(&r)
	return s.JSONHandler.Handle(ctx, r)
}

func (s *SlogJSONHandler) addTraceID(ctx context.Context, r *slog.Record) {
	traceID := trace.GetTraceIDFromContext(ctx)
	r.AddAttrs(
		slog.Attr{Key: traceIDKey, Value: slog.StringValue(traceID)},
	)
}

func (s *SlogJSONHandler) printStackTrace(r *slog.Record) {
	if r.Level >= slog.LevelError {
		// Add a stack trace as an attribute
		r.AddAttrs(slog.Attr{
			Key:   "stacktrace",
			Value: slog.StringValue(string(debug.Stack())),
		})
	}
}

// WithAttrs overrides the WithAttrs method from slog.JSONHandler.
func (s *SlogJSONHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &SlogJSONHandler{
		JSONHandler: s.JSONHandler.WithAttrs(attrs).(*slog.JSONHandler),
	}
}

// WithGroup overrides the WithGroup method from slog.JSONHandler.
func (s *SlogJSONHandler) WithGroup(name string) slog.Handler {
	return &SlogJSONHandler{
		JSONHandler: s.JSONHandler.WithGroup(name).(*slog.JSONHandler),
	}
}

// NewSlogLogger creates an instance of slog.Logger
// using SlogJSONHandler as the Logger. It will write log to stdout.
func NewSlogLogger(svc string) *slog.Logger {
	h := NewSlogJSONHandler(os.Stdout, nil)
	l := slog.New(h)
	l = l.With(
		slog.Attr{Key: "service", Value: slog.StringValue(svc)},
	)
	return l
}
