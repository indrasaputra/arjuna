package postgres

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/indrasaputra/arjuna/pkg/sdk/uow"
	"github.com/indrasaputra/arjuna/service/user/entity"
	"github.com/indrasaputra/arjuna/service/user/internal/app"
	"github.com/indrasaputra/arjuna/service/user/internal/repository/db"
)

// User is responsible to connect user entity with users table in PostgreSQL.
type User struct {
	db      uow.Tr
	getter  uow.TxGetter
	queries *db.Queries
}

// NewUser creates an instance of User.
func NewUser(tr uow.Tr, g uow.TxGetter) *User {
	tx := uow.NewTxDB(tr, g)
	q := db.New(tx)
	return &User{db: tr, getter: g, queries: q}
}

// Insert inserts the user into users table.
func (u *User) Insert(ctx context.Context, user *entity.User) error {
	if user == nil {
		return entity.ErrEmptyUser()
	}

	param := db.CreateUserParams{
		ID:        user.ID,
		Name:      user.Name,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		CreatedBy: user.CreatedBy,
		UpdatedBy: user.UpdatedBy,
	}
	err := u.queries.CreateUser(ctx, param)

	if uow.IsUniqueViolationError(err) {
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
	if err == pgx.ErrNoRows {
		return nil, entity.ErrNotFound()
	}
	if err != nil {
		app.Logger.Errorf(ctx, "[PostgresUser-GetByID] fail get user: %v", err)
		return nil, entity.ErrInternal(err.Error())
	}
	return &res, nil
}

// GetAll gets all users in users table.
func (u *User) GetAll(ctx context.Context, limit uint) ([]*entity.User, error) {
	users, err := u.queries.GetAllUsers(ctx, int32(limit))
	if err != nil {
		app.Logger.Errorf(ctx, "[PostgresUser-GetAll] fail get all users: %v", err)
		return []*entity.User{}, entity.ErrInternal(err.Error())
	}

	result := make([]*entity.User, 0, len(users))
	for _, user := range users {
		res := &entity.User{}
		res.ID = user.ID
		res.Name = user.Name
		res.CreatedAt = user.CreatedAt
		res.UpdatedAt = user.UpdatedAt
		res.CreatedBy = user.CreatedBy
		res.UpdatedBy = user.UpdatedBy
		result = append(result, res)
	}
	return result, err
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
