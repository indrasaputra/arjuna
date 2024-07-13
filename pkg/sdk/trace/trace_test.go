package trace_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/indrasaputra/arjuna/pkg/sdk/trace"
)

func TestGetTraceIDFromContext(t *testing.T) {
	t.Run("get default trace id from context without trace", func(t *testing.T) {
		id := trace.GetTraceIDFromContext(context.Background())

		assert.NotEmpty(t, id)
		assert.Equal(t, "00000000000000000000000000000000", id)
	})
}

func TestNewProvider(t *testing.T) {
	ctx := context.Background()

	t.Run("success create a new provider", func(t *testing.T) {
		cfg := trace.Config{
			OtelCollectorAddress: "localhost:4317",
			AppEnv:               "production",
		}

		prov, err := trace.NewProvider(ctx, cfg)

		assert.NoError(t, err)
		assert.NotNil(t, prov)
	})
}
