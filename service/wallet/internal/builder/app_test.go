package builder_test

import (
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/pashagolub/pgxmock/v4"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	mock_uow "github.com/indrasaputra/arjuna/pkg/sdk/test/mock/uow"
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

func TestBuildQueries(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("success create queries", func(t *testing.T) {
		pool, err := pgxmock.NewPool()
		if err != nil {
			t.Fatalf("error opening a stub database connection: %v\n", err)
		}
		g := mock_uow.NewMockTxGetter(ctrl)

		queries := builder.BuildQueries(pool, g)

		assert.NotNil(t, queries)
	})
}
