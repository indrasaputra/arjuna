package builder

import (
	sdkpostgres "github.com/indrasaputra/arjuna/pkg/sdk/database/postgres"
	"github.com/indrasaputra/arjuna/pkg/sdk/uow"
	"github.com/indrasaputra/arjuna/service/transaction/internal/config"
	"github.com/indrasaputra/arjuna/service/transaction/internal/grpc/handler"
	"github.com/indrasaputra/arjuna/service/transaction/internal/repository/db"
	"github.com/indrasaputra/arjuna/service/transaction/internal/repository/postgres"
	"github.com/indrasaputra/arjuna/service/transaction/internal/service"
)

// Dependency holds any dependency to build full use cases.
type Dependency struct {
	Config  *config.Config
	Queries *db.Queries
}

// BuildTransactionCommandHandler builds transaction command handler including all of its dependencies.
func BuildTransactionCommandHandler(dep *Dependency) *handler.TransactionCommand {
	p := postgres.NewTransaction(dep.Queries)

	c := service.NewTransactionCreator(p)

	return handler.NewTransactionCommand(c)
}

// BuildQueries builds sqlc queries.
func BuildQueries(tr uow.Tr, getter uow.TxGetter) *db.Queries {
	tx := sdkpostgres.NewTxDB(tr, getter)
	return db.New(tx)
}
