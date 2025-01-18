package builder

import (
	"context"

	goredis "github.com/redis/go-redis/v9"

	"github.com/indrasaputra/arjuna/pkg/sdk/uow"
	"github.com/indrasaputra/arjuna/service/wallet/internal/config"
	"github.com/indrasaputra/arjuna/service/wallet/internal/grpc/handler"
	"github.com/indrasaputra/arjuna/service/wallet/internal/repository/db"
	"github.com/indrasaputra/arjuna/service/wallet/internal/repository/postgres"
	"github.com/indrasaputra/arjuna/service/wallet/internal/repository/redis"
	"github.com/indrasaputra/arjuna/service/wallet/internal/service"
)

// Dependency holds any dependency to build full use cases.
type Dependency struct {
	Config      *config.Config
	RedisClient goredis.Cmdable
	TxManager   uow.TxManager
	Queries     *db.Queries
}

// BuildWalletCommandHandler builds wallet command handler including all of its dependencies.
func BuildWalletCommandHandler(dep *Dependency) *handler.WalletCommand {
	p := postgres.NewWallet(dep.Queries)
	k := redis.NewIdempotencyKey(dep.RedisClient)
	c := service.NewWalletCreator(p)
	t := service.NewWalletTopup(p, k)
	f := service.NewWalletTransferer(p, dep.TxManager)
	return handler.NewWalletCommand(c, t, f)
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
