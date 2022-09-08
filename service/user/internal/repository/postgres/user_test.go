package postgres_test

import (
	"context"
	"errors"
	"log"
	"testing"

	"github.com/jackc/pgconn"
	"github.com/pashagolub/pgxmock"
	"github.com/stretchr/testify/assert"

	"github.com/indrasaputra/arjuna/service/user/entity"
	"github.com/indrasaputra/arjuna/service/user/internal/repository/postgres"
)

var (
	testCtx             = context.Background()
	testUser            = &entity.User{Name: "Zlatan Ibrahimovic", Email: "zlatan@ibrahimovic.com"}
	errPostgresInternal = errors.New("error")
)

type UserExecutor struct {
	user *postgres.User
	pgx  pgxmock.PgxPoolIface
}

func TestNewUser(t *testing.T) {
	t.Run("successfully create an instance of User", func(t *testing.T) {
		exec := createUserExecutor()
		assert.NotNil(t, exec.user)
	})
}

func TestUser_Insert(t *testing.T) {
	t.Run("nil toggle is prohibited", func(t *testing.T) {
		exec := createUserExecutor()

		err := exec.user.Insert(testCtx, nil)

		assert.Error(t, err)
		assert.Equal(t, entity.ErrEmptyUser(), err)
	})

	query := `INSERT INTO users \(id, name, email, username, created_at, updated_at, created_by, updated_by\) VALUES \(\$1, \$2, \$3, \$4, \$5, \$6, \$7, \$8\)`

	t.Run("insert duplicate user", func(t *testing.T) {
		exec := createUserExecutor()
		exec.pgx.
			ExpectExec(query).
			WillReturnError(&pgconn.PgError{Code: "23505"})

		err := exec.user.Insert(testCtx, testUser)

		assert.Error(t, err)
		assert.Equal(t, entity.ErrAlreadyExists(), err)
	})

	t.Run("postgres returns error", func(t *testing.T) {
		exec := createUserExecutor()
		exec.pgx.
			ExpectExec(query).
			WillReturnError(errPostgresInternal)

		err := exec.user.Insert(testCtx, testUser)

		assert.Error(t, err)
	})

	t.Run("success save to postgres", func(t *testing.T) {
		exec := createUserExecutor()
		exec.pgx.
			ExpectExec(query).
			WillReturnResult(pgxmock.NewResult("INSERT", 1))

		err := exec.user.Insert(testCtx, testUser)

		assert.NoError(t, err)
	})
}

func createUserExecutor() *UserExecutor {
	mock, err := pgxmock.NewPool(pgxmock.MonitorPingsOption(true))
	if err != nil {
		log.Panicf("error opening a stub database connection: %v\n", err)
	}

	user := postgres.NewUser(mock)
	return &UserExecutor{
		user: user,
		pgx:  mock,
	}
}
