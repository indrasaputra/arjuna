package builder_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	sdkpg "github.com/indrasaputra/arjuna/pkg/sdk/database/postgres"
	"github.com/indrasaputra/arjuna/service/auth/internal/builder"
	"github.com/indrasaputra/arjuna/service/auth/internal/config"
)

func TestBuildAuthHandler(t *testing.T) {
	t.Run("fail create auth handler", func(t *testing.T) {
		dep := &builder.Dependency{
			Config: &config.Config{},
		}

		handler, err := builder.BuildAuthHandler(dep)

		assert.NoError(t, err)
		assert.NotNil(t, handler)
	})
}

func TestBuildBunDB(t *testing.T) {
	t.Run("success create bundb", func(t *testing.T) {
		cfg := sdkpg.Config{}

		db, err := builder.BuildBunDB(cfg)

		assert.NoError(t, err)
		assert.NotNil(t, db)
	})
}
