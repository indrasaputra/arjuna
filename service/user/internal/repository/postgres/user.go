package postgres

import (
	"context"

	pgsdk "github.com/indrasaputra/arjuna/pkg/sdk/database/postgres"
	"github.com/indrasaputra/arjuna/pkg/sdk/uow"
	"github.com/indrasaputra/arjuna/service/user/entity"
	"github.com/indrasaputra/arjuna/service/user/internal/app"
)

// User is responsible to connect user entity with users table in PostgreSQL.
type User struct {
	db uow.DB
}

// NewUser creates an instance of User.
func NewUser(db uow.DB) *User {
	return &User{db: db}
}

// InsertWithTx inserts the user into users table using transaction.
func (u *User) InsertWithTx(ctx context.Context, tx uow.Tx, user *entity.User) error {
	if tx == nil {
		app.Logger.Errorf(ctx, "[PostgresUser-InsertWithTx] transaction is not set")
		return entity.ErrInternal("transaction is not set")
	}
	if user == nil {
		return entity.ErrEmptyUser()
	}

	query := "INSERT INTO " +
		"users (id, keycloak_id, name, email, created_at, updated_at, created_by, updated_by) " +
		"VALUES (?, ?, ?, ?, ?, ?, ?, ?)"

	_, err := tx.Exec(ctx, query,
		user.ID,
		user.ID, // Set Keycloak ID as ID to maintenance uniqueness
		user.Name,
		user.Email,
		user.CreatedAt,
		user.UpdatedAt,
		user.CreatedBy,
		user.UpdatedBy,
	)

	if err == pgsdk.ErrAlreadyExist {
		return entity.ErrAlreadyExists()
	}
	if err != nil {
		app.Logger.Errorf(ctx, "[PostgresUser-InsertWithTx] fail insert user with tx: %v", err)
		return entity.ErrInternal(err.Error())
	}
	return nil
}

// GetByID gets a user from database.
// It returns entity.ErrNotFound if user can't be found.
func (u *User) GetByID(ctx context.Context, id string) (*entity.User, error) {
	query := "SELECT id, keycloak_id, name, email, created_at, updated_at, created_by, updated_by FROM users WHERE id = ? LIMIT 1"
	var res entity.User
	err := u.db.Query(ctx, &res, query, id)
	if err != nil {
		app.Logger.Errorf(ctx, "[PostgresUser-GetByID] fail get user: %v", err)
		return nil, entity.ErrInternal(err.Error())
	}
	return &res, nil
}

// GetAll gets all users in users table.
func (u *User) GetAll(ctx context.Context, limit uint) ([]*entity.User, error) {
	query := "SELECT id, keycloak_id, name, email, created_at, updated_at, created_by, updated_by FROM users LIMIT ?"
	res := []*entity.User{}
	err := u.db.Query(ctx, &res, query, limit)
	if err != nil {
		app.Logger.Errorf(ctx, "[PostgresUser-GetAll] fail get all users: %v", err)
		return []*entity.User{}, entity.ErrInternal(err.Error())
	}
	return res, nil
}

// UpdateKeycloakID updates user's keycloak id in database.
func (u *User) UpdateKeycloakID(ctx context.Context, id, keycloakID string) error {
	query := "UPDATE users SET keycloak_id = ? WHERE id = ?"
	_, err := u.db.Exec(ctx, query, keycloakID, id)
	if err != nil {
		app.Logger.Errorf(ctx, "[PostgresUser-UpdateKeycloakID] fail update user's keycloak id: %v", err)
		return entity.ErrInternal(err.Error())
	}
	return nil
}

// HardDeleteWithTx deletes a user from database.
// If the user doesn't exist, it doesn't returns error.
func (u *User) HardDeleteWithTx(ctx context.Context, tx uow.Tx, id string) error {
	if tx == nil {
		app.Logger.Errorf(ctx, "[PostgresUser-HardDeleteWithTx] transaction is not set")
		return entity.ErrInternal("transaction is not set")
	}

	query := "DELETE FROM users WHERE id = ?"
	_, err := tx.Exec(ctx, query, id)
	if err != nil {
		app.Logger.Errorf(ctx, "[PostgresUser-HardDeleteWithTx] fail delete user: %v", err)
		return entity.ErrInternal(err.Error())
	}
	return nil
}
