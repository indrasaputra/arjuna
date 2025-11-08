package builder

import (
	sdkpostgres "github.com/indrasaputra/arjuna/pkg/sdk/database/postgres"
	"github.com/indrasaputra/arjuna/pkg/sdk/uow"
	"github.com/indrasaputra/arjuna/service/wallet/internal/config"
	"github.com/indrasaputra/arjuna/service/wallet/internal/grpc/handler"
	"github.com/indrasaputra/arjuna/service/wallet/internal/repository/db"
	"github.com/indrasaputra/arjuna/service/wallet/internal/repository/postgres"
	"github.com/indrasaputra/arjuna/service/wallet/internal/service"
)

// Dependency holds any dependency to build full use cases.
type Dependency struct {
	Config    *config.Config
	TxManager uow.TxManager
	Queries   *db.Queries
}

// BuildWalletCommandHandler builds wallet command handler including all of its dependencies.
func BuildWalletCommandHandler(dep *Dependency) *handler.WalletCommand {
	p := postgres.NewWallet(dep.Queries)
	c := service.NewWalletCreator(p)
	t := service.NewWalletTopup(p)
	f := service.NewWalletTransferer(p, dep.TxManager)
	return handler.NewWalletCommand(c, t, f)
}

// BuildQueries builds sqlc queries.
func BuildQueries(tr uow.Tr, getter uow.TxGetter) *db.Queries {
	tx := sdkpostgres.NewTxDB(tr, getter)
	return db.New(tx)
}
