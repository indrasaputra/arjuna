package postgres

import (
	"context"

	sdkpg "github.com/indrasaputra/arjuna/pkg/sdk/database/postgres"
	"github.com/indrasaputra/arjuna/pkg/sdk/uow"
	"github.com/indrasaputra/arjuna/service/wallet/entity"
	"github.com/indrasaputra/arjuna/service/wallet/internal/app"
)

// Wallet is responsible to connect wallet entity with wallets table in PostgreSQL.
type Wallet struct {
	db uow.DB
}

// NewWallet creates an instance of Wallet.
func NewWallet(db uow.DB) *Wallet {
	return &Wallet{db: db}
}

// Insert inserts a wallet to the database.
func (w *Wallet) Insert(ctx context.Context, trx *entity.Wallet) error {
	if trx == nil {
		return entity.ErrEmptyWallet()
	}

	query := "INSERT INTO " +
		"wallets (id, user_id, balance, created_at, updated_at, created_by, updated_by) " +
		"VALUES (?, ?, ?, ?, ?, ?, ?)"

	_, err := w.db.Exec(ctx, query,
		trx.ID,
		trx.UserID,
		trx.Balance,
		trx.CreatedAt,
		trx.UpdatedAt,
		trx.CreatedBy,
		trx.UpdatedBy,
	)

	if err == sdkpg.ErrAlreadyExist {
		return entity.ErrAlreadyExists()
	}
	if err != nil {
		app.Logger.Errorf(ctx, "[PostgresWallet-Insert] fail insert wallet with tx: %v", err)
		return entity.ErrInternal(err.Error())
	}
	return nil
}
