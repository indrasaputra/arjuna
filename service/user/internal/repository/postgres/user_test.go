package postgres_test

import (
	"context"
	"errors"
	"log"
	"testing"
	"time"

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
	columns             = []string{"id", "keycloak_id", "name", "email", "created_at", "updated_at", "created_by", "updated_by"}
	testUserID          = "1"
	testUserKeycloakID  = "1"
	testUserName        = "Zlatan Ibrahimovic"
	testUserEmail       = "zlatan@ibrahimovic.com"
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
	t.Run("nil user is prohibited", func(t *testing.T) {
		exec := createUserExecutor()

		err := exec.user.Insert(testCtx, nil)

		assert.Error(t, err)
		assert.Equal(t, entity.ErrEmptyUser(), err)
	})

	query := `INSERT INTO users \(id, keycloak_id, name, email, created_at, updated_at, created_by, updated_by\) VALUES \(\$1, \$2, \$3, \$4, \$5, \$6, \$7, \$8\)`

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

func TestToggle_GetAll(t *testing.T) {
	query := `SELECT id, keycloak_id, name, email, created_at, updated_at, created_by, updated_by FROM users LIMIT \$1`

	t.Run("select all query returns error", func(t *testing.T) {
		exec := createUserExecutor()
		exec.pgx.
			ExpectQuery(query).
			WillReturnError(errPostgresInternal)

		res, err := exec.user.GetAll(testCtx, 10)

		assert.Error(t, err)
		assert.Empty(t, res)
	})

	t.Run("select all rows scan returns error", func(t *testing.T) {
		exec := createUserExecutor()
		exec.pgx.
			ExpectQuery(query).
			WillReturnRows(pgxmock.
				NewRows(columns).
				AddRow(testUserID, testUserKeycloakID, testUserName, testUserEmail, time.Now(), time.Now(), testUserID, testUserID).
				AddRow(testUserID, testUserKeycloakID, testUserName, testUserEmail, "time.Now()", "time.Now()", testUserID, testUserID),
			)

		res, err := exec.user.GetAll(testCtx, 10)

		assert.Nil(t, err)
		assert.Equal(t, 1, len(res))
	})

	t.Run("select all rows error occurs after scanning", func(t *testing.T) {
		exec := createUserExecutor()
		exec.pgx.
			ExpectQuery(query).
			WillReturnRows(pgxmock.
				NewRows(columns).
				AddRow(testUserID, testUserKeycloakID, testUserName, testUserEmail, time.Now(), time.Now(), testUserID, testUserID).
				AddRow(testUserID, testUserKeycloakID, testUserName, testUserEmail, "time.Now()", "time.Now()", testUserID, testUserID).
				RowError(2, errPostgresInternal),
			)

		res, err := exec.user.GetAll(testCtx, 10)

		assert.Error(t, err)
		assert.Empty(t, res)
	})

	t.Run("successfully retrieve all rows", func(t *testing.T) {
		exec := createUserExecutor()
		exec.pgx.
			ExpectQuery(query).
			WillReturnRows(pgxmock.
				NewRows(columns).
				AddRow(testUserID, testUserKeycloakID, testUserName, testUserEmail, time.Now(), time.Now(), testUserID, testUserID).
				AddRow(testUserID, testUserKeycloakID, testUserName, testUserEmail, time.Now(), time.Now(), testUserID, testUserID),
			)

		res, err := exec.user.GetAll(testCtx, 10)

		assert.Nil(t, err)
		assert.Equal(t, 2, len(res))
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
