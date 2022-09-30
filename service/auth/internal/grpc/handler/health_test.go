package handler_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/health/grpc_health_v1"

	"github.com/indrasaputra/arjuna/service/auth/internal/grpc/handler"
)

var (
	testCtx                = context.Background()
	testHealthCheckRequest = &grpc_health_v1.HealthCheckRequest{Service: "auth"}
)

func TestNewHealth(t *testing.T) {
	t.Run("successful create an instance of Health", func(t *testing.T) {
		health := handler.NewHealth()
		assert.NotNil(t, health)
	})
}

func TestHealth_Check(t *testing.T) {
	t.Run("nil request is prohibited", func(t *testing.T) {
		health := handler.NewHealth()

		resp, err := health.Check(testCtx, nil)

		assert.NotNil(t, err)
		assert.Equal(t, grpc_health_v1.HealthCheckResponse_UNKNOWN, resp.GetStatus())
	})

	t.Run("system is healthy", func(t *testing.T) {
		health := handler.NewHealth()

		resp, err := health.Check(testCtx, testHealthCheckRequest)

		assert.Nil(t, err)
		assert.Equal(t, grpc_health_v1.HealthCheckResponse_SERVING, resp.GetStatus())
	})
}
