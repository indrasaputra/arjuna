package postgres_test

import (
	"context"
	"testing"

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

	query := `UPDATE wallets SET balance = balance + ? WHERE id = ? AND user_id = ?`

	t.Run("add account balance returns internal error", func(t *testing.T) {
		st := createWalletSuite(ctrl)
		topup := createTestTopupWallet()
		st.db.EXPECT().
			Exec(testCtx, query, topup.Amount, topup.WalletID, topup.UserID).
			Return(int64(1), entity.ErrInternal(""))

		err := st.wallet.AddWalletBalance(testCtx, topup)

		assert.Error(t, err)
	})

	t.Run("add account balance returns success", func(t *testing.T) {
		st := createWalletSuite(ctrl)
		topup := createTestTopupWallet()
		st.db.EXPECT().
			Exec(testCtx, query, topup.Amount, topup.WalletID, topup.UserID).
			Return(int64(1), nil)

		err := st.wallet.AddWalletBalance(testCtx, topup)

		assert.NoError(t, err)
	})
}

func createTestWallet() *entity.Wallet {
	b, _ := decimal.NewFromString("10.23")
	return &entity.Wallet{
		ID:      "123",
		UserID:  "1",
		Balance: b,
	}
}

func createTestTopupWallet() *entity.TopupWallet {
	a, _ := decimal.NewFromString("10.23")
	return &entity.TopupWallet{
		WalletID: "1",
		UserID:   "2",
		Amount:   a,
	}
}

func createWalletSuite(ctrl *gomock.Controller) *WalletSuite {
	db := mock_uow.NewMockDB(ctrl)
	t := postgres.NewWallet(db)
	return &WalletSuite{
		wallet: t,
		db:     db,
	}
}
