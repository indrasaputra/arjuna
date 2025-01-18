package builder

import (
	"context"

	goredis "github.com/redis/go-redis/v9"
	"go.temporal.io/sdk/client"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/indrasaputra/arjuna/pkg/sdk/uow"
	sdkauth "github.com/indrasaputra/arjuna/service/auth/pkg/sdk/auth"
	"github.com/indrasaputra/arjuna/service/user/internal/config"
	"github.com/indrasaputra/arjuna/service/user/internal/grpc/handler"
	"github.com/indrasaputra/arjuna/service/user/internal/repository/db"
	"github.com/indrasaputra/arjuna/service/user/internal/repository/postgres"
	"github.com/indrasaputra/arjuna/service/user/internal/repository/redis"
	"github.com/indrasaputra/arjuna/service/user/internal/service"
	sdkwallet "github.com/indrasaputra/arjuna/service/wallet/pkg/sdk/wallet"
)

// Dependency holds any dependency to build full use cases.
type Dependency struct {
	Config         *config.Config
	TemporalClient client.Client
	RedisClient    goredis.Cmdable
	TxManager      uow.TxManager
	Queries        *db.Queries
}

// BuildUserCommandHandler builds user command handler including all of its dependencies.
func BuildUserCommandHandler(dep *Dependency) *handler.UserCommand {
	pu := postgres.NewUser(dep.Queries)
	puo := postgres.NewUserOutbox(dep.Queries)
	ik := redis.NewIdempotencyKey(dep.RedisClient)

	rg := service.NewUserRegistrar(dep.TxManager, pu, puo, ik)
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

// BuildRedisClient builds an instance of redis client.
func BuildRedisClient(cfg *config.Redis) (*goredis.Client, error) {
	opt := &goredis.Options{
		Addr: cfg.Address,
	}

	client := goredis.NewClient(opt)
	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		return nil, err
	}

	return client, nil
}

// BuildQueries builds sqlc queries.
func BuildQueries(tr uow.Tr, getter uow.TxGetter) *db.Queries {
	tx := uow.NewTxDB(tr, getter)
	return db.New(tx)
}
