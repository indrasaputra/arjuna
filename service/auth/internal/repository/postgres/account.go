package postgres

import (
	"context"

	"github.com/jackc/pgx/v5"

	sdkpostgres "github.com/indrasaputra/arjuna/pkg/sdk/database/postgres"
	"github.com/indrasaputra/arjuna/service/auth/entity"
	"github.com/indrasaputra/arjuna/service/auth/internal/app"
	"github.com/indrasaputra/arjuna/service/auth/internal/repository/db"
)

// Account is responsible to connect account entity with accounts table in PostgreSQL.
type Account struct {
	queries *db.Queries
}

// NewAccount creates an instance of Account.
func NewAccount(q *db.Queries) *Account {
	return &Account{queries: q}
}

// Insert inserts an account to the database.
func (a *Account) Insert(ctx context.Context, account *entity.Account) error {
	if account == nil {
		return entity.ErrEmptyAccount()
	}

	param := db.CreateAccountParams{
		ID:        account.ID,
		UserID:    account.UserID,
		Email:     account.Email,
		Password:  account.Password,
		CreatedAt: account.CreatedAt,
		UpdatedAt: account.UpdatedAt,
		CreatedBy: account.CreatedBy,
		UpdatedBy: account.UpdatedBy,
	}
	err := a.queries.CreateAccount(ctx, param)
	if sdkpostgres.IsUniqueViolationError(err) {
		return entity.ErrAlreadyExists()
	}
	if err != nil {
		app.Logger.Errorf(ctx, "[PostgresAccount-Insert] fail insert account with tx: %v", err)
		return entity.ErrInternal(err.Error())
	}
	return nil
}

// GetByEmail gets an account by email.
func (a *Account) GetByEmail(ctx context.Context, email string) (*entity.Account, error) {
	account, err := a.queries.GetAccountByEmail(ctx, email)
	if err == pgx.ErrNoRows {
		return nil, entity.ErrNotFound()
	}
	if err != nil {
		app.Logger.Errorf(ctx, "[PostgresAccount-GetByEmail] fail get account: %v", err)
		return nil, entity.ErrInternal(err.Error())
	}

	return &entity.Account{
		ID:       account.ID,
		UserID:   account.UserID,
		Email:    account.Email,
		Password: account.Password,
	}, nil
}
