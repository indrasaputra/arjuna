package builder_test

import (
	"testing"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/stretchr/testify/assert"

	"github.com/indrasaputra/arjuna/service/user/internal/builder"
	"github.com/indrasaputra/arjuna/service/user/internal/config"
)

func TestBuildToggleCommandHandler(t *testing.T) {
	t.Run("success create toggle command handler", func(t *testing.T) {
		dep := &builder.Dependency{
			PgxPool: &pgxpool.Pool{},
		}

		handler := builder.BuildUserCommandHandler(dep)

		assert.NotNil(t, handler)
	})
}

func TestBuildPostgrePgxPool(t *testing.T) {
	cfg := &config.Postgres{
		Host:            "localhost",
		Port:            "5432",
		Name:            "users",
		User:            "user",
		Password:        "password",
		MaxOpenConns:    "10",
		MaxConnLifetime: "10m",
		MaxIdleLifetime: "5m",
		SSLMode:         "disable",
	}

	t.Run("fail build postgres pgxpool client", func(t *testing.T) {
		client, err := builder.BuildPostgrePgxPool(cfg)

		assert.Error(t, err)
		assert.Nil(t, client)
	})
}
