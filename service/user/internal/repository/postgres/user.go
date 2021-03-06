package postgres

import (
	"context"
	"time"

	"github.com/jackc/pgconn"

	"github.com/indrasaputra/arjuna/service/user/entity"
)

const (
	// errCodeUniqueViolation is derived from https://www.postgresql.org/docs/11/errcodes-appendix.html
	errCodeUniqueViolation = "23505"
)

// User is responsible to connect user entity with users table in PostgreSQL.
type User struct {
	pool PgxPoolIface
}

// NewUser creates an instance of User.
func NewUser(pool PgxPoolIface) *User {
	return &User{pool: pool}
}

// Insert inserts the user into the users table.
func (u *User) Insert(ctx context.Context, user *entity.User) error {
	if user == nil {
		return entity.ErrEmptyUser()
	}
	user.CreatedAt = time.Now().UTC()
	user.UpdatedAt = time.Now().UTC()

	query := "INSERT INTO " +
		"users (id, name, email, password, created_at, updated_at) " +
		"VALUES ($1, $2, $3, $4, $5, $6)"

	_, err := u.pool.Exec(ctx, query,
		user.ID,
		user.Name,
		user.Email,
		user.Password,
		user.CreatedAt,
		user.UpdatedAt,
	)

	if err != nil && isUniqueViolationErr(err) {
		return entity.ErrAlreadyExists()
	}
	if err != nil {
		return entity.ErrInternal(err.Error())
	}
	return nil
}

func isUniqueViolationErr(err error) bool {
	pgerr, ok := err.(*pgconn.PgError)
	if !ok {
		return false
	}
	return pgerr.Code == errCodeUniqueViolation
}
