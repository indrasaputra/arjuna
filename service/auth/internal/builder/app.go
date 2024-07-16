package builder

import (
	sdkpg "github.com/indrasaputra/arjuna/pkg/sdk/database/postgres"
	"github.com/indrasaputra/arjuna/pkg/sdk/uow"
	"github.com/indrasaputra/arjuna/service/auth/internal/config"
	"github.com/indrasaputra/arjuna/service/auth/internal/grpc/handler"
	"github.com/indrasaputra/arjuna/service/auth/internal/repository/postgres"
	"github.com/indrasaputra/arjuna/service/auth/internal/service"
)

// Dependency holds any dependency to build full use cases.
type Dependency struct {
	Config             *config.Config
	DB                 uow.DB
	SigningKey         string
	ExpiryTimeInMinute int
}

// BuildAuthHandler builds auth handler including all of its dependencies.
func BuildAuthHandler(dep *Dependency) (*handler.Auth, error) {
	acc := postgres.NewAccount(dep.DB)
	auth := service.NewAuth(acc, []byte(dep.SigningKey), dep.ExpiryTimeInMinute)
	return handler.NewAuth(auth), nil
}

// BuildBunDB builds BunDB.
func BuildBunDB(cfg sdkpg.Config) (*sdkpg.BunDB, error) {
	pdb, err := sdkpg.NewDBWithPgx(cfg)
	if err != nil {
		return nil, err
	}
	return sdkpg.NewBunDB(pdb)
}
