package builder

import (
	"net/http"
	"time"

	kcsdk "github.com/indrasaputra/arjuna/pkg/sdk/keycloak"
	"github.com/indrasaputra/arjuna/service/auth/internal/config"
	"github.com/indrasaputra/arjuna/service/auth/internal/grpc/handler"
	"github.com/indrasaputra/arjuna/service/auth/internal/repository/keycloak"
	"github.com/indrasaputra/arjuna/service/auth/internal/service"
)

// Dependency holds any dependency to build full use cases.
type Dependency struct {
	Config         *config.Config
	KeycloakClient kcsdk.Keycloak
}

// BuildAuthHandler builds auth handler including all of its dependencies.
func BuildAuthHandler(dep *Dependency) (*handler.Auth, error) {
	kcConfig := &keycloak.Config{
		Client: dep.KeycloakClient,
		Realm:  dep.Config.Keycloak.Realm,
	}
	kc, err := keycloak.NewAuth(kcConfig)
	if err != nil {
		return nil, err
	}
	auth := service.NewAuth(kc)
	return handler.NewAuth(auth), nil
}

// BuildKeycloakClient builds a keycloak client.
func BuildKeycloakClient(cfg config.Keycloak) kcsdk.Keycloak {
	hc := &http.Client{Timeout: time.Duration(cfg.Timeout) * time.Second}
	client := kcsdk.NewClient(hc, cfg.Address)
	return client
}
