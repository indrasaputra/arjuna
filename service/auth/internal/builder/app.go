package builder

import (
	"net/http"
	"time"

	sdkpg "github.com/indrasaputra/arjuna/pkg/sdk/database/postgres"
	kcsdk "github.com/indrasaputra/arjuna/pkg/sdk/keycloak"
	"github.com/indrasaputra/arjuna/pkg/sdk/uow"
	"github.com/indrasaputra/arjuna/service/auth/internal/config"
	"github.com/indrasaputra/arjuna/service/auth/internal/grpc/handler"
	"github.com/indrasaputra/arjuna/service/auth/internal/repository/postgres"
	"github.com/indrasaputra/arjuna/service/auth/internal/service"
)

// Dependency holds any dependency to build full use cases.
type Dependency struct {
	Config         *config.Config
	KeycloakClient kcsdk.Keycloak
	DB             uow.DB
}

// BuildAuthHandler builds auth handler including all of its dependencies.
func BuildAuthHandler(dep *Dependency) (*handler.Auth, error) {
	acc := postgres.NewAccount(dep.DB)
	auth := service.NewAuth(acc)
	return handler.NewAuth(auth), nil
}

// BuildKeycloakClient builds a keycloak client.
func BuildKeycloakClient(cfg config.Keycloak) kcsdk.Keycloak {
	hc := &http.Client{Timeout: time.Duration(cfg.Timeout) * time.Second}
	client := kcsdk.NewClient(hc, cfg.Address)
	return client
}

// BuildBunDB builds BunDB.
func BuildBunDB(cfg sdkpg.Config) (*sdkpg.BunDB, error) {
	pdb, err := sdkpg.NewDBWithPgx(cfg)
	if err != nil {
		return nil, err
	}
	return sdkpg.NewBunDB(pdb)
}
