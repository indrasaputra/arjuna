package postgres_test

import (
	"context"
	"database/sql"
	"testing"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	sdkpg "github.com/indrasaputra/arjuna/pkg/sdk/database/postgres"
	sdklog "github.com/indrasaputra/arjuna/pkg/sdk/log"
	mock_uow "github.com/indrasaputra/arjuna/pkg/sdk/test/mock/uow"
	"github.com/indrasaputra/arjuna/service/wallet/entity"
	"github.com/indrasaputra/arjuna/service/wallet/internal/app"
	"github.com/indrasaputra/arjuna/service/wallet/internal/repository/postgres"
)

var (
	testCtx = context.Background()
	testEnv = "development"
)

type WalletSuite struct {
	wallet *postgres.Wallet
	db     *mock_uow.MockDB
	tx     *mock_uow.MockTx
}

func TestNewWallet(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("successfully create an instance of Wallet", func(t *testing.T) {
		st := createWalletSuite(ctrl)
		assert.NotNil(t, st.wallet)
	})
}

func TestWallet_Insert(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	app.Logger = sdklog.NewLogger(testEnv)
	query := "INSERT INTO " +
		"wallets (id, user_id, balance, created_at, updated_at, created_by, updated_by) " +
		"VALUES (?, ?, ?, ?, ?, ?, ?)"

	t.Run("nil wallets is prohibited", func(t *testing.T) {
		st := createWalletSuite(ctrl)

		err := st.wallet.Insert(testCtx, nil)

		assert.Error(t, err)
		assert.Equal(t, entity.ErrEmptyWallet(), err)
	})

	t.Run("insert duplicate wallets", func(t *testing.T) {
		wallet := createTestWallet()
		st := createWalletSuite(ctrl)
		st.db.EXPECT().
			Exec(testCtx, query, wallet.ID, wallet.UserID, wallet.Balance, wallet.CreatedAt, wallet.UpdatedAt, wallet.CreatedBy, wallet.UpdatedBy).
			Return(int64(0), sdkpg.ErrAlreadyExist)

		err := st.wallet.Insert(testCtx, wallet)

		assert.Error(t, err)
		assert.Equal(t, entity.ErrAlreadyExists(), err)
	})

	t.Run("insert returns error", func(t *testing.T) {
		wallet := createTestWallet()
		st := createWalletSuite(ctrl)
		st.db.EXPECT().
			Exec(testCtx, query, wallet.ID, wallet.UserID, wallet.Balance, wallet.CreatedAt, wallet.UpdatedAt, wallet.CreatedBy, wallet.UpdatedBy).
			Return(int64(0), entity.ErrInternal(""))

		err := st.wallet.Insert(testCtx, wallet)

		assert.Error(t, err)
	})

	t.Run("success insert wallets", func(t *testing.T) {
		wallet := createTestWallet()
		st := createWalletSuite(ctrl)
		st.db.EXPECT().
			Exec(testCtx, query, wallet.ID, wallet.UserID, wallet.Balance, wallet.CreatedAt, wallet.UpdatedAt, wallet.CreatedBy, wallet.UpdatedBy).
			Return(int64(1), nil)

		err := st.wallet.Insert(testCtx, wallet)

		assert.NoError(t, err)
	})
}

func TestWallet_AddWalletBalance(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	app.Logger = sdklog.NewLogger(testEnv)

	query := `UPDATE wallets SET balance = balance + ? WHERE id = ?`

	t.Run("add account balance returns internal error", func(t *testing.T) {
		st := createWalletSuite(ctrl)
		id := uuid.Must(uuid.NewV7())
		amount, _ := decimal.NewFromString("4.56")
		st.db.EXPECT().
			Exec(testCtx, query, amount, id).
			Return(int64(1), entity.ErrInternal(""))

		err := st.wallet.AddWalletBalance(testCtx, id, amount)

		assert.Error(t, err)
	})

	t.Run("add account balance returns success", func(t *testing.T) {
		st := createWalletSuite(ctrl)
		id := uuid.Must(uuid.NewV7())
		amount, _ := decimal.NewFromString("4.56")
		st.db.EXPECT().
			Exec(testCtx, query, amount, id).
			Return(int64(1), nil)

		err := st.wallet.AddWalletBalance(testCtx, id, amount)

		assert.NoError(t, err)
	})
}

func TestWallet_AddWalletBalanceWithTx(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	app.Logger = sdklog.NewLogger(testEnv)

	query := `UPDATE wallets SET balance = balance + ? WHERE id = ?`

	t.Run("tx is nil", func(t *testing.T) {
		st := createWalletSuite(ctrl)
		id := uuid.Must(uuid.NewV7())
		amount, _ := decimal.NewFromString("4.56")

		err := st.wallet.AddWalletBalanceWithTx(testCtx, nil, id, amount)

		assert.Error(t, err)
	})

	t.Run("add account balance with tx returns internal error", func(t *testing.T) {
		st := createWalletSuite(ctrl)
		id := uuid.Must(uuid.NewV7())
		amount, _ := decimal.NewFromString("4.56")
		st.tx.EXPECT().
			Exec(testCtx, query, amount, id).
			Return(int64(1), entity.ErrInternal(""))

		err := st.wallet.AddWalletBalanceWithTx(testCtx, st.tx, id, amount)

		assert.Error(t, err)
	})

	t.Run("add account balance with tx returns success", func(t *testing.T) {
		st := createWalletSuite(ctrl)
		id := uuid.Must(uuid.NewV7())
		amount, _ := decimal.NewFromString("4.56")
		st.tx.EXPECT().
			Exec(testCtx, query, amount, id).
			Return(int64(1), nil)

		err := st.wallet.AddWalletBalanceWithTx(testCtx, st.tx, id, amount)

		assert.NoError(t, err)
	})
}

func TestWallet_GetUserWalletWithTx(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	app.Logger = sdklog.NewLogger(testEnv)

	query := `SELECT id, user_id, balance FROM wallets WHERE id = ? AND user_id = ? LIMIT 1 FOR NO KEY UPDATE`

	t.Run("tx is nil", func(t *testing.T) {
		st := createWalletSuite(ctrl)
		id := uuid.Must(uuid.NewV7())
		userID := uuid.Must(uuid.NewV7())

		res, err := st.wallet.GetUserWalletWithTx(testCtx, nil, id, userID)

		assert.Error(t, err)
		assert.Nil(t, res)
	})

	t.Run("add account balance returns no rows", func(t *testing.T) {
		st := createWalletSuite(ctrl)
		id := uuid.Must(uuid.NewV7())
		userID := uuid.Must(uuid.NewV7())
		st.tx.EXPECT().
			Query(testCtx, gomock.Any(), query, id, userID).
			Return(sql.ErrNoRows)

		res, err := st.wallet.GetUserWalletWithTx(testCtx, st.tx, id, userID)

		assert.Error(t, err)
		assert.Equal(t, entity.ErrEmptyWallet(), err)
		assert.Nil(t, res)
	})

	t.Run("add account balance returns internal error", func(t *testing.T) {
		st := createWalletSuite(ctrl)
		id := uuid.Must(uuid.NewV7())
		userID := uuid.Must(uuid.NewV7())
		st.tx.EXPECT().
			Query(testCtx, gomock.Any(), query, id, userID).
			Return(entity.ErrInternal(""))

		res, err := st.wallet.GetUserWalletWithTx(testCtx, st.tx, id, userID)

		assert.Error(t, err)
		assert.Nil(t, res)
	})

	t.Run("success add balance", func(t *testing.T) {
		st := createWalletSuite(ctrl)
		id := uuid.Must(uuid.NewV7())
		userID := uuid.Must(uuid.NewV7())
		st.tx.EXPECT().
			Query(testCtx, gomock.Any(), query, id, userID).
			Return(nil)

		res, err := st.wallet.GetUserWalletWithTx(testCtx, st.tx, id, userID)

		assert.NoError(t, err)
		assert.NotNil(t, res)
	})
}

func createTestWallet() *entity.Wallet {
	b, _ := decimal.NewFromString("10.23")
	return &entity.Wallet{
		ID:      uuid.Must(uuid.NewV7()),
		UserID:  uuid.Must(uuid.NewV7()),
		Balance: b,
	}
}

func createWalletSuite(ctrl *gomock.Controller) *WalletSuite {
	db := mock_uow.NewMockDB(ctrl)
	t := mock_uow.NewMockTx(ctrl)
	w := postgres.NewWallet(db)
	return &WalletSuite{
		wallet: w,
		db:     db,
		tx:     t,
	}
}
