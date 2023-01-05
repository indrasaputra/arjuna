package postgres_test

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	pgsdk "github.com/indrasaputra/arjuna/pkg/sdk/database/postgres"
	mock_uow "github.com/indrasaputra/arjuna/pkg/sdk/test/mock/uow"
	"github.com/indrasaputra/arjuna/service/user/entity"
	"github.com/indrasaputra/arjuna/service/user/internal/repository/postgres"
)

var (
	testCtx             = context.Background()
	errPostgresInternal = errors.New("error")
)

type UserExecutor struct {
	user *postgres.User
	db   *mock_uow.MockDB
	tx   *mock_uow.MockTx
}

func TestNewUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("successfully create an instance of User", func(t *testing.T) {
		exec := createUserExecutor(ctrl)
		assert.NotNil(t, exec.user)
	})
}

func TestUser_Insert(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	query := "INSERT INTO " +
		"users (id, keycloak_id, name, email, created_at, updated_at, created_by, updated_by) " +
		"VALUES (?, ?, ?, ?, ?, ?, ?, ?)"

	t.Run("nil user is prohibited", func(t *testing.T) {
		exec := createUserExecutor(ctrl)

		err := exec.user.Insert(testCtx, nil)

		assert.Error(t, err)
		assert.Equal(t, entity.ErrEmptyUser(), err)
	})

	t.Run("insert duplicate user", func(t *testing.T) {
		user := createTestUser()
		exec := createUserExecutor(ctrl)
		exec.db.EXPECT().
			Exec(testCtx, query, user.ID, user.KeycloakID, user.Name, user.Email, user.CreatedAt, user.UpdatedAt, user.CreatedBy, user.UpdatedBy).
			Return(int64(0), pgsdk.ErrAlreadyExist)

		err := exec.user.Insert(testCtx, user)

		assert.Error(t, err)
		assert.Equal(t, entity.ErrAlreadyExists(), err)
	})

	t.Run("insert returns error", func(t *testing.T) {
		user := createTestUser()
		exec := createUserExecutor(ctrl)
		exec.db.EXPECT().
			Exec(testCtx, query, user.ID, user.KeycloakID, user.Name, user.Email, user.CreatedAt, user.UpdatedAt, user.CreatedBy, user.UpdatedBy).
			Return(int64(0), entity.ErrInternal(""))

		err := exec.user.Insert(testCtx, user)

		assert.Error(t, err)
	})

	t.Run("success insert user", func(t *testing.T) {
		user := createTestUser()
		exec := createUserExecutor(ctrl)
		exec.db.EXPECT().
			Exec(testCtx, query, user.ID, user.KeycloakID, user.Name, user.Email, user.CreatedAt, user.UpdatedAt, user.CreatedBy, user.UpdatedBy).
			Return(int64(1), nil)

		err := exec.user.Insert(testCtx, user)

		assert.NoError(t, err)
	})
}

func TestUser_GetByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	query := `SELECT id, keycloak_id, name, email, created_at, updated_at, created_by, updated_by FROM users WHERE id = ? LIMIT 1`

	t.Run("select by id returns empty row", func(t *testing.T) {
		user := createTestUser()
		exec := createUserExecutor(ctrl)
		exec.db.EXPECT().
			Query(testCtx, gomock.Any(), query, user.ID).
			Return(entity.ErrNotFound())

		res, err := exec.user.GetByID(testCtx, user.ID)

		assert.Error(t, err)
		assert.Nil(t, res)
	})

	t.Run("select by id returns empty row", func(t *testing.T) {
		user := createTestUser()
		exec := createUserExecutor(ctrl)
		exec.db.EXPECT().
			Query(testCtx, gomock.Any(), query, user.ID).
			Return(errPostgresInternal)

		res, err := exec.user.GetByID(testCtx, user.ID)

		assert.Error(t, err)
		assert.Nil(t, res)
	})

	t.Run("success select by id", func(t *testing.T) {
		user := createTestUser()
		exec := createUserExecutor(ctrl)
		exec.db.EXPECT().
			Query(testCtx, gomock.Any(), query, user.ID).
			Return(nil)

		res, err := exec.user.GetByID(testCtx, user.ID)

		assert.NoError(t, err)
		assert.NotNil(t, res)
	})
}

func TestUser_GetAll(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	query := `SELECT id, keycloak_id, name, email, created_at, updated_at, created_by, updated_by FROM users LIMIT ?`
	limit := uint(10)

	t.Run("get all returns error", func(t *testing.T) {
		exec := createUserExecutor(ctrl)
		exec.db.EXPECT().
			Query(testCtx, gomock.Any(), query, limit).
			Return(errPostgresInternal)

		res, err := exec.user.GetAll(testCtx, limit)

		assert.Error(t, err)
		assert.Empty(t, res)
	})

	t.Run("success get all", func(t *testing.T) {
		exec := createUserExecutor(ctrl)
		exec.db.EXPECT().
			Query(testCtx, gomock.Any(), query, limit).
			Return(nil)

		res, err := exec.user.GetAll(testCtx, limit)

		assert.NoError(t, err)
		assert.Empty(t, res)
	})
}

func TestUser_HardDelete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	query := "DELETE FROM users WHERE id = ?"

	t.Run("tx is not set", func(t *testing.T) {
		user := createTestUser()
		exec := createUserExecutor(ctrl)
		// exec.db.EXPECT().
		// 	Exec(testCtx, query, user.ID).
		// 	Return(0, )

		err := exec.user.HardDelete(testCtx, nil, user.ID)

		assert.Error(t, err)
	})

	t.Run("hard delete returns error", func(t *testing.T) {
		user := createTestUser()
		exec := createUserExecutor(ctrl)
		exec.tx.EXPECT().
			Exec(testCtx, query, user.ID).
			Return(int64(0), errPostgresInternal)

		err := exec.user.HardDelete(testCtx, exec.tx, user.ID)

		assert.Error(t, err)
	})

	t.Run("success hard delete", func(t *testing.T) {
		user := createTestUser()
		exec := createUserExecutor(ctrl)
		exec.tx.EXPECT().
			Exec(testCtx, query, user.ID).
			Return(int64(0), nil)

		err := exec.user.HardDelete(testCtx, exec.tx, user.ID)

		assert.NoError(t, err)
	})
}

func createTestUser() *entity.User {
	return &entity.User{
		ID:         "1",
		KeycloakID: "1",
		Name:       "First User",
		Email:      "first@user.com",
	}
}

func createUserExecutor(ctrl *gomock.Controller) *UserExecutor {
	db := mock_uow.NewMockDB(ctrl)
	tx := mock_uow.NewMockTx(ctrl)
	user := postgres.NewUser(db)
	return &UserExecutor{
		user: user,
		db:   db,
		tx:   tx,
	}
}
