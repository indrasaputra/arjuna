package builder

import (
	"context"

	goredis "github.com/redis/go-redis/v9"

	"github.com/indrasaputra/arjuna/pkg/sdk/uow"
	"github.com/indrasaputra/arjuna/service/transaction/internal/config"
	"github.com/indrasaputra/arjuna/service/transaction/internal/grpc/handler"
	"github.com/indrasaputra/arjuna/service/transaction/internal/repository/db"
	"github.com/indrasaputra/arjuna/service/transaction/internal/repository/postgres"
	"github.com/indrasaputra/arjuna/service/transaction/internal/repository/redis"
	"github.com/indrasaputra/arjuna/service/transaction/internal/service"
)

// Dependency holds any dependency to build full use cases.
type Dependency struct {
	Config      *config.Config
	RedisClient goredis.Cmdable
	Queries     *db.Queries
}

// BuildTransactionCommandHandler builds transaction command handler including all of its dependencies.
func BuildTransactionCommandHandler(dep *Dependency) *handler.TransactionCommand {
	p := postgres.NewTransaction(dep.Queries)
	i := redis.NewIdempotencyKey(dep.RedisClient)

	c := service.NewTransactionCreator(p, i)

	return handler.NewTransactionCommand(c)
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
