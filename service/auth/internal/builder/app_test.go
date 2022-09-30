package builder_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	mock_keycloak "github.com/indrasaputra/arjuna/pkg/sdk/test/mock/keycloak"
	"github.com/indrasaputra/arjuna/service/auth/internal/builder"
	"github.com/indrasaputra/arjuna/service/auth/internal/config"
)

func TestBuildAuthHandler(t *testing.T) {
	t.Run("fail create auth handler", func(t *testing.T) {
		dep := &builder.Dependency{
			KeycloakClient: nil,
			Config: &config.Config{
				Keycloak: config.Keycloak{},
			},
		}

		handler, err := builder.BuildAuthHandler(dep)

		assert.Error(t, err)
		assert.Nil(t, handler)
	})

	t.Run("success create auth handler", func(t *testing.T) {
		dep := &builder.Dependency{
			KeycloakClient: &mock_keycloak.MockKeycloak{},
			Config: &config.Config{
				Keycloak: config.Keycloak{
					Realm: "realm",
				},
			},
		}

		handler, err := builder.BuildAuthHandler(dep)

		assert.NoError(t, err)
		assert.NotNil(t, handler)
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
