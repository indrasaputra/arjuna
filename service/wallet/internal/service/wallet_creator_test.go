package service_test

import (
	"context"
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
	testCtx        = context.Background()
	testUserID     = uuid.Must(uuid.NewV7())
	testBalance, _ = decimal.NewFromString("10.23")
)

type WalletCreatorSuite struct {
	wallet     *service.WalletCreator
	walletRepo *mock_service.MockCreateWalletRepository
}

func TestNewWalletCreator(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("successfully create an instance of Wallet", func(t *testing.T) {
		st := createWalletCreatorSuite(ctrl)
		assert.NotNil(t, st.wallet)
	})
}

func TestWalletCreator_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("empty wallet is prohibited", func(t *testing.T) {
		st := createWalletCreatorSuite(ctrl)

		err := st.wallet.Create(testCtx, nil)

		assert.Error(t, err)
		assert.Equal(t, entity.ErrEmptyWallet(), err)
	})

	t.Run("user id is invalid", func(t *testing.T) {
		st := createWalletCreatorSuite(ctrl)
		wallet := createTestWallet()
		wallet.UserID = uuid.Nil

		err := st.wallet.Create(testCtx, wallet)

		assert.Error(t, err)
	})

	t.Run("wallet repo insert returns error", func(t *testing.T) {
		st := createWalletCreatorSuite(ctrl)
		wallet := createTestWallet()
		st.walletRepo.EXPECT().Insert(testCtx, wallet).Return(assert.AnError)

		err := st.wallet.Create(testCtx, wallet)

		assert.Error(t, err)
	})

	t.Run("success create a wallet", func(t *testing.T) {
		st := createWalletCreatorSuite(ctrl)
		wallet := createTestWallet()
		st.walletRepo.EXPECT().Insert(testCtx, wallet).Return(nil)

		err := st.wallet.Create(testCtx, wallet)

		assert.NoError(t, err)
	})
}

func createWalletCreatorSuite(ctrl *gomock.Controller) *WalletCreatorSuite {
	r := mock_service.NewMockCreateWalletRepository(ctrl)
	w := service.NewWalletCreator(r)
	return &WalletCreatorSuite{
		wallet:     w,
		walletRepo: r,
	}
}

func createTestWallet() *entity.Wallet {
	return &entity.Wallet{
		ID:      uuid.Must(uuid.NewV7()),
		UserID:  testUserID,
		Balance: testBalance,
	}
}
