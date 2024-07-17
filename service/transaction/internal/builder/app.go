package builder

import (
	"context"

	goredis "github.com/redis/go-redis/v9"

	sdkpg "github.com/indrasaputra/arjuna/pkg/sdk/database/postgres"
	"github.com/indrasaputra/arjuna/pkg/sdk/uow"
	"github.com/indrasaputra/arjuna/service/transaction/internal/config"
	"github.com/indrasaputra/arjuna/service/transaction/internal/grpc/handler"
	"github.com/indrasaputra/arjuna/service/transaction/internal/repository/postgres"
	"github.com/indrasaputra/arjuna/service/transaction/internal/repository/redis"
	"github.com/indrasaputra/arjuna/service/transaction/internal/service"
)

// Dependency holds any dependency to build full use cases.
type Dependency struct {
	Config      *config.Config
	DB          uow.DB
	RedisClient goredis.Cmdable
}

// BuildTransactionCommandHandler builds transaction command handler including all of its dependencies.
func BuildTransactionCommandHandler(dep *Dependency) *handler.TransactionCommand {
	p := postgres.NewTransaction(dep.DB)
	i := redis.NewIdempotencyKey(dep.RedisClient)

	t := service.NewTransactionCreator(p, i)
	return handler.NewTransactionCommand(t)
}

// BuildBunDB builds BunDB.
func BuildBunDB(cfg sdkpg.Config) (*sdkpg.BunDB, error) {
	pdb, err := sdkpg.NewDBWithPgx(cfg)
	if err != nil {
		return nil, err
	}
	return sdkpg.NewBunDB(pdb)
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
