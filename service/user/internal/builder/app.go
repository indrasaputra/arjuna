package builder

import (
	"go.temporal.io/sdk/client"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	sdkpostgres "github.com/indrasaputra/arjuna/pkg/sdk/database/postgres"
	"github.com/indrasaputra/arjuna/pkg/sdk/uow"
	sdkauth "github.com/indrasaputra/arjuna/service/auth/pkg/sdk/auth"
	"github.com/indrasaputra/arjuna/service/user/internal/config"
	"github.com/indrasaputra/arjuna/service/user/internal/grpc/handler"
	"github.com/indrasaputra/arjuna/service/user/internal/repository/db"
	"github.com/indrasaputra/arjuna/service/user/internal/repository/postgres"
	"github.com/indrasaputra/arjuna/service/user/internal/service"
	sdkwallet "github.com/indrasaputra/arjuna/service/wallet/pkg/sdk/wallet"
)

// Dependency holds any dependency to build full use cases.
type Dependency struct {
	Config         *config.Config
	TemporalClient client.Client
	TxManager      uow.TxManager
	Queries        *db.Queries
}

// BuildUserCommandHandler builds user command handler including all of its dependencies.
func BuildUserCommandHandler(dep *Dependency) *handler.UserCommand {
	pu := postgres.NewUser(dep.Queries)
	puo := postgres.NewUserOutbox(dep.Queries)

	rg := service.NewUserRegistrar(dep.TxManager, pu, puo)
	return handler.NewUserCommand(rg)
}

// BuildUserCommandInternalHandler builds user command handler including all of its dependencies.
func BuildUserCommandInternalHandler(dep *Dependency) *handler.UserCommandInternal {
	pg := postgres.NewUser(dep.Queries)
	d := service.NewUserDeleter(pg)
	return handler.NewUserCommandInternal(d)
}

// BuildUserQueryHandler builds user query handler including all of its dependencies.
func BuildUserQueryHandler(dep *Dependency) *handler.UserQuery {
	pg := postgres.NewUser(dep.Queries)
	g := service.NewUserGetter(pg)
	return handler.NewUserQuery(g)
}

// BuildTemporalClient builds temporal client.
func BuildTemporalClient(address string) (client.Client, error) {
	return client.Dial(client.Options{HostPort: address})
}

// BuildAuthClient builds auth service client.
func BuildAuthClient(host, username, password string) (*sdkauth.Client, error) {
	dc := &sdkauth.Config{
		Host:     host,
		Options:  []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())},
		Username: username,
		Password: password,
	}
	return sdkauth.NewClient(dc)
}

// BuildWalletClient builds wallet service client.
func BuildWalletClient(host, username, password string) (*sdkwallet.Client, error) {
	dc := &sdkwallet.Config{
		Host:     host,
		Options:  []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())},
		Username: username,
		Password: password,
	}
	return sdkwallet.NewClient(dc)
}

// BuildQueries builds sqlc queries.
func BuildQueries(tr uow.Tr, getter uow.TxGetter) *db.Queries {
	tx := sdkpostgres.NewTxDB(tr, getter)
	return db.New(tx)
}
