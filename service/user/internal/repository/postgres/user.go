package postgres

import (
	"context"
	"log"

	"github.com/google/uuid"

	sdkpg "github.com/indrasaputra/arjuna/pkg/sdk/database/postgres"
	"github.com/indrasaputra/arjuna/pkg/sdk/uow"
	"github.com/indrasaputra/arjuna/service/user/entity"
	"github.com/indrasaputra/arjuna/service/user/internal/app"
)

// User is responsible to connect user entity with users table in PostgreSQL.
type User struct {
	db     uow.Tr
	getter uow.TxGetter
}

// NewUser creates an instance of User.
func NewUser(db uow.Tr, g uow.TxGetter) *User {
	return &User{db: db, getter: g}
}

// Insert inserts the user into users table.
func (u *User) Insert(ctx context.Context, user *entity.User) error {
	if user == nil {
		return entity.ErrEmptyUser()
	}

	tx := u.getter.DefaultTrOrDB(ctx, u.db)
	query := "INSERT INTO " +
		"users (id, name, created_at, updated_at, created_by, updated_by) " +
		"VALUES ($1, $2, $3, $4, $5, $6)"

	_, err := tx.Exec(ctx, query,
		user.ID,
		user.Name,
		user.CreatedAt,
		user.UpdatedAt,
		user.CreatedBy,
		user.UpdatedBy,
	)

	if err == sdkpg.ErrAlreadyExist {
		return entity.ErrAlreadyExists()
	}
	if err != nil {
		app.Logger.Errorf(ctx, "[PostgresUser-Insert] fail insert user: %v", err)
		return entity.ErrInternal(err.Error())
	}
	return nil
}

// GetByID gets a user from database.
// It returns entity.ErrNotFound if user can't be found.
func (u *User) GetByID(ctx context.Context, id uuid.UUID) (*entity.User, error) {
	tx := u.getter.DefaultTrOrDB(ctx, u.db)

	query := "SELECT id, name, created_at, updated_at, created_by, updated_by FROM users WHERE id = $1 LIMIT 1"
	var res entity.User
	row := tx.QueryRow(ctx, query, id)
	err := row.Scan(&res.ID, &res.Name, &res.CreatedAt, &res.UpdatedAt, &res.CreatedBy, &res.UpdatedBy)
	if err != nil {
		app.Logger.Errorf(ctx, "[PostgresUser-GetByID] fail get user: %v", err)
		return nil, entity.ErrInternal(err.Error())
	}
	return &res, nil
}

// GetAll gets all users in users table.
func (u *User) GetAll(ctx context.Context, limit uint) ([]*entity.User, error) {
	tx := u.getter.DefaultTrOrDB(ctx, u.db)

	query := "SELECT id, name, created_at, updated_at, created_by, updated_by FROM users LIMIT $1"
	res := []*entity.User{}
	rows, err := tx.Query(ctx, query, limit)
	if err != nil {
		app.Logger.Errorf(ctx, "[PostgresUser-GetAll] fail get all users: %v", err)
		return []*entity.User{}, entity.ErrInternal(err.Error())
	}
	defer rows.Close()

	for rows.Next() {
		var tmp entity.User
		if err := rows.Scan(&tmp.ID, &tmp.Name, &tmp.CreatedAt, &tmp.UpdatedAt, &tmp.CreatedBy, &tmp.UpdatedBy); err != nil {
			log.Printf("[PostgresUser-GetAll] scan rows error: %s", err.Error())
			continue
		}
		res = append(res, &tmp)
	}

	return res, nil
}

// HardDelete deletes a user from database.
// If the user doesn't exist, it doesn't returns error.
func (u *User) HardDelete(ctx context.Context, id uuid.UUID) error {
	tx := u.getter.DefaultTrOrDB(ctx, u.db)
	query := "DELETE FROM users WHERE id = $1"
	_, err := tx.Exec(ctx, query, id)
	if err != nil {
		app.Logger.Errorf(ctx, "[PostgresUser-doHardDelete] fail delete user: %v", err)
		return entity.ErrInternal(err.Error())
	}
	return nil
}
