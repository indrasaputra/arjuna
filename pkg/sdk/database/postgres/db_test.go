package postgres_test

import (
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"

	"github.com/indrasaputra/arjuna/pkg/sdk/database/postgres"
)

func TestNewTxGetter(t *testing.T) {
	t.Run("success create tx getter", func(t *testing.T) {
		g := postgres.NewTxGetter()

		assert.NotNil(t, g)
	})
}

func TestNewTxManager(t *testing.T) {
	t.Run("success create tx manager", func(t *testing.T) {
		tx, err := postgres.NewTxManager(&pgxpool.Pool{})

		assert.NoError(t, err)
		assert.NotNil(t, tx)
	})
}

func TestNewPgxPool(t *testing.T) {
	t.Run("success create pgx pool", func(t *testing.T) {
		pool, err := postgres.NewPgxPool(postgres.Config{})

		assert.NoError(t, err)
		assert.NotNil(t, pool)
	})
}
