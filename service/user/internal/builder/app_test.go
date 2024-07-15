package builder_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	pgsdk "github.com/indrasaputra/arjuna/pkg/sdk/database/postgres"
	mock_keycloak "github.com/indrasaputra/arjuna/pkg/sdk/test/mock/keycloak"
	"github.com/indrasaputra/arjuna/service/user/internal/builder"
	"github.com/indrasaputra/arjuna/service/user/internal/config"
)

func TestBuildUserCommandHandler(t *testing.T) {
	t.Run("success create user command handler", func(t *testing.T) {
		dep := &builder.Dependency{
			KeycloakClient: &mock_keycloak.MockKeycloak{},
			Config: &config.Config{
				Keycloak: config.Keycloak{
					Realm:         "realm",
					AdminUser:     "admin",
					AdminPassword: "admin",
				},
			},
		}

		handler := builder.BuildUserCommandHandler(dep)

		assert.NotNil(t, handler)
	})
}

func TestBuildUserCommandInternalHandler(t *testing.T) {
	t.Run("success create user command internal handler", func(t *testing.T) {
		dep := &builder.Dependency{
			KeycloakClient: &mock_keycloak.MockKeycloak{},
			Config: &config.Config{
				Keycloak: config.Keycloak{
					Realm:         "realm",
					AdminUser:     "admin",
					AdminPassword: "admin",
				},
			},
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

func TestBuildBunDB(t *testing.T) {
	t.Run("success create bundb", func(t *testing.T) {
		cfg := pgsdk.Config{}

		db, err := builder.BuildBunDB(cfg)

		assert.NoError(t, err)
		assert.NotNil(t, db)
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

func TestBuildTemporalClient(t *testing.T) {
	t.Run("fail build a temporal client", func(t *testing.T) {
		client, err := builder.BuildTemporalClient("localhost:7233")

		assert.Error(t, err)
		assert.Nil(t, client)
	})
}
