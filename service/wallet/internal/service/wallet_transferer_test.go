package service_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	sdklog "github.com/indrasaputra/arjuna/pkg/sdk/log"
	mock_uow "github.com/indrasaputra/arjuna/pkg/sdk/test/mock/uow"
	"github.com/indrasaputra/arjuna/service/wallet/entity"
	"github.com/indrasaputra/arjuna/service/wallet/internal/app"
	"github.com/indrasaputra/arjuna/service/wallet/internal/service"
	mock_service "github.com/indrasaputra/arjuna/service/wallet/test/mock/service"
)

type ctxKey string

var testCtxTx = context.WithValue(testCtx, ctxKey("tx"), true)

type WalletTransfererSuite struct {
	wallet    *service.WalletTransferer
	repo      *mock_service.MockWalletTransfererRepository
	txManager *mock_uow.MockTxManager
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
		trf := createTestTransferWallet("01917a52-86af-73aa-817f-46baf900d0e8", "01917a52-86af-7d6f-994f-771bcf2ffa8b")
		trf.ReceiverID = trf.SenderID

		err := st.wallet.TransferBalance(testCtx, trf)

		assert.Error(t, err)
		assert.Equal(t, entity.ErrSameAccount(), err)
	})

	t.Run("invalid amount", func(t *testing.T) {
		st := createWalletTransfererSuite(ctrl)
		trf := createTestTransferWallet("01917a52-86af-73aa-817f-46baf900d0e8", "01917a52-86af-7d6f-994f-771bcf2ffa8b")
		trf.Amount = decimal.Zero

		err := st.wallet.TransferBalance(testCtx, trf)

		assert.Error(t, err)
		assert.Equal(t, entity.ErrInvalidAmount(), err)
	})

	t.Run("get sender returns error; swid < rwid", func(t *testing.T) {
		st := createWalletTransfererSuite(ctrl)
		trf := createTestTransferWallet("01917a52-86af-73aa-817f-46baf900d0e8", "01917a52-86af-7d6f-994f-771bcf2ffa8b")
		st.repo.EXPECT().GetUserWalletForUpdate(testCtxTx, trf.SenderWalletID, trf.SenderID).Return(nil, assert.AnError)
		st.txManager.EXPECT().Do(gomock.Any(), gomock.Any()).
			DoAndReturn(func(_ context.Context, fn func(context.Context) error) error {
				assert.Error(t, fn(testCtxTx))
				return assert.AnError
			})

		err := st.wallet.TransferBalance(testCtx, trf)

		assert.Error(t, err)
	})

	t.Run("get receiver returns error; swid < rwid", func(t *testing.T) {
		st := createWalletTransfererSuite(ctrl)
		trf := createTestTransferWallet("01917a52-86af-73aa-817f-46baf900d0e8", "01917a52-86af-7d6f-994f-771bcf2ffa8b")
		sw := createTestWallet()
		st.repo.EXPECT().GetUserWalletForUpdate(testCtxTx, trf.SenderWalletID, trf.SenderID).Return(sw, nil)
		st.repo.EXPECT().GetUserWalletForUpdate(testCtxTx, trf.ReceiverID, trf.ReceiverID).Return(nil, assert.AnError)
		st.txManager.EXPECT().Do(gomock.Any(), gomock.Any()).
			DoAndReturn(func(_ context.Context, fn func(context.Context) error) error {
				assert.Error(t, fn(testCtxTx))
				return assert.AnError
			})

		err := st.wallet.TransferBalance(testCtx, trf)

		assert.Error(t, err)
	})

	t.Run("get receiver returns error; swid >= rwid", func(t *testing.T) {
		st := createWalletTransfererSuite(ctrl)
		trf := createTestTransferWallet("01917a52-86af-7d6f-994f-771bcf2ffa8b", "01917a52-86af-73aa-817f-46baf900d0e8")
		rw := createTestWallet()
		st.repo.EXPECT().GetUserWalletForUpdate(testCtxTx, trf.ReceiverWalletID, trf.ReceiverID).Return(rw, assert.AnError)
		st.txManager.EXPECT().Do(gomock.Any(), gomock.Any()).
			DoAndReturn(func(_ context.Context, fn func(context.Context) error) error {
				assert.Error(t, fn(testCtxTx))
				return assert.AnError
			})

		err := st.wallet.TransferBalance(testCtx, trf)

		assert.Error(t, err)
	})

	t.Run("get sender returns error; swid >= rwid", func(t *testing.T) {
		st := createWalletTransfererSuite(ctrl)
		trf := createTestTransferWallet("01917a52-86af-7d6f-994f-771bcf2ffa8b", "01917a52-86af-73aa-817f-46baf900d0e8")
		rw := createTestWallet()
		st.repo.EXPECT().GetUserWalletForUpdate(testCtxTx, trf.ReceiverWalletID, trf.ReceiverID).Return(rw, nil)
		st.repo.EXPECT().GetUserWalletForUpdate(testCtxTx, trf.SenderWalletID, trf.SenderID).Return(nil, assert.AnError)
		st.txManager.EXPECT().Do(gomock.Any(), gomock.Any()).
			DoAndReturn(func(_ context.Context, fn func(context.Context) error) error {
				assert.Error(t, fn(testCtxTx))
				return assert.AnError
			})

		err := st.wallet.TransferBalance(testCtx, trf)

		assert.Error(t, err)
	})

	t.Run("sender is nil", func(t *testing.T) {
		st := createWalletTransfererSuite(ctrl)
		trf := createTestTransferWallet("01917a52-86af-7d6f-994f-771bcf2ffa8b", "01917a52-86af-73aa-817f-46baf900d0e8")
		rw := createTestWallet()
		st.repo.EXPECT().GetUserWalletForUpdate(testCtxTx, trf.ReceiverWalletID, trf.ReceiverID).Return(rw, nil)
		st.repo.EXPECT().GetUserWalletForUpdate(testCtxTx, trf.SenderWalletID, trf.SenderID).Return(nil, nil)
		st.txManager.EXPECT().Do(gomock.Any(), gomock.Any()).
			DoAndReturn(func(_ context.Context, fn func(context.Context) error) error {
				assert.Error(t, fn(testCtxTx))
				return entity.ErrInvalidUser()
			})

		err := st.wallet.TransferBalance(testCtx, trf)

		assert.Error(t, err)
		assert.Equal(t, entity.ErrInvalidUser(), err)
	})

	t.Run("receiver is nil", func(t *testing.T) {
		st := createWalletTransfererSuite(ctrl)
		trf := createTestTransferWallet("01917a52-86af-7d6f-994f-771bcf2ffa8b", "01917a52-86af-73aa-817f-46baf900d0e8")
		sw := createTestWallet()
		st.repo.EXPECT().GetUserWalletForUpdate(testCtxTx, trf.ReceiverWalletID, trf.ReceiverID).Return(nil, nil)
		st.repo.EXPECT().GetUserWalletForUpdate(testCtxTx, trf.SenderWalletID, trf.SenderID).Return(sw, nil)
		st.txManager.EXPECT().Do(gomock.Any(), gomock.Any()).
			DoAndReturn(func(_ context.Context, fn func(context.Context) error) error {
				assert.Error(t, fn(testCtxTx))
				return entity.ErrInvalidUser()
			})

		err := st.wallet.TransferBalance(testCtx, trf)

		assert.Error(t, err)
		assert.Equal(t, entity.ErrInvalidUser(), err)
	})

	t.Run("sender balance is insufficient", func(t *testing.T) {
		st := createWalletTransfererSuite(ctrl)
		trf := createTestTransferWallet("01917a52-86af-73aa-817f-46baf900d0e8", "01917a52-86af-7d6f-994f-771bcf2ffa8b")
		trf.Amount, _ = decimal.NewFromString("100.98")
		sw := createTestWallet()
		rw := createTestWallet()
		st.repo.EXPECT().GetUserWalletForUpdate(testCtxTx, trf.SenderWalletID, trf.SenderID).Return(sw, nil)
		st.repo.EXPECT().GetUserWalletForUpdate(testCtxTx, trf.ReceiverWalletID, trf.ReceiverID).Return(rw, nil)
		st.txManager.EXPECT().Do(gomock.Any(), gomock.Any()).
			DoAndReturn(func(_ context.Context, fn func(context.Context) error) error {
				assert.Error(t, fn(testCtxTx))
				return entity.ErrInsufficientBalance()
			})

		err := st.wallet.TransferBalance(testCtx, trf)

		assert.Error(t, err)
		assert.Equal(t, entity.ErrInsufficientBalance(), err)
	})

	t.Run("add sender wallet returns error", func(t *testing.T) {
		st := createWalletTransfererSuite(ctrl)
		trf := createTestTransferWallet("01917a52-86af-73aa-817f-46baf900d0e8", "01917a52-86af-7d6f-994f-771bcf2ffa8b")
		sw := createTestWallet()
		rw := createTestWallet()
		st.repo.EXPECT().GetUserWalletForUpdate(testCtxTx, trf.SenderWalletID, trf.SenderID).Return(sw, nil)
		st.repo.EXPECT().GetUserWalletForUpdate(testCtxTx, trf.ReceiverWalletID, trf.ReceiverID).Return(rw, nil)
		st.repo.EXPECT().AddWalletBalance(testCtxTx, trf.SenderWalletID, trf.Amount.Neg()).Return(assert.AnError)
		st.txManager.EXPECT().Do(gomock.Any(), gomock.Any()).
			DoAndReturn(func(_ context.Context, fn func(context.Context) error) error {
				assert.Error(t, fn(testCtxTx))
				return assert.AnError
			})

		err := st.wallet.TransferBalance(testCtx, trf)

		assert.Error(t, err)
	})

	t.Run("add receiver wallet returns error", func(t *testing.T) {
		st := createWalletTransfererSuite(ctrl)
		trf := createTestTransferWallet("01917a52-86af-73aa-817f-46baf900d0e8", "01917a52-86af-7d6f-994f-771bcf2ffa8b")
		sw := createTestWallet()
		rw := createTestWallet()
		st.repo.EXPECT().GetUserWalletForUpdate(testCtxTx, trf.SenderWalletID, trf.SenderID).Return(sw, nil)
		st.repo.EXPECT().GetUserWalletForUpdate(testCtxTx, trf.ReceiverWalletID, trf.ReceiverID).Return(rw, nil)
		st.repo.EXPECT().AddWalletBalance(testCtxTx, trf.SenderWalletID, trf.Amount.Neg()).Return(nil)
		st.repo.EXPECT().AddWalletBalance(testCtxTx, trf.ReceiverWalletID, trf.Amount).Return(assert.AnError)
		st.txManager.EXPECT().Do(gomock.Any(), gomock.Any()).
			DoAndReturn(func(_ context.Context, fn func(context.Context) error) error {
				assert.Error(t, fn(testCtxTx))
				return assert.AnError
			})

		err := st.wallet.TransferBalance(testCtx, trf)

		assert.Error(t, err)
	})

	t.Run("tx manager returns error", func(t *testing.T) {
		st := createWalletTransfererSuite(ctrl)
		trf := createTestTransferWallet("01917a52-86af-73aa-817f-46baf900d0e8", "01917a52-86af-7d6f-994f-771bcf2ffa8b")
		sw := createTestWallet()
		rw := createTestWallet()
		st.repo.EXPECT().GetUserWalletForUpdate(testCtxTx, trf.SenderWalletID, trf.SenderID).Return(sw, nil)
		st.repo.EXPECT().GetUserWalletForUpdate(testCtxTx, trf.ReceiverWalletID, trf.ReceiverID).Return(rw, nil)
		st.repo.EXPECT().AddWalletBalance(testCtxTx, trf.SenderWalletID, trf.Amount.Neg()).Return(nil)
		st.repo.EXPECT().AddWalletBalance(testCtxTx, trf.ReceiverWalletID, trf.Amount).Return(nil)
		st.txManager.EXPECT().Do(gomock.Any(), gomock.Any()).
			DoAndReturn(func(_ context.Context, fn func(context.Context) error) error {
				assert.NoError(t, fn(testCtxTx))
				return assert.AnError
			})

		err := st.wallet.TransferBalance(testCtx, trf)

		assert.Error(t, err)
	})

	t.Run("success transfer balance", func(t *testing.T) {
		st := createWalletTransfererSuite(ctrl)
		trf := createTestTransferWallet("01917a52-86af-73aa-817f-46baf900d0e8", "01917a52-86af-7d6f-994f-771bcf2ffa8b")
		sw := createTestWallet()
		rw := createTestWallet()
		st.repo.EXPECT().GetUserWalletForUpdate(testCtxTx, trf.SenderWalletID, trf.SenderID).Return(sw, nil)
		st.repo.EXPECT().GetUserWalletForUpdate(testCtxTx, trf.ReceiverWalletID, trf.ReceiverID).Return(rw, nil)
		st.repo.EXPECT().AddWalletBalance(testCtxTx, trf.SenderWalletID, trf.Amount.Neg()).Return(nil)
		st.repo.EXPECT().AddWalletBalance(testCtxTx, trf.ReceiverWalletID, trf.Amount).Return(nil)
		st.txManager.EXPECT().Do(gomock.Any(), gomock.Any()).
			DoAndReturn(func(_ context.Context, fn func(context.Context) error) error {
				assert.NoError(t, fn(testCtxTx))
				return nil
			})

		err := st.wallet.TransferBalance(testCtx, trf)

		assert.NoError(t, err)
	})
}

func createTestTransferWallet(swid, rwid string) *entity.TransferWallet {
	amount, _ := decimal.NewFromString("3.4")
	return &entity.TransferWallet{
		SenderID:         uuid.MustParse("01917a52-86af-73aa-817f-46baf900d0e8"),
		SenderWalletID:   uuid.MustParse(swid),
		ReceiverID:       uuid.MustParse("01917a52-86af-7d6f-994f-771bcf2ffa8b"),
		ReceiverWalletID: uuid.MustParse(rwid),
		Amount:           amount,
	}
}

func createWalletTransfererSuite(ctrl *gomock.Controller) *WalletTransfererSuite {
	r := mock_service.NewMockWalletTransfererRepository(ctrl)
	m := mock_uow.NewMockTxManager(ctrl)
	w := service.NewWalletTransferer(r, m)
	return &WalletTransfererSuite{
		wallet:    w,
		repo:      r,
		txManager: m,
	}
}
