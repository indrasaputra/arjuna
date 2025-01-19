package builder_test

import (
	"testing"

	"github.com/pashagolub/pgxmock/v4"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	mock_uow "github.com/indrasaputra/arjuna/pkg/sdk/test/mock/uow"
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
