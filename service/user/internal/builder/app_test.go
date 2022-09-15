package builder_test

import (
	"testing"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/stretchr/testify/assert"

	mock_keycloak "github.com/indrasaputra/arjuna/pkg/sdk/test/mock/keycloak"
	"github.com/indrasaputra/arjuna/service/user/internal/builder"
	"github.com/indrasaputra/arjuna/service/user/internal/config"
)

func TestBuildUserCommandHandler(t *testing.T) {
	t.Run("fail create user command handler", func(t *testing.T) {
		dep := &builder.Dependency{
			PgxPool:        &pgxpool.Pool{},
			KeycloakClient: nil,
			Config: &config.Config{
				Keycloak: config.Keycloak{},
			},
		}

		handler, err := builder.BuildUserCommandHandler(dep)

		assert.Error(t, err)
		assert.Nil(t, handler)
	})

	t.Run("success create user command handler", func(t *testing.T) {
		dep := &builder.Dependency{
			PgxPool:        &pgxpool.Pool{},
			KeycloakClient: &mock_keycloak.MockKeycloak{},
			Config: &config.Config{
				Keycloak: config.Keycloak{
					Realm:         "realm",
					AdminUser:     "admin",
					AdminPassword: "admin",
				},
			},
		}

		handler, err := builder.BuildUserCommandHandler(dep)

		assert.NoError(t, err)
		assert.NotNil(t, handler)
	})
}

func TestBuildUserCommandInternalHandler(t *testing.T) {
	t.Run("fail create user command internal handler", func(t *testing.T) {
		dep := &builder.Dependency{
			PgxPool:        &pgxpool.Pool{},
			KeycloakClient: nil,
			Config: &config.Config{
				Keycloak: config.Keycloak{},
			},
		}

		handler, err := builder.BuildUserCommandInternalHandler(dep)

		assert.Error(t, err)
		assert.Nil(t, handler)
	})

	t.Run("success create user command internal handler", func(t *testing.T) {
		dep := &builder.Dependency{
			PgxPool:        &pgxpool.Pool{},
			KeycloakClient: &mock_keycloak.MockKeycloak{},
			Config: &config.Config{
				Keycloak: config.Keycloak{
					Realm:         "realm",
					AdminUser:     "admin",
					AdminPassword: "admin",
				},
			},
		}

		handler, err := builder.BuildUserCommandInternalHandler(dep)

		assert.NoError(t, err)
		assert.NotNil(t, handler)
	})
}

func TestBuildUserQueryHandler(t *testing.T) {
	t.Run("success create user query handler", func(t *testing.T) {
		dep := &builder.Dependency{
			PgxPool: &pgxpool.Pool{},
		}

		handler := builder.BuildUserQueryHandler(dep)

		assert.NotNil(t, handler)
	})
}

func TestBuildPostgrePgxPool(t *testing.T) {
	cfg := config.Postgres{
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

func TestBuildKeycloakClient(t *testing.T) {
	cfg := config.Keycloak{
		Timeout: 5,
	}

	t.Run("success build a keycloak client", func(t *testing.T) {
		client := builder.BuildKeycloakClient(cfg)
		assert.NotNil(t, client)
	})
}
