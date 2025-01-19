package log

import (
	"context"
	"io"
	"log/slog"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/indrasaputra/arjuna/pkg/sdk/trace"
)

const (
	// EnvDevelopment is set to development.
	EnvDevelopment = "development"
	// EnvProduction is set to production.
	EnvProduction = "production"

	traceIDKey = "trace_id"
)

// Logger provides logging functionality.
type Logger struct {
	*zap.SugaredLogger
}

// Debugf adds `trace_id` from context to the log.
func (l *Logger) Debugf(ctx context.Context, template string, args ...interface{}) {
	id := trace.GetTraceIDFromContext(ctx)
	l.With(traceIDKey, id).Debugf(template, args...)
}

// Errorf adds `trace_id` from context to the log.
func (l *Logger) Errorf(ctx context.Context, template string, args ...interface{}) {
	id := trace.GetTraceIDFromContext(ctx)
	l.With(traceIDKey, id).Errorf(template, args...)
}

// Infof adds `trace_id` from context to the log.
func (l *Logger) Infof(ctx context.Context, template string, args ...interface{}) {
	id := trace.GetTraceIDFromContext(ctx)
	l.With(traceIDKey, id).Infof(template, args...)
}

// Warnf adds `trace_id` from context to the log.
func (l *Logger) Warnf(ctx context.Context, template string, args ...interface{}) {
	id := trace.GetTraceIDFromContext(ctx)
	l.With(traceIDKey, id).Warnf(template, args...)
}

// NewLogger creates an instance of Logger.
func NewLogger(env string) *Logger {
	c := newLoggerConfig(env)
	l, _ := c.Build(zap.AddCaller(), zap.AddCallerSkip(1))
	return &Logger{l.Sugar()}
}

func newLoggerConfig(env string) zap.Config {
	if env == EnvProduction {
		return zap.NewProductionConfig()
	}

	c := zap.NewDevelopmentConfig()
	c.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	c.EncoderConfig.TimeKey = ""
	return c
}

type SlogJSONHandler struct {
	*slog.JSONHandler
}

func NewSlogJSONHandler(w io.Writer, o *slog.HandlerOptions) *SlogJSONHandler {
	return &SlogJSONHandler{slog.NewJSONHandler(w, o)}
}

func (s *SlogJSONHandler) Handle(ctx context.Context, r slog.Record) error {
	s.addTraceId(ctx, &r)
	return s.JSONHandler.Handle(ctx, r)
}

func (s *SlogJSONHandler) addTraceId(ctx context.Context, r *slog.Record) {
	traceID := trace.GetTraceIDFromContext(ctx)
	r.AddAttrs(
		slog.Attr{Key: traceIDKey, Value: slog.StringValue(traceID)},
	)
}

func (s *SlogJSONHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &SlogJSONHandler{
		JSONHandler: s.JSONHandler.WithAttrs(attrs).(*slog.JSONHandler),
	}
}

func (s *SlogJSONHandler) WithGroup(name string) slog.Handler {
	return &SlogJSONHandler{
		JSONHandler: s.JSONHandler.WithGroup(name).(*slog.JSONHandler),
	}
}

func NewSlogLogger(svc string) *slog.Logger {
	h := NewSlogJSONHandler(os.Stdout, nil)
	l := slog.New(h)
	l = l.With(
		slog.Attr{Key: "service", Value: slog.StringValue(svc)},
	)
	return l
}
