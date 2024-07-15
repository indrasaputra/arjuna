package postgres_test

import (
	"context"
	"database/sql"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	sdkpg "github.com/indrasaputra/arjuna/pkg/sdk/database/postgres"
	sdklog "github.com/indrasaputra/arjuna/pkg/sdk/log"
	mock_uow "github.com/indrasaputra/arjuna/pkg/sdk/test/mock/uow"
	"github.com/indrasaputra/arjuna/service/auth/entity"
	"github.com/indrasaputra/arjuna/service/auth/internal/app"
	"github.com/indrasaputra/arjuna/service/auth/internal/repository/postgres"
)

var (
	testCtx = context.Background()
	testEnv = "development"
)

type AccountSuite struct {
	account *postgres.Account
	db      *mock_uow.MockDB
}

func TestNewAccount(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("successfully create an instance of Account", func(t *testing.T) {
		st := createAccountSuite(ctrl)
		assert.NotNil(t, st.account)
	})
}

func TestAccount_Insert(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	app.Logger = sdklog.NewLogger(testEnv)
	query := "INSERT INTO " +
		"accounts (id, user_id, email, password, created_at, updated_at, created_by, updated_by) " +
		"VALUES (?, ?, ?, ?, ?, ?, ?, ?)"

	t.Run("nil account is prohibited", func(t *testing.T) {
		st := createAccountSuite(ctrl)

		err := st.account.Insert(testCtx, nil)

		assert.Error(t, err)
		assert.Equal(t, entity.ErrEmptyAccount(), err)
	})

	t.Run("insert duplicate account", func(t *testing.T) {
		account := createTestAccount()
		st := createAccountSuite(ctrl)
		st.db.EXPECT().
			Exec(testCtx, query, account.ID, account.UserID, account.Email, account.Password, account.CreatedAt, account.UpdatedAt, account.CreatedBy, account.UpdatedBy).
			Return(int64(0), sdkpg.ErrAlreadyExist)

		err := st.account.Insert(testCtx, account)

		assert.Error(t, err)
		assert.Equal(t, entity.ErrAlreadyExists(), err)
	})

	t.Run("insert returns error", func(t *testing.T) {
		account := createTestAccount()
		st := createAccountSuite(ctrl)
		st.db.EXPECT().
			Exec(testCtx, query, account.ID, account.UserID, account.Email, account.Password, account.CreatedAt, account.UpdatedAt, account.CreatedBy, account.UpdatedBy).
			Return(int64(0), entity.ErrInternal(""))

		err := st.account.Insert(testCtx, account)

		assert.Error(t, err)
	})

	t.Run("success insert account", func(t *testing.T) {
		account := createTestAccount()
		st := createAccountSuite(ctrl)
		st.db.EXPECT().
			Exec(testCtx, query, account.ID, account.UserID, account.Email, account.Password, account.CreatedAt, account.UpdatedAt, account.CreatedBy, account.UpdatedBy).
			Return(int64(1), nil)

		err := st.account.Insert(testCtx, account)

		assert.NoError(t, err)
	})
}

func TestAccount_GetByEmail(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	app.Logger = sdklog.NewLogger(testEnv)
	query := `SELECT id, user_id, email, password FROM accounts WHERE email = ? LIMIT 1`

	t.Run("get by email returns empty row", func(t *testing.T) {
		acc := createTestAccount()
		st := createAccountSuite(ctrl)
		st.db.EXPECT().
			Query(testCtx, gomock.Any(), query, acc.Email).
			Return(sql.ErrNoRows)

		res, err := st.account.GetByEmail(testCtx, acc.Email)

		assert.Error(t, err)
		assert.Equal(t, entity.ErrNotFound(), err)
		assert.Nil(t, res)
	})

	t.Run("get by email returns error", func(t *testing.T) {
		acc := createTestAccount()
		st := createAccountSuite(ctrl)
		st.db.EXPECT().
			Query(testCtx, gomock.Any(), query, acc.Email).
			Return(entity.ErrInternal("error"))

		res, err := st.account.GetByEmail(testCtx, acc.Email)

		assert.Error(t, err)
		assert.Nil(t, res)
	})

	t.Run("success select by email", func(t *testing.T) {
		acc := createTestAccount()
		st := createAccountSuite(ctrl)
		st.db.EXPECT().
			Query(testCtx, gomock.Any(), query, acc.Email).
			Return(nil)

		res, err := st.account.GetByEmail(testCtx, acc.Email)

		assert.NoError(t, err)
		assert.NotNil(t, res)
	})
}

func createTestAccount() *entity.Account {
	return &entity.Account{
		ID:       "1",
		UserID:   "1",
		Email:    "first@account.com",
		Password: "password",
	}
}

func createAccountSuite(ctrl *gomock.Controller) *AccountSuite {
	db := mock_uow.NewMockDB(ctrl)
	ac := postgres.NewAccount(db)
	return &AccountSuite{
		account: ac,
		db:      db,
	}
}
