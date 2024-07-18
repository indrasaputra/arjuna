package builder_test

import (
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/stretchr/testify/assert"

	sdkpg "github.com/indrasaputra/arjuna/pkg/sdk/database/postgres"
	"github.com/indrasaputra/arjuna/service/wallet/internal/builder"
	"github.com/indrasaputra/arjuna/service/wallet/internal/config"
)

func TestBuildWalletCommandHandler(t *testing.T) {
	t.Run("success create wallet command handler", func(t *testing.T) {
		dep := &builder.Dependency{
			Config: &config.Config{},
		}

		handler := builder.BuildWalletCommandHandler(dep)

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
