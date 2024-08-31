package postgres

import (
	"context"
	"database/sql"

	sdkpg "github.com/indrasaputra/arjuna/pkg/sdk/database/postgres"
	"github.com/indrasaputra/arjuna/pkg/sdk/uow"
	"github.com/indrasaputra/arjuna/service/auth/entity"
	"github.com/indrasaputra/arjuna/service/auth/internal/app"
)

// Account is responsible to connect account entity with accounts table in PostgreSQL.
type Account struct {
	db uow.DB
}

// NewAccount creates an instance of Account.
func NewAccount(db uow.DB) *Account {
	return &Account{db: db}
}

// Insert inserts an account to the database.
func (a *Account) Insert(ctx context.Context, account *entity.Account) error {
	if account == nil {
		return entity.ErrEmptyAccount()
	}

	query := "INSERT INTO " +
		"accounts (id, user_id, email, password, created_at, updated_at, created_by, updated_by) " +
		"VALUES (?, ?, ?, ?, ?, ?, ?, ?)"

	_, err := a.db.Exec(ctx, query,
		account.ID,
		account.UserID,
		account.Email,
		account.Password,
		account.CreatedAt,
		account.UpdatedAt,
		account.CreatedBy,
		account.UpdatedBy,
	)

	if err == sdkpg.ErrAlreadyExist {
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
	query := "SELECT id, user_id, email, password FROM accounts WHERE email = ? LIMIT 1"
	var res entity.Account
	err := a.db.Query(ctx, &res, query, email)
	if err == sql.ErrNoRows {
		return nil, entity.ErrNotFound()
	}
	if err != nil {
		app.Logger.Errorf(ctx, "[PostgresAccount-GetByEmail] fail get account: %v", err)
		return nil, entity.ErrInternal(err.Error())
	}
	return &res, nil
}
