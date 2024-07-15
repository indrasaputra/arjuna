package builder_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	sdkpg "github.com/indrasaputra/arjuna/pkg/sdk/database/postgres"
	"github.com/indrasaputra/arjuna/service/user/internal/builder"
	"github.com/indrasaputra/arjuna/service/user/internal/config"
)

func TestBuildUserCommandHandler(t *testing.T) {
	t.Run("success create user command handler", func(t *testing.T) {
		dep := &builder.Dependency{
			Config: &config.Config{},
		}

		handler := builder.BuildUserCommandHandler(dep)

		assert.NotNil(t, handler)
	})
}

func TestBuildUserCommandInternalHandler(t *testing.T) {
	t.Run("success create user command internal handler", func(t *testing.T) {
		dep := &builder.Dependency{
			Config: &config.Config{},
		}

		handler := builder.BuildUserCommandInternalHandler(dep)

		assert.NotNil(t, handler)
	})
}

func TestBuildUserQueryHandler(t *testing.T) {
	t.Run("success create user query handler", func(t *testing.T) {
		dep := &builder.Dependency{}

		handler := builder.BuildUserQueryHandler(dep)

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

func TestBuildTemporalClient(t *testing.T) {
	t.Run("fail build a temporal client", func(t *testing.T) {
		client, err := builder.BuildTemporalClient("localhost:7233")

		assert.Error(t, err)
		assert.Nil(t, client)
	})
}

func TestBuildAuthClient(t *testing.T) {
	t.Run("success build an auth client", func(t *testing.T) {
		client, err := builder.BuildAuthClient("localhost:8002")

		assert.NoError(t, err)
		assert.NotNil(t, client)
	})
}
