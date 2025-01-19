package log_test

import (
	"context"
	"log/slog"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/indrasaputra/arjuna/pkg/sdk/log"
)

var (
	testCtx = context.Background()
)

func TestNewSlogJSONHandler(t *testing.T) {
	t.Run("successfully create an instance of SlogJSONHandler", func(t *testing.T) {
		handler := log.NewSlogJSONHandler(os.Stdout, nil)

		assert.NotNil(t, handler)
	})
}

func TestSlogJSONHandler_Handle(t *testing.T) {
	t.Run("successfully handle a log", func(t *testing.T) {
		handler := log.NewSlogJSONHandler(os.Stdout, nil)
		err := handler.Handle(testCtx, slog.Record{Level: slog.LevelError, Message: "error"})

		assert.NoError(t, err)
	})
}

func TestSlogJSONHandler_WithAttrs(t *testing.T) {
	t.Run("successfully add attributes to the log", func(t *testing.T) {
		handler := log.NewSlogJSONHandler(os.Stdout, nil)
		newHandler := handler.WithAttrs([]slog.Attr{{Key: "key", Value: slog.StringValue("value")}})

		assert.NotNil(t, newHandler)
	})
}

func TestSlogJSONHandler_WithGroup(t *testing.T) {
	t.Run("successfully add group to the log", func(t *testing.T) {
		handler := log.NewSlogJSONHandler(os.Stdout, nil)
		newHandler := handler.WithGroup("group")

		assert.NotNil(t, newHandler)
	})
}

func TestNewSlogLogger(t *testing.T) {
	t.Run("successfully create an instance of SlogLogger", func(t *testing.T) {
		logger := log.NewSlogLogger("my-service")

		assert.NotNil(t, logger)
	})
}
