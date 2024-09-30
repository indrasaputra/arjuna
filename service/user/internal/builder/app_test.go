package builder_test

import (
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/stretchr/testify/assert"

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

func TestBuildTemporalClient(t *testing.T) {
	t.Run("fail build a temporal client", func(t *testing.T) {
		client, err := builder.BuildTemporalClient("localhost:7233")

		assert.Error(t, err)
		assert.Nil(t, client)
	})
}

func TestBuildAuthClient(t *testing.T) {
	t.Run("success build an auth client", func(t *testing.T) {
		client, err := builder.BuildAuthClient("localhost:8002", "user", "pass")

		assert.NoError(t, err)
		assert.NotNil(t, client)
	})
}

func TestBuildWalletClient(t *testing.T) {
	t.Run("success build a wallet client", func(t *testing.T) {
		client, err := builder.BuildWalletClient("localhost:8004", "user", "pass")

		assert.NoError(t, err)
		assert.NotNil(t, client)
	})
}

func TestBuildRedisClient(t *testing.T) {
	t.Run("fail create redis client", func(t *testing.T) {
		server, _ := miniredis.Run()

		cfg := &config.Redis{
			Address: server.Addr(),
		}

		server.Close()
		client, err := builder.BuildRedisClient(cfg)

		assert.NotNil(t, err)
		assert.Nil(t, client)
	})

	t.Run("success create redis client", func(t *testing.T) {
		server, _ := miniredis.Run()
		defer server.Close()

		cfg := &config.Redis{
			Address: server.Addr(),
		}

		client, err := builder.BuildRedisClient(cfg)

		assert.Nil(t, err)
		assert.NotNil(t, client)
	})
}
