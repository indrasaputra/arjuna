package postgres_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	sdkpg "github.com/indrasaputra/arjuna/pkg/sdk/database/postgres"
	sdklog "github.com/indrasaputra/arjuna/pkg/sdk/log"
	mock_uow "github.com/indrasaputra/arjuna/pkg/sdk/test/mock/uow"
	"github.com/indrasaputra/arjuna/service/user/entity"
	"github.com/indrasaputra/arjuna/service/user/internal/app"
	"github.com/indrasaputra/arjuna/service/user/internal/repository/postgres"
)

var (
	testCtx             = context.Background()
	errPostgresInternal = errors.New("error")
	testEnv             = "development"
)

type UserSuite struct {
	user *postgres.User
	db   *mock_uow.MockDB
	tx   *mock_uow.MockTx
}

func TestNewUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("successfully create an instance of User", func(t *testing.T) {
		st := createUserSuite(ctrl)
		assert.NotNil(t, st.user)
	})
}

func TestUser_InsertWithTx(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	app.Logger = sdklog.NewLogger(testEnv)

	query := "INSERT INTO " +
		"users (id, name, created_at, updated_at, created_by, updated_by) " +
		"VALUES (?, ?, ?, ?, ?, ?)"

	t.Run("nil tx is prohibited", func(t *testing.T) {
		st := createUserSuite(ctrl)

		err := st.user.InsertWithTx(testCtx, nil, nil)

		assert.Error(t, err)
	})

	t.Run("nil user is prohibited", func(t *testing.T) {
		st := createUserSuite(ctrl)

		err := st.user.InsertWithTx(testCtx, st.tx, nil)

		assert.Error(t, err)
		assert.Equal(t, entity.ErrEmptyUser(), err)
	})

	t.Run("insert duplicate user", func(t *testing.T) {
		user := createTestUser()
		st := createUserSuite(ctrl)
		st.tx.EXPECT().
			Exec(testCtx, query, user.ID, user.Name, user.CreatedAt, user.UpdatedAt, user.CreatedBy, user.UpdatedBy).
			Return(int64(0), sdkpg.ErrAlreadyExist)

		err := st.user.InsertWithTx(testCtx, st.tx, user)

		assert.Error(t, err)
		assert.Equal(t, entity.ErrAlreadyExists(), err)
	})

	t.Run("insert returns error", func(t *testing.T) {
		user := createTestUser()
		st := createUserSuite(ctrl)
		st.tx.EXPECT().
			Exec(testCtx, query, user.ID, user.Name, user.CreatedAt, user.UpdatedAt, user.CreatedBy, user.UpdatedBy).
			Return(int64(0), entity.ErrInternal(""))

		err := st.user.InsertWithTx(testCtx, st.tx, user)

		assert.Error(t, err)
	})

	t.Run("success insert user", func(t *testing.T) {
		user := createTestUser()
		st := createUserSuite(ctrl)
		st.tx.EXPECT().
			Exec(testCtx, query, user.ID, user.Name, user.CreatedAt, user.UpdatedAt, user.CreatedBy, user.UpdatedBy).
			Return(int64(1), nil)

		err := st.user.InsertWithTx(testCtx, st.tx, user)

		assert.NoError(t, err)
	})
}

func TestUser_GetByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	app.Logger = sdklog.NewLogger(testEnv)

	query := `SELECT id, name, created_at, updated_at, created_by, updated_by FROM users WHERE id = ? LIMIT 1`

	t.Run("select by id returns empty row", func(t *testing.T) {
		user := createTestUser()
		st := createUserSuite(ctrl)
		st.db.EXPECT().
			Query(testCtx, gomock.Any(), query, user.ID).
			Return(entity.ErrNotFound())

		res, err := st.user.GetByID(testCtx, user.ID)

		assert.Error(t, err)
		assert.Nil(t, res)
	})

	t.Run("select by id returns empty row", func(t *testing.T) {
		user := createTestUser()
		st := createUserSuite(ctrl)
		st.db.EXPECT().
			Query(testCtx, gomock.Any(), query, user.ID).
			Return(errPostgresInternal)

		res, err := st.user.GetByID(testCtx, user.ID)

		assert.Error(t, err)
		assert.Nil(t, res)
	})

	t.Run("success select by id", func(t *testing.T) {
		user := createTestUser()
		st := createUserSuite(ctrl)
		st.db.EXPECT().
			Query(testCtx, gomock.Any(), query, user.ID).
			Return(nil)

		res, err := st.user.GetByID(testCtx, user.ID)

		assert.NoError(t, err)
		assert.NotNil(t, res)
	})
}

func TestUser_GetAll(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	app.Logger = sdklog.NewLogger(testEnv)

	query := `SELECT id, name, created_at, updated_at, created_by, updated_by FROM users LIMIT ?`
	limit := uint(10)

	t.Run("get all returns error", func(t *testing.T) {
		st := createUserSuite(ctrl)
		st.db.EXPECT().
			Query(testCtx, gomock.Any(), query, limit).
			Return(errPostgresInternal)

		res, err := st.user.GetAll(testCtx, limit)

		assert.Error(t, err)
		assert.Empty(t, res)
	})

	t.Run("success get all", func(t *testing.T) {
		st := createUserSuite(ctrl)
		st.db.EXPECT().
			Query(testCtx, gomock.Any(), query, limit).
			Return(nil)

		res, err := st.user.GetAll(testCtx, limit)

		assert.NoError(t, err)
		assert.Empty(t, res)
	})
}

func TestUser_HardDeleteWithTx(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	app.Logger = sdklog.NewLogger(testEnv)

	query := "DELETE FROM users WHERE id = ?"

	t.Run("tx is not set", func(t *testing.T) {
		user := createTestUser()
		st := createUserSuite(ctrl)

		err := st.user.HardDeleteWithTx(testCtx, nil, user.ID)

		assert.Error(t, err)
	})

	t.Run("hard delete returns error", func(t *testing.T) {
		user := createTestUser()
		st := createUserSuite(ctrl)
		st.tx.EXPECT().
			Exec(testCtx, query, user.ID).
			Return(int64(0), errPostgresInternal)

		err := st.user.HardDeleteWithTx(testCtx, st.tx, user.ID)

		assert.Error(t, err)
	})

	t.Run("success hard delete", func(t *testing.T) {
		user := createTestUser()
		st := createUserSuite(ctrl)
		st.tx.EXPECT().
			Exec(testCtx, query, user.ID).
			Return(int64(0), nil)

		err := st.user.HardDeleteWithTx(testCtx, st.tx, user.ID)

		assert.NoError(t, err)
	})
}

func TestUser_HardDelete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	app.Logger = sdklog.NewLogger(testEnv)

	query := "DELETE FROM users WHERE id = ?"

	t.Run("hard delete returns error", func(t *testing.T) {
		user := createTestUser()
		st := createUserSuite(ctrl)
		st.db.EXPECT().
			Exec(testCtx, query, user.ID).
			Return(int64(0), errPostgresInternal)

		err := st.user.HardDelete(testCtx, user.ID)

		assert.Error(t, err)
	})

	t.Run("success hard delete", func(t *testing.T) {
		user := createTestUser()
		st := createUserSuite(ctrl)
		st.db.EXPECT().
			Exec(testCtx, query, user.ID).
			Return(int64(0), nil)

		err := st.user.HardDelete(testCtx, user.ID)

		assert.NoError(t, err)
	})
}

func createTestUser() *entity.User {
	return &entity.User{
		ID:    "1",
		Name:  "First User",
		Email: "first@user.com",
	}
}

func createUserSuite(ctrl *gomock.Controller) *UserSuite {
	db := mock_uow.NewMockDB(ctrl)
	tx := mock_uow.NewMockTx(ctrl)
	user := postgres.NewUser(db)
	return &UserSuite{
		user: user,
		db:   db,
		tx:   tx,
	}
}
