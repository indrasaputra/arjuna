package builder

import (
	sdkpostgres "github.com/indrasaputra/arjuna/pkg/sdk/database/postgres"
	"github.com/indrasaputra/arjuna/pkg/sdk/uow"
	"github.com/indrasaputra/arjuna/service/auth/internal/config"
	"github.com/indrasaputra/arjuna/service/auth/internal/grpc/handler"
	"github.com/indrasaputra/arjuna/service/auth/internal/repository/db"
	"github.com/indrasaputra/arjuna/service/auth/internal/repository/postgres"
	"github.com/indrasaputra/arjuna/service/auth/internal/service"
)

// Dependency holds any dependency to build full use cases.
type Dependency struct {
	Config             *config.Config
	Queries            *db.Queries
	SigningKey         string
	ExpiryTimeInMinute int
}

// BuildAuthHandler builds auth handler including all of its dependencies.
func BuildAuthHandler(dep *Dependency) (*handler.Auth, error) {
	acc := postgres.NewAccount(dep.Queries)
	auth := service.NewAuth(acc, []byte(dep.SigningKey), dep.ExpiryTimeInMinute)
	return handler.NewAuth(auth), nil
}

// BuildQueries builds sqlc queries.
func BuildQueries(tr uow.Tr, getter uow.TxGetter) *db.Queries {
	tx := sdkpostgres.NewTxDB(tr, getter)
	return db.New(tx)
}
