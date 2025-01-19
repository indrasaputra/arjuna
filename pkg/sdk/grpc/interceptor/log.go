package interceptor

import (
	"context"
	"log/slog"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
)

// SlogLogger intercepts logs with log/slog.
// Copied from https://github.com/grpc-ecosystem/go-grpc-middleware/blob/main/interceptors/logging/examples/slog/example_test.go
func SlogLogger(l *slog.Logger) logging.Logger {
	return logging.LoggerFunc(func(ctx context.Context, lvl logging.Level, msg string, fields ...any) {
		l.Log(ctx, slog.Level(lvl), msg, fields...)
	})
}
