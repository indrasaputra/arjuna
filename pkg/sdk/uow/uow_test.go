package uow_test

import (
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"

	"github.com/indrasaputra/arjuna/pkg/sdk/uow"
)

func TestNewTxManager(t *testing.T) {
	t.Run("success create tx manager", func(t *testing.T) {
		tx, err := uow.NewTxManager(&pgxpool.Pool{})

		assert.NoError(t, err)
		assert.NotNil(t, tx)
	})
}

func TestNewTxGetter(t *testing.T) {
	t.Run("success create tx getter", func(t *testing.T) {
		g := uow.NewTxGetter()

		assert.NotNil(t, g)
	})
}
