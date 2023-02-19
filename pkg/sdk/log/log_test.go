package log_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/indrasaputra/arjuna/pkg/sdk/log"
)

var (
	testCtx = context.Background()
)

func TestNewLogger(t *testing.T) {
	t.Run("success create development logger", func(t *testing.T) {
		l := log.NewLogger(log.EnvDevelopment)
		assert.NotNil(t, l)
	})

	t.Run("success create production logger", func(t *testing.T) {
		l := log.NewLogger(log.EnvProduction)
		assert.NotNil(t, l)
	})
}

func TestLogger_Debugf(t *testing.T) {
	t.Run("success emit log for Debugf", func(t *testing.T) {
		l := log.NewLogger(log.EnvDevelopment)
		assert.NotPanics(t, func() { l.Debugf(testCtx, "Debugf") })
	})
}

func TestLogger_Errorf(t *testing.T) {
	t.Run("success emit log for Errorf", func(t *testing.T) {
		l := log.NewLogger(log.EnvDevelopment)
		assert.NotPanics(t, func() { l.Errorf(testCtx, "Errorf") })
	})
}

func TestLogger_Infof(t *testing.T) {
	t.Run("success emit log for Infof", func(t *testing.T) {
		l := log.NewLogger(log.EnvDevelopment)
		assert.NotPanics(t, func() { l.Infof(testCtx, "Infof") })
	})
}

func TestLogger_Warnf(t *testing.T) {
	t.Run("success emit log for Warnf", func(t *testing.T) {
		l := log.NewLogger(log.EnvDevelopment)
		assert.NotPanics(t, func() { l.Warnf(testCtx, "debugf") })
	})
}
