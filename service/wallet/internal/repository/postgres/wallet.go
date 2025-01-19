package postgres

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/shopspring/decimal"

	sdkpostgres "github.com/indrasaputra/arjuna/pkg/sdk/database/postgres"
	"github.com/indrasaputra/arjuna/service/wallet/entity"
	"github.com/indrasaputra/arjuna/service/wallet/internal/app"
	"github.com/indrasaputra/arjuna/service/wallet/internal/repository/db"
)

// Wallet is responsible to connect wallet entity with wallets table in PostgreSQL.
type Wallet struct {
	queries *db.Queries
}

// NewWallet creates an instance of Wallet.
func NewWallet(q *db.Queries) *Wallet {
	return &Wallet{queries: q}
}

// Insert inserts a wallet to the database.
func (w *Wallet) Insert(ctx context.Context, wallet *entity.Wallet) error {
	if wallet == nil {
		return entity.ErrEmptyWallet()
	}

	param := db.CreateWalletParams{
		ID:        wallet.ID,
		UserID:    wallet.UserID,
		Balance:   wallet.Balance,
		CreatedAt: wallet.CreatedAt,
		UpdatedAt: wallet.UpdatedAt,
		CreatedBy: wallet.CreatedBy,
		UpdatedBy: wallet.UpdatedBy,
	}
	err := w.queries.CreateWallet(ctx, param)

	if sdkpostgres.IsUniqueViolationError(err) {
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
	param := db.AddWalletBalanceParams{ID: id, Amount: amount}
	err := w.queries.AddWalletBalance(ctx, param)
	if err != nil {
		app.Logger.Errorf(ctx, "[WalletPostgres-addWalletBalance] internal error: %v", err)
		return entity.ErrInternal(err.Error())
	}
	return nil
}

// GetUserWalletForUpdate gets user's wallet for update.
func (w *Wallet) GetUserWalletForUpdate(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*entity.Wallet, error) {
	param := db.GetUserWalletForUpdateParams{ID: id, UserID: userID}
	wallet, err := w.queries.GetUserWalletForUpdate(ctx, param)
	if err == pgx.ErrNoRows {
		return nil, entity.ErrEmptyWallet()
	}
	if err != nil {
		app.Logger.Errorf(ctx, "[WalletPostgres-GetUserWalletWithTx] internal error: %v", err)
		return nil, entity.ErrInternal(err.Error())
	}
	return &entity.Wallet{
		ID:      wallet.ID,
		UserID:  wallet.UserID,
		Balance: wallet.Balance,
	}, nil
}
