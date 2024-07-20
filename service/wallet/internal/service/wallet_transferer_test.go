package service_test

import (
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	sdklog "github.com/indrasaputra/arjuna/pkg/sdk/log"
	"github.com/indrasaputra/arjuna/service/wallet/entity"
	"github.com/indrasaputra/arjuna/service/wallet/internal/app"
	"github.com/indrasaputra/arjuna/service/wallet/internal/service"
	mock_service "github.com/indrasaputra/arjuna/service/wallet/test/mock/service"
)

type WalletTransfererSuite struct {
	wallet *service.WalletTransferer
	repo   *mock_service.MockWalletTransfererRepository
}

func TestNewWalletTransferer(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("successfully create an instance of WalletTransferer", func(t *testing.T) {
		st := createWalletTransfererSuite(ctrl)
		assert.NotNil(t, st.wallet)
	})
}

func TestWalletTransferer_TransferBalance(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	app.Logger = sdklog.NewLogger(testEnv)

	t.Run("transfer is nil", func(t *testing.T) {
		st := createWalletTransfererSuite(ctrl)

		err := st.wallet.TransferBalance(testCtx, nil)

		assert.Error(t, err)
		assert.Equal(t, entity.ErrInvalidTransfer(), err)
	})

	t.Run("sender and receiver are same", func(t *testing.T) {
		st := createWalletTransfererSuite(ctrl)
		trf := createTestTransferWallet("1", "2")
		trf.ReceiverID = trf.SenderID

		err := st.wallet.TransferBalance(testCtx, trf)

		assert.Error(t, err)
		assert.Equal(t, entity.ErrSameAccount(), err)
	})

	t.Run("invalid amount", func(t *testing.T) {
		st := createWalletTransfererSuite(ctrl)
		trf := createTestTransferWallet("1", "2")
		trf.Amount = decimal.Zero

		err := st.wallet.TransferBalance(testCtx, trf)

		assert.Error(t, err)
		assert.Equal(t, entity.ErrInvalidAmount(), err)
	})

	t.Run("get sender returns error; swid < rwid", func(t *testing.T) {
		st := createWalletTransfererSuite(ctrl)
		trf := createTestTransferWallet("1", "2")
		st.repo.EXPECT().GetUserWallet(testCtx, trf.SenderWalletID, trf.SenderID).Return(nil, entity.ErrInternal(""))

		err := st.wallet.TransferBalance(testCtx, trf)

		assert.Error(t, err)
	})

	t.Run("get receiver returns error; swid < rwid", func(t *testing.T) {
		st := createWalletTransfererSuite(ctrl)
		trf := createTestTransferWallet("1", "2")
		sw := createTestWallet()
		st.repo.EXPECT().GetUserWallet(testCtx, trf.SenderWalletID, trf.SenderID).Return(sw, nil)
		st.repo.EXPECT().GetUserWallet(testCtx, trf.ReceiverID, trf.ReceiverID).Return(nil, entity.ErrInternal(""))

		err := st.wallet.TransferBalance(testCtx, trf)

		assert.Error(t, err)
	})

	t.Run("get receiver returns error; swid >= rwid", func(t *testing.T) {
		st := createWalletTransfererSuite(ctrl)
		trf := createTestTransferWallet("2", "1")
		rw := createTestWallet()
		st.repo.EXPECT().GetUserWallet(testCtx, trf.ReceiverWalletID, trf.ReceiverID).Return(rw, entity.ErrInternal(""))

		err := st.wallet.TransferBalance(testCtx, trf)

		assert.Error(t, err)
	})

	t.Run("get sender returns error; swid >= rwid", func(t *testing.T) {
		st := createWalletTransfererSuite(ctrl)
		trf := createTestTransferWallet("2", "1")
		rw := createTestWallet()
		st.repo.EXPECT().GetUserWallet(testCtx, trf.ReceiverWalletID, trf.ReceiverID).Return(rw, nil)
		st.repo.EXPECT().GetUserWallet(testCtx, trf.SenderWalletID, trf.SenderID).Return(nil, entity.ErrInternal(""))

		err := st.wallet.TransferBalance(testCtx, trf)

		assert.Error(t, err)
	})

	t.Run("sender is nil", func(t *testing.T) {
		st := createWalletTransfererSuite(ctrl)
		trf := createTestTransferWallet("2", "1")
		rw := createTestWallet()
		st.repo.EXPECT().GetUserWallet(testCtx, trf.ReceiverWalletID, trf.ReceiverID).Return(rw, nil)
		st.repo.EXPECT().GetUserWallet(testCtx, trf.SenderWalletID, trf.SenderID).Return(nil, nil)

		err := st.wallet.TransferBalance(testCtx, trf)

		assert.Error(t, err)
		assert.Equal(t, entity.ErrInvalidUser(), err)
	})

	t.Run("receiver is nil", func(t *testing.T) {
		st := createWalletTransfererSuite(ctrl)
		trf := createTestTransferWallet("2", "1")
		sw := createTestWallet()
		st.repo.EXPECT().GetUserWallet(testCtx, trf.ReceiverWalletID, trf.ReceiverID).Return(nil, nil)
		st.repo.EXPECT().GetUserWallet(testCtx, trf.SenderWalletID, trf.SenderID).Return(sw, nil)

		err := st.wallet.TransferBalance(testCtx, trf)

		assert.Error(t, err)
		assert.Equal(t, entity.ErrInvalidUser(), err)
	})

	t.Run("sender balance is insufficient", func(t *testing.T) {
		st := createWalletTransfererSuite(ctrl)
		trf := createTestTransferWallet("1", "2")
		trf.Amount, _ = decimal.NewFromString("100.98")
		sw := createTestWallet()
		rw := createTestWallet()
		st.repo.EXPECT().GetUserWallet(testCtx, trf.SenderWalletID, trf.SenderID).Return(sw, nil)
		st.repo.EXPECT().GetUserWallet(testCtx, trf.ReceiverWalletID, trf.ReceiverID).Return(rw, nil)

		err := st.wallet.TransferBalance(testCtx, trf)

		assert.Error(t, err)
		assert.Equal(t, entity.ErrInsufficientBalance(), err)
	})

	t.Run("add sender wallet returns error", func(t *testing.T) {
		st := createWalletTransfererSuite(ctrl)
		trf := createTestTransferWallet("1", "2")
		sw := createTestWallet()
		rw := createTestWallet()
		st.repo.EXPECT().GetUserWallet(testCtx, trf.SenderWalletID, trf.SenderID).Return(sw, nil)
		st.repo.EXPECT().GetUserWallet(testCtx, trf.ReceiverWalletID, trf.ReceiverID).Return(rw, nil)
		st.repo.EXPECT().AddWalletBalance(testCtx, trf.SenderWalletID, trf.Amount.Neg()).Return(entity.ErrInternal(""))

		err := st.wallet.TransferBalance(testCtx, trf)

		assert.Error(t, err)
	})

	t.Run("add receiver wallet returns error", func(t *testing.T) {
		st := createWalletTransfererSuite(ctrl)
		trf := createTestTransferWallet("1", "2")
		sw := createTestWallet()
		rw := createTestWallet()
		st.repo.EXPECT().GetUserWallet(testCtx, trf.SenderWalletID, trf.SenderID).Return(sw, nil)
		st.repo.EXPECT().GetUserWallet(testCtx, trf.ReceiverWalletID, trf.ReceiverID).Return(rw, nil)
		st.repo.EXPECT().AddWalletBalance(testCtx, trf.SenderWalletID, trf.Amount.Neg()).Return(nil)
		st.repo.EXPECT().AddWalletBalance(testCtx, trf.ReceiverWalletID, trf.Amount).Return(entity.ErrInternal(""))

		err := st.wallet.TransferBalance(testCtx, trf)

		assert.Error(t, err)
	})

	t.Run("success transfer balance", func(t *testing.T) {
		st := createWalletTransfererSuite(ctrl)
		trf := createTestTransferWallet("1", "2")
		sw := createTestWallet()
		rw := createTestWallet()
		st.repo.EXPECT().GetUserWallet(testCtx, trf.SenderWalletID, trf.SenderID).Return(sw, nil)
		st.repo.EXPECT().GetUserWallet(testCtx, trf.ReceiverWalletID, trf.ReceiverID).Return(rw, nil)
		st.repo.EXPECT().AddWalletBalance(testCtx, trf.SenderWalletID, trf.Amount.Neg()).Return(nil)
		st.repo.EXPECT().AddWalletBalance(testCtx, trf.ReceiverWalletID, trf.Amount).Return(nil)

		err := st.wallet.TransferBalance(testCtx, trf)

		assert.NoError(t, err)
	})
}
func createTestTransferWallet(swid, rwid string) *entity.TransferWallet {
	amount, _ := decimal.NewFromString("3.4")
	return &entity.TransferWallet{
		SenderID:         "1",
		SenderWalletID:   swid,
		ReceiverID:       "2",
		ReceiverWalletID: rwid,
		Amount:           amount,
	}
}

func createWalletTransfererSuite(ctrl *gomock.Controller) *WalletTransfererSuite {
	r := mock_service.NewMockWalletTransfererRepository(ctrl)
	w := service.NewWalletTransferer(r)
	return &WalletTransfererSuite{
		wallet: w,
		repo:   r,
	}
}
