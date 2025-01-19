package postgres_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/indrasaputra/arjuna/pkg/sdk/database/postgres"
)

func TestNewPgxPool(t *testing.T) {
	t.Run("success create pgx pool", func(t *testing.T) {
		pool, err := postgres.NewPgxPool(postgres.Config{})

		assert.NoError(t, err)
		assert.NotNil(t, pool)
	})
}
