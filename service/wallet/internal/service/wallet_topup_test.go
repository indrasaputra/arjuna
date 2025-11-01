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
	testWalletID       = uuid.Must(uuid.NewV7())
	testAmount, _      = decimal.NewFromString("10.23")
	testIdempotencyKey = "idempotency-key"
)

type WalletTopupSuite struct {
	topup     *service.WalletTopup
	topupRepo *mock_service.MockTopupWalletRepository
	keyRepo   *mock_service.MockIdempotencyKeyRepository
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

	t.Run("validate idempotency key returns error", func(t *testing.T) {
		st := createWalletTopupSuite(ctrl)
		topup := createTestTopupWallet()
		st.keyRepo.EXPECT().Exists(testCtx, testIdempotencyKey).Return(false, assert.AnError)

		wallet, err := st.topup.Topup(testCtx, topup)

		assert.Error(t, err)
		assert.Nil(t, wallet)
	})

	t.Run("idempotency key has been used", func(t *testing.T) {
		st := createWalletTopupSuite(ctrl)
		topup := createTestTopupWallet()
		st.keyRepo.EXPECT().Exists(testCtx, testIdempotencyKey).Return(true, nil)

		wallet, err := st.topup.Topup(testCtx, topup)

		assert.Error(t, err)
		assert.Nil(t, wallet)
	})

	t.Run("wallet id is invalid", func(t *testing.T) {
		st := createWalletTopupSuite(ctrl)
		topup := createTestTopupWallet()
		topup.WalletID = uuid.Nil
		st.keyRepo.EXPECT().Exists(testCtx, testIdempotencyKey).Return(false, nil)

		wallet, err := st.topup.Topup(testCtx, topup)

		assert.Error(t, err)
		assert.Nil(t, wallet)
	})

	t.Run("user id is invalid", func(t *testing.T) {
		st := createWalletTopupSuite(ctrl)
		topup := createTestTopupWallet()
		topup.UserID = uuid.Nil
		st.keyRepo.EXPECT().Exists(testCtx, testIdempotencyKey).Return(false, nil)

		wallet, err := st.topup.Topup(testCtx, topup)

		assert.Error(t, err)
		assert.Nil(t, wallet)
	})

	t.Run("amount is invalid", func(t *testing.T) {
		st := createWalletTopupSuite(ctrl)
		topup := createTestTopupWallet()
		topup.Amount = decimal.Zero
		st.keyRepo.EXPECT().Exists(testCtx, testIdempotencyKey).Return(false, nil)

		wallet, err := st.topup.Topup(testCtx, topup)

		assert.Error(t, err)
		assert.Nil(t, wallet)
	})

	t.Run("wallet repo update balance returns error", func(t *testing.T) {
		st := createWalletTopupSuite(ctrl)
		topup := createTestTopupWallet()
		st.keyRepo.EXPECT().Exists(testCtx, testIdempotencyKey).Return(false, nil)
		st.topupRepo.EXPECT().AddWalletBalance(testCtx, topup.WalletID, topup.Amount).Return(nil, assert.AnError)

		wallet, err := st.topup.Topup(testCtx, topup)

		assert.Error(t, err)
		assert.Nil(t, wallet)
	})

	t.Run("success create a topup", func(t *testing.T) {
		st := createWalletTopupSuite(ctrl)
		topup := createTestTopupWallet()
		st.keyRepo.EXPECT().Exists(testCtx, testIdempotencyKey).Return(false, nil)
		st.topupRepo.EXPECT().AddWalletBalance(testCtx, topup.WalletID, topup.Amount).Return(createTestWallet(), nil)

		wallet, err := st.topup.Topup(testCtx, topup)

		assert.NoError(t, err)
		assert.NotNil(t, wallet)
	})
}

func createWalletTopupSuite(ctrl *gomock.Controller) *WalletTopupSuite {
	r := mock_service.NewMockTopupWalletRepository(ctrl)
	i := mock_service.NewMockIdempotencyKeyRepository(ctrl)
	t := service.NewWalletTopup(r, i)
	return &WalletTopupSuite{
		topup:     t,
		topupRepo: r,
		keyRepo:   i,
	}
}

func createTestTopupWallet() *entity.TopupWallet {
	return &entity.TopupWallet{
		WalletID:       testWalletID,
		UserID:         testUserID,
		Amount:         testAmount,
		IdempotencyKey: testIdempotencyKey,
	}
}
