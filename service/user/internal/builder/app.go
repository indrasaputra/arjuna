package builder

import (
	"go.temporal.io/sdk/client"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	sdkpg "github.com/indrasaputra/arjuna/pkg/sdk/database/postgres"
	"github.com/indrasaputra/arjuna/pkg/sdk/uow"
	sdkauth "github.com/indrasaputra/arjuna/service/auth/pkg/sdk/auth"
	"github.com/indrasaputra/arjuna/service/user/internal/config"
	"github.com/indrasaputra/arjuna/service/user/internal/grpc/handler"
	"github.com/indrasaputra/arjuna/service/user/internal/repository/postgres"
	"github.com/indrasaputra/arjuna/service/user/internal/service"
)

// Dependency holds any dependency to build full use cases.
type Dependency struct {
	Config         *config.Config
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
func BuildUserCommandInternalHandler(dep *Dependency) *handler.UserCommandInternal {
	pg := postgres.NewUser(dep.DB)
	u := uow.NewUnitWorker(dep.DB)
	d := service.NewUserDeleter(u, pg)
	return handler.NewUserCommandInternal(d)
}

// BuildUserQueryHandler builds user query handler including all of its dependencies.
func BuildUserQueryHandler(dep *Dependency) *handler.UserQuery {
	pg := postgres.NewUser(dep.DB)
	g := service.NewUserGetter(pg)
	return handler.NewUserQuery(g)
}

// BuildBunDB builds BunDB.
func BuildBunDB(cfg sdkpg.Config) (*sdkpg.BunDB, error) {
	pdb, err := sdkpg.NewDBWithPgx(cfg)
	if err != nil {
		return nil, err
	}
	return sdkpg.NewBunDB(pdb)
}

// BuildTemporalClient builds temporal client.
func BuildTemporalClient(address string) (client.Client, error) {
	return client.Dial(client.Options{HostPort: address})
}

// BuildAuthClient builds auth service client.
func BuildAuthClient(host string) (*sdkauth.Client, error) {
	dc := &sdkauth.DialConfig{
		Host:    host,
		Options: []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())},
	}
	return sdkauth.NewClient(dc)
}
