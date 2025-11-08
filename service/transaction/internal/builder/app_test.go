package builder_test

import (
	"testing"

	"github.com/pashagolub/pgxmock/v4"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	mock_uow "github.com/indrasaputra/arjuna/pkg/sdk/test/mock/uow"
	"github.com/indrasaputra/arjuna/service/transaction/internal/builder"
	"github.com/indrasaputra/arjuna/service/transaction/internal/config"
)

func TestBuildTransactionCommandHandler(t *testing.T) {
	t.Run("success create transaction command handler", func(t *testing.T) {
		dep := &builder.Dependency{
			Config: &config.Config{},
		}

		handler := builder.BuildTransactionCommandHandler(dep)

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
