package service_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/indrasaputra/arjuna/service/wallet/entity"
	"github.com/indrasaputra/arjuna/service/wallet/internal/service"
	mock_service "github.com/indrasaputra/arjuna/service/wallet/test/mock/service"
)

var (
	testWalletID  = uuid.Must(uuid.NewV7())
	testAmount, _ = decimal.NewFromString("10.23")
)

type WalletTopupSuite struct {
	topup     *service.WalletTopup
	topupRepo *mock_service.MockTopupWalletRepository
}

func TestNewWalletTopup(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("successfully create an instance of Wallet", func(t *testing.T) {
		st := createWalletTopupSuite(ctrl)
		assert.NotNil(t, st.topup)
	})
}

func TestWalletTopup_Topup(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("empty wallet is prohibited", func(t *testing.T) {
		st := createWalletTopupSuite(ctrl)

		wallet, err := st.topup.Topup(testCtx, nil)

		assert.Error(t, err)
		assert.Equal(t, entity.ErrEmptyWallet(), err)
		assert.Nil(t, wallet)
	})

	t.Run("wallet id is invalid", func(t *testing.T) {
		st := createWalletTopupSuite(ctrl)
		topup := createTestTopupWallet()
		topup.WalletID = uuid.Nil

		wallet, err := st.topup.Topup(testCtx, topup)

		assert.Error(t, err)
		assert.Nil(t, wallet)
	})

	t.Run("user id is invalid", func(t *testing.T) {
		st := createWalletTopupSuite(ctrl)
		topup := createTestTopupWallet()
		topup.UserID = uuid.Nil

		wallet, err := st.topup.Topup(testCtx, topup)

		assert.Error(t, err)
		assert.Nil(t, wallet)
	})

	t.Run("amount is invalid", func(t *testing.T) {
		st := createWalletTopupSuite(ctrl)
		topup := createTestTopupWallet()
		topup.Amount = decimal.Zero

		wallet, err := st.topup.Topup(testCtx, topup)

		assert.Error(t, err)
		assert.Nil(t, wallet)
	})

	t.Run("wallet repo update balance returns error", func(t *testing.T) {
		st := createWalletTopupSuite(ctrl)
		topup := createTestTopupWallet()
		st.topupRepo.EXPECT().AddWalletBalance(testCtx, topup.WalletID, topup.Amount).Return(nil, assert.AnError)

		wallet, err := st.topup.Topup(testCtx, topup)

		assert.Error(t, err)
		assert.Nil(t, wallet)
	})

	t.Run("success create a topup", func(t *testing.T) {
		st := createWalletTopupSuite(ctrl)
		topup := createTestTopupWallet()
		st.topupRepo.EXPECT().AddWalletBalance(testCtx, topup.WalletID, topup.Amount).Return(createTestWallet(), nil)

		wallet, err := st.topup.Topup(testCtx, topup)

		assert.NoError(t, err)
		assert.NotNil(t, wallet)
	})
}

func createWalletTopupSuite(ctrl *gomock.Controller) *WalletTopupSuite {
	r := mock_service.NewMockTopupWalletRepository(ctrl)
	t := service.NewWalletTopup(r)
	return &WalletTopupSuite{
		topup:     t,
		topupRepo: r,
	}
}

func createTestTopupWallet() *entity.TopupWallet {
	return &entity.TopupWallet{
		WalletID: testWalletID,
		UserID:   testUserID,
		Amount:   testAmount,
	}
}
