package builder_test

import (
	"testing"

	"github.com/pashagolub/pgxmock/v4"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	mock_uow "github.com/indrasaputra/arjuna/pkg/sdk/test/mock/uow"
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
