package builder

import (
	"net/http"
	"time"

	"go.temporal.io/sdk/client"

	pgsdk "github.com/indrasaputra/arjuna/pkg/sdk/database/postgres"
	kcsdk "github.com/indrasaputra/arjuna/pkg/sdk/keycloak"
	"github.com/indrasaputra/arjuna/pkg/sdk/uow"
	"github.com/indrasaputra/arjuna/service/user/internal/config"
	"github.com/indrasaputra/arjuna/service/user/internal/grpc/handler"
	"github.com/indrasaputra/arjuna/service/user/internal/repository/keycloak"
	"github.com/indrasaputra/arjuna/service/user/internal/repository/postgres"
	"github.com/indrasaputra/arjuna/service/user/internal/service"
)

// Dependency holds any dependency to build full use cases.
type Dependency struct {
	Config         *config.Config
	KeycloakClient kcsdk.Keycloak
	TemporalClient client.Client
	DB             uow.DB
}

// BuildUserCommandHandler builds user command handler including all of its dependencies.
func BuildUserCommandHandler(dep *Dependency) *handler.UserCommand {
	pu := postgres.NewUser(dep.DB)
	puo := postgres.NewUserOutbox(dep.DB)
	u := uow.NewUnitWorker(dep.DB)

	rg := service.NewUserRegistrar(pu, puo, u)
	return handler.NewUserCommand(rg)
}

// BuildUserCommandInternalHandler builds user command handler including all of its dependencies.
func BuildUserCommandInternalHandler(dep *Dependency) (*handler.UserCommandInternal, error) {
	kcconf := &keycloak.Config{
		Client:        dep.KeycloakClient,
		Realm:         dep.Config.Keycloak.Realm,
		AdminUsername: dep.Config.Keycloak.AdminUser,
		AdminPassword: dep.Config.Keycloak.AdminPassword,
	}
	kc, err := keycloak.NewUser(kcconf)
	if err != nil {
		return nil, err
	}

	pg := postgres.NewUser(dep.DB)
	u := uow.NewUnitWorker(dep.DB)
	d := service.NewUserDeleter(u, pg, kc)
	return handler.NewUserCommandInternal(d), nil
}

// BuildUserQueryHandler builds user query handler including all of its dependencies.
func BuildUserQueryHandler(dep *Dependency) *handler.UserQuery {
	pg := postgres.NewUser(dep.DB)
	g := service.NewUserGetter(pg)
	return handler.NewUserQuery(g)
}

// BuildBunDB builds BunDB.
func BuildBunDB(cfg pgsdk.Config) (*pgsdk.BunDB, error) {
	pdb, err := pgsdk.NewDBWithPgx(cfg)
	if err != nil {
		return nil, err
	}
	return pgsdk.NewBunDB(pdb)
}

// BuildKeycloakClient builds a keycloak client.
func BuildKeycloakClient(cfg config.Keycloak) kcsdk.Keycloak {
	hc := &http.Client{Timeout: time.Duration(cfg.Timeout) * time.Second}
	client := kcsdk.NewClient(hc, cfg.Address)
	return client
}

// BuildTemporalClient builds temporal client.
func BuildTemporalClient(address string) (client.Client, error) {
	return client.Dial(client.Options{HostPort: address})
}
