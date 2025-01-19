package postgres_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/pashagolub/pgxmock/v4"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	sdkpostgres "github.com/indrasaputra/arjuna/pkg/sdk/database/postgres"
	sdklog "github.com/indrasaputra/arjuna/pkg/sdk/log"
	mock_uow "github.com/indrasaputra/arjuna/pkg/sdk/test/mock/uow"
	"github.com/indrasaputra/arjuna/service/user/entity"
	"github.com/indrasaputra/arjuna/service/user/internal/app"
	"github.com/indrasaputra/arjuna/service/user/internal/repository/db"
	"github.com/indrasaputra/arjuna/service/user/internal/repository/postgres"
)

var (
	testCtx = context.Background()
	testEnv = "development"
)

type UserSuite struct {
	user   *postgres.User
	db     pgxmock.PgxPoolIface
	getter *mock_uow.MockTxGetter
}

func TestNewUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("successfully create an instance of User", func(t *testing.T) {
		st := createUserSuite(t, ctrl)
		assert.NotNil(t, st.user)
	})
}

func TestUser_Insert(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	app.Logger = sdklog.NewLogger(testEnv)

	query := `INSERT INTO users \(id, name, created_at, updated_at, created_by, updated_by\) VALUES \(\$1, \$2, \$3, \$4, \$5, \$6\)`

	t.Run("nil user is prohibited", func(t *testing.T) {
		st := createUserSuite(t, ctrl)

		err := st.user.Insert(testCtx, nil)

		assert.Error(t, err)
		assert.Equal(t, entity.ErrEmptyUser(), err)
	})

	t.Run("insert duplicate user", func(t *testing.T) {
		user := createTestUser()
		st := createUserSuite(t, ctrl)
		st.getter.EXPECT().DefaultTrOrDB(testCtx, st.db).Return(st.db)
		st.db.ExpectExec(query).
			WithArgs(user.ID, user.Name, user.CreatedAt, user.UpdatedAt, user.CreatedBy, user.UpdatedBy).
			WillReturnError(&pgconn.PgError{Code: "23505"})

		err := st.user.Insert(testCtx, user)

		assert.Error(t, err)
		assert.Equal(t, entity.ErrAlreadyExists(), err)
	})

	t.Run("insert returns error", func(t *testing.T) {
		user := createTestUser()
		st := createUserSuite(t, ctrl)
		st.getter.EXPECT().DefaultTrOrDB(testCtx, st.db).Return(st.db)
		st.db.ExpectExec(query).
			WithArgs(user.ID, user.Name, user.CreatedAt, user.UpdatedAt, user.CreatedBy, user.UpdatedBy).
			WillReturnError(assert.AnError)

		err := st.user.Insert(testCtx, user)

		assert.Error(t, err)
	})

	t.Run("success insert user", func(t *testing.T) {
		user := createTestUser()
		st := createUserSuite(t, ctrl)
		st.getter.EXPECT().DefaultTrOrDB(testCtx, st.db).Return(st.db)
		st.db.ExpectExec(query).
			WithArgs(user.ID, user.Name, user.CreatedAt, user.UpdatedAt, user.CreatedBy, user.UpdatedBy).
			WillReturnResult(pgxmock.NewResult("INSERT", 1))

		err := st.user.Insert(testCtx, user)

		assert.NoError(t, err)
	})
}

func TestUser_GetByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	app.Logger = sdklog.NewLogger(testEnv)

	query := `SELECT id, name, created_at, updated_at, deleted_at, created_by, updated_by, deleted_by FROM users WHERE id = \$1 LIMIT 1`

	t.Run("select by id returns empty row", func(t *testing.T) {
		user := createTestUser()
		st := createUserSuite(t, ctrl)
		st.getter.EXPECT().DefaultTrOrDB(testCtx, st.db).Return(st.db)
		st.db.ExpectQuery(query).WithArgs(user.ID).WillReturnError(pgx.ErrNoRows)

		res, err := st.user.GetByID(testCtx, user.ID)

		assert.Error(t, err)
		assert.Equal(t, entity.ErrNotFound(), err)
		assert.Nil(t, res)
	})

	t.Run("select by id returns error", func(t *testing.T) {
		user := createTestUser()
		st := createUserSuite(t, ctrl)
		st.getter.EXPECT().DefaultTrOrDB(testCtx, st.db).Return(st.db)
		st.db.ExpectQuery(query).WithArgs(user.ID).WillReturnError(assert.AnError)

		res, err := st.user.GetByID(testCtx, user.ID)

		assert.Error(t, err)
		assert.Nil(t, res)
	})

	t.Run("success select by id", func(t *testing.T) {
		user := createTestUser()
		st := createUserSuite(t, ctrl)
		st.getter.EXPECT().DefaultTrOrDB(testCtx, st.db).Return(st.db)
		st.db.ExpectQuery(query).WithArgs(user.ID).WillReturnRows(pgxmock.
			NewRows([]string{"id", "name", "created_at", "updated_at", "deleted_at", "created_by", "updated_by", "deleted_by"}).
			AddRow(user.ID, user.Name, user.CreatedAt, user.UpdatedAt, user.DeletedAt, user.CreatedBy, user.UpdatedBy, user.DeletedBy))

		res, err := st.user.GetByID(testCtx, user.ID)

		assert.NoError(t, err)
		assert.NotNil(t, res)
	})
}

func TestUser_GetAll(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	app.Logger = sdklog.NewLogger(testEnv)

	query := `SELECT id, name, created_at, updated_at, deleted_at, created_by, updated_by, deleted_by FROM users LIMIT \$1`
	limit := uint(10)

	t.Run("get all returns error", func(t *testing.T) {
		st := createUserSuite(t, ctrl)
		st.getter.EXPECT().DefaultTrOrDB(testCtx, st.db).Return(st.db)
		st.db.ExpectQuery(query).WithArgs(limit).WillReturnError(assert.AnError)

		res, err := st.user.GetAll(testCtx, limit)

		assert.Error(t, err)
		assert.Empty(t, res)
	})

	t.Run("success get all", func(t *testing.T) {
		user := createTestUser()
		st := createUserSuite(t, ctrl)
		st.getter.EXPECT().DefaultTrOrDB(testCtx, st.db).Return(st.db)
		st.db.ExpectQuery(query).WithArgs(int32(limit)).WillReturnRows(pgxmock.
			NewRows([]string{"id", "name", "created_at", "updated_at", "deleted_at", "created_by", "updated_by", "deleted_by"}).
			AddRow(user.ID, user.Name, user.CreatedAt, user.UpdatedAt, user.DeletedAt, user.CreatedBy, user.UpdatedBy, user.DeletedBy))

		res, err := st.user.GetAll(testCtx, limit)

		assert.NoError(t, err)
		assert.Len(t, res, 1)
	})
}

func TestUser_HardDelete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	app.Logger = sdklog.NewLogger(testEnv)

	query := `DELETE FROM users WHERE id = \$1`

	t.Run("hard delete returns error", func(t *testing.T) {
		user := createTestUser()
		st := createUserSuite(t, ctrl)
		st.getter.EXPECT().DefaultTrOrDB(testCtx, st.db).Return(st.db)
		st.db.ExpectExec(query).WithArgs(user.ID).WillReturnError(assert.AnError)

		err := st.user.HardDelete(testCtx, user.ID)

		assert.Error(t, err)
	})

	t.Run("success hard delete", func(t *testing.T) {
		user := createTestUser()
		st := createUserSuite(t, ctrl)
		st.getter.EXPECT().DefaultTrOrDB(testCtx, st.db).Return(st.db)
		st.db.ExpectExec(query).WithArgs(user.ID).WillReturnResult(pgxmock.NewResult("DELETE", 1))

		err := st.user.HardDelete(testCtx, user.ID)

		assert.NoError(t, err)
	})
}

func createTestUser() *entity.User {
	return &entity.User{
		ID:    uuid.Must(uuid.NewV7()),
		Name:  "First User",
		Email: "first@user.com",
	}
}

func createUserSuite(t *testing.T, ctrl *gomock.Controller) *UserSuite {
	pool, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("error opening a stub database connection: %v\n", err)
	}
	g := mock_uow.NewMockTxGetter(ctrl)
	tx := sdkpostgres.NewTxDB(pool, g)
	q := db.New(tx)
	user := postgres.NewUser(q)
	return &UserSuite{
		user:   user,
		db:     pool,
		getter: g,
	}
}
