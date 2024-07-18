package builder

import (
	"context"

	goredis "github.com/redis/go-redis/v9"

	sdkpg "github.com/indrasaputra/arjuna/pkg/sdk/database/postgres"
	"github.com/indrasaputra/arjuna/pkg/sdk/uow"
	"github.com/indrasaputra/arjuna/service/wallet/internal/config"
	"github.com/indrasaputra/arjuna/service/wallet/internal/grpc/handler"
	"github.com/indrasaputra/arjuna/service/wallet/internal/repository/postgres"
	"github.com/indrasaputra/arjuna/service/wallet/internal/service"
)

// Dependency holds any dependency to build full use cases.
type Dependency struct {
	Config      *config.Config
	DB          uow.DB
	RedisClient goredis.Cmdable
}

// BuildWalletCommandHandler builds wallet command handler including all of its dependencies.
func BuildWalletCommandHandler(dep *Dependency) *handler.WalletCommand {
	p := postgres.NewWallet(dep.DB)
	t := service.NewWalletCreator(p)
	return handler.NewWalletCommand(t)
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
