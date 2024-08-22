package postgres

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"

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

// AddWalletBalance adds some amount to specific user's wallet.
func (w *Wallet) AddWalletBalance(ctx context.Context, id uuid.UUID, amount decimal.Decimal) error {
	return w.addWalletBalance(ctx, w.db, id, amount)
}

// AddWalletBalanceWithTx adds some amount to specific user's wallet.
func (w *Wallet) AddWalletBalanceWithTx(ctx context.Context, tx uow.Tx, id uuid.UUID, amount decimal.Decimal) error {
	if tx == nil {
		app.Logger.Errorf(ctx, "[WalletPostgres-AddWalletBalanceWithTx] internal error: tx is nil")
		return entity.ErrInternal("something went wrong")
	}
	return w.addWalletBalance(ctx, tx, id, amount)
}

// AddWalletBalance adds some amount to specific user's wallet.
func (w *Wallet) addWalletBalance(ctx context.Context, db uow.DB, id uuid.UUID, amount decimal.Decimal) error {
	query := "UPDATE wallets SET balance = balance + ? WHERE id = ?"
	_, err := db.Exec(ctx, query, amount, id)
	if err != nil {
		app.Logger.Errorf(ctx, "[WalletPostgres-addWalletBalance] internal error: %v", err)
		return entity.ErrInternal("something went wrong")
	}
	return nil
}

// GetUserWalletWithTx gets user's wallet from repository.
func (w *Wallet) GetUserWalletWithTx(ctx context.Context, tx uow.Tx, id uuid.UUID, userID uuid.UUID) (*entity.Wallet, error) {
	if tx == nil {
		app.Logger.Errorf(ctx, "[WalletPostgres-GetUserWalletWithTx] internal error: tx is nil")
		return nil, entity.ErrInternal("something went wrong")
	}

	query := "SELECT id, user_id, balance FROM wallets WHERE id = ? AND user_id = ? LIMIT 1 FOR NO KEY UPDATE"
	res := entity.Wallet{}
	err := tx.Query(ctx, &res, query, id, userID)

	if err == sql.ErrNoRows {
		return nil, entity.ErrEmptyWallet()
	}
	if err != nil {
		app.Logger.Errorf(ctx, "[WalletPostgres-GetUserWalletWithTx] internal error: %v", err)
		return nil, entity.ErrInternal("something went wrong")
	}
	return &res, nil
}
