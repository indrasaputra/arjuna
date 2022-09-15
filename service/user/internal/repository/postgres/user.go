package postgres

import (
	"context"
	"log"

	"github.com/jackc/pgconn"
	pgx "github.com/jackc/pgx/v4"

	"github.com/indrasaputra/arjuna/service/user/entity"
)

const (
	// errCodeUniqueViolation is derived from https://www.postgresql.org/docs/11/errcodes-appendix.html
	errCodeUniqueViolation = "23505"
)

// User is responsible to connect user entity with users table in PostgreSQL.
type User struct {
	pool PgxPool
}

// NewUser creates an instance of User.
func NewUser(pool PgxPool) *User {
	return &User{pool: pool}
}

// Insert inserts the user into the users table.
func (u *User) Insert(ctx context.Context, user *entity.User) error {
	if user == nil {
		return entity.ErrEmptyUser()
	}

	query := "INSERT INTO " +
		"users (id, keycloak_id, name, email, created_at, updated_at, created_by, updated_by) " +
		"VALUES ($1, $2, $3, $4, $5, $6, $7, $8)"

	_, err := u.pool.Exec(ctx, query,
		user.ID,
		user.KeycloakID,
		user.Name,
		user.Email,
		user.CreatedAt,
		user.UpdatedAt,
		user.CreatedBy,
		user.UpdatedBy,
	)

	if err != nil && isUniqueViolationErr(err) {
		return entity.ErrAlreadyExists()
	}
	if err != nil {
		return entity.ErrInternal(err.Error())
	}
	return nil
}

// GetByID gets a user from database.
// It returns entity.ErrNotFound if user can't be found.
func (u *User) GetByID(ctx context.Context, id string) (*entity.User, error) {
	query := "SELECT id, keycloak_id, name, email, created_at, updated_at, created_by, updated_by FROM users WHERE id = $1 LIMIT 1"
	row := u.pool.QueryRow(ctx, query, id)

	res := entity.User{}
	err := row.Scan(&res.ID, &res.KeycloakID, &res.Name, &res.Email, &res.CreatedAt, &res.UpdatedAt, &res.CreatedBy, &res.UpdatedBy)
	if err == pgx.ErrNoRows {
		return nil, entity.ErrNotFound()
	}
	if err != nil {
		return nil, entity.ErrInternal(err.Error())
	}
	return &res, nil
}

// GetAll gets all users in users table.
func (u *User) GetAll(ctx context.Context, limit uint) ([]*entity.User, error) {
	query := "SELECT id, keycloak_id, name, email, created_at, updated_at, created_by, updated_by FROM users LIMIT $1"
	rows, err := u.pool.Query(ctx, query, limit)
	if err != nil {
		return []*entity.User{}, entity.ErrInternal(err.Error())
	}
	defer rows.Close()

	users := []*entity.User{}
	for rows.Next() {
		var tmp entity.User
		if err := rows.Scan(&tmp.ID, &tmp.KeycloakID, &tmp.Name, &tmp.Email, &tmp.CreatedAt, &tmp.UpdatedAt, &tmp.CreatedBy, &tmp.UpdatedBy); err != nil {
			log.Printf("[User-GetAll] postgres scan rows error: %s", err.Error())
			continue
		}
		users = append(users, &tmp)
	}
	if rows.Err() != nil {
		return []*entity.User{}, entity.ErrInternal(rows.Err().Error())
	}
	return users, nil
}

// HardDelete deletes a user from database.
// If the user doesn't exist, it doesn't returns error.
func (u *User) HardDelete(ctx context.Context, id string) error {
	query := "DELETE FROM users WHERE id = $1"
	_, err := u.pool.Exec(ctx, query, id)
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
