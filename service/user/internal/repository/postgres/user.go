package postgres

import (
	"context"
	"errors"

	"github.com/google/uuid"

	sdkpostgres "github.com/indrasaputra/arjuna/pkg/sdk/database/postgres"
	"github.com/indrasaputra/arjuna/service/user/entity"
	"github.com/indrasaputra/arjuna/service/user/internal/app"
	"github.com/indrasaputra/arjuna/service/user/internal/repository/db"
)

// User is responsible to connect user entity with users table in PostgreSQL.
type User struct {
	queries *db.Queries
}

// NewUser creates an instance of User.
func NewUser(q *db.Queries) *User {
	return &User{queries: q}
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

	if sdkpostgres.IsUniqueViolationError(err) {
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
	user, err := u.queries.GetUserByID(ctx, id)
	if errors.Is(err, sdkpostgres.ErrNotFound) {
		return nil, entity.ErrNotFound()
	}
	if err != nil {
		app.Logger.Errorf(ctx, "[PostgresUser-GetByID] fail get user: %v", err)
		return nil, entity.ErrInternal(err.Error())
	}

	return &entity.User{
		ID:   user.ID,
		Name: user.Name,
		Auditable: entity.Auditable{
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			CreatedBy: user.CreatedBy,
			UpdatedBy: user.UpdatedBy,
		},
	}, nil
}

// GetAll gets all users in users table.
func (u *User) GetAll(ctx context.Context, limit uint) ([]*entity.User, error) {
	users, err := u.queries.GetAllUsers(ctx, int32(limit))
	if err != nil {
		app.Logger.Errorf(ctx, "[PostgresUser-GetAll] fail get all users: %v", err)
		return []*entity.User{}, entity.ErrInternal(err.Error())
	}

	result := make([]*entity.User, len(users))
	for i, user := range users {
		res := &entity.User{}
		res.ID = user.ID
		res.Name = user.Name
		res.CreatedAt = user.CreatedAt
		res.UpdatedAt = user.UpdatedAt
		res.CreatedBy = user.CreatedBy
		res.UpdatedBy = user.UpdatedBy
		result[i] = res
	}
	return result, err
}

// HardDelete deletes a user from database.
// If the user doesn't exist, it doesn't returns error.
func (u *User) HardDelete(ctx context.Context, id uuid.UUID) error {
	if err := u.queries.HardDeleteUserByID(ctx, id); err != nil {
		app.Logger.Errorf(ctx, "[PostgresUser-doHardDelete] fail delete user: %v", err)
		return entity.ErrInternal(err.Error())
	}
	return nil
}
