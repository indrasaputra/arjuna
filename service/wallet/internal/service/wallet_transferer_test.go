package service_test

import (
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

var (
	testErrInternal = entity.ErrInternal("")
)

type WalletTransfererSuite struct {
	wallet *service.WalletTransferer
	repo   *mock_service.MockWalletTransfererRepository
	uow    *mock_uow.MockUnitOfWork
	tx     *mock_uow.MockTx
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

	t.Run("transaction begin returns error", func(t *testing.T) {
		st := createWalletTransfererSuite(ctrl)
		trf := createTestTransferWallet("01917a52-86af-73aa-817f-46baf900d0e8", "01917a52-86af-7d6f-994f-771bcf2ffa8b")
		st.uow.EXPECT().Begin(testCtx).Return(nil, testErrInternal)

		err := st.wallet.TransferBalance(testCtx, trf)

		assert.Error(t, err)
	})

	t.Run("get sender returns error; swid < rwid", func(t *testing.T) {
		st := createWalletTransfererSuite(ctrl)
		trf := createTestTransferWallet("01917a52-86af-73aa-817f-46baf900d0e8", "01917a52-86af-7d6f-994f-771bcf2ffa8b")
		st.uow.EXPECT().Begin(testCtx).Return(st.tx, nil)
		st.repo.EXPECT().GetUserWalletWithTx(testCtx, st.tx, trf.SenderWalletID, trf.SenderID).Return(nil, testErrInternal)
		st.uow.EXPECT().Finish(testCtx, st.tx, testErrInternal)

		err := st.wallet.TransferBalance(testCtx, trf)

		assert.Error(t, err)
	})

	t.Run("get receiver returns error; swid < rwid", func(t *testing.T) {
		st := createWalletTransfererSuite(ctrl)
		trf := createTestTransferWallet("01917a52-86af-73aa-817f-46baf900d0e8", "01917a52-86af-7d6f-994f-771bcf2ffa8b")
		sw := createTestWallet()
		st.uow.EXPECT().Begin(testCtx).Return(st.tx, nil)
		st.repo.EXPECT().GetUserWalletWithTx(testCtx, st.tx, trf.SenderWalletID, trf.SenderID).Return(sw, nil)
		st.repo.EXPECT().GetUserWalletWithTx(testCtx, st.tx, trf.ReceiverID, trf.ReceiverID).Return(nil, testErrInternal)
		st.uow.EXPECT().Finish(testCtx, st.tx, testErrInternal)

		err := st.wallet.TransferBalance(testCtx, trf)

		assert.Error(t, err)
	})

	t.Run("get receiver returns error; swid >= rwid", func(t *testing.T) {
		st := createWalletTransfererSuite(ctrl)
		trf := createTestTransferWallet("01917a52-86af-7d6f-994f-771bcf2ffa8b", "01917a52-86af-73aa-817f-46baf900d0e8")
		rw := createTestWallet()
		st.uow.EXPECT().Begin(testCtx).Return(st.tx, nil)
		st.repo.EXPECT().GetUserWalletWithTx(testCtx, st.tx, trf.ReceiverWalletID, trf.ReceiverID).Return(rw, testErrInternal)
		st.uow.EXPECT().Finish(testCtx, st.tx, testErrInternal)

		err := st.wallet.TransferBalance(testCtx, trf)

		assert.Error(t, err)
	})

	t.Run("get sender returns error; swid >= rwid", func(t *testing.T) {
		st := createWalletTransfererSuite(ctrl)
		trf := createTestTransferWallet("01917a52-86af-7d6f-994f-771bcf2ffa8b", "01917a52-86af-73aa-817f-46baf900d0e8")
		rw := createTestWallet()
		st.uow.EXPECT().Begin(testCtx).Return(st.tx, nil)
		st.repo.EXPECT().GetUserWalletWithTx(testCtx, st.tx, trf.ReceiverWalletID, trf.ReceiverID).Return(rw, nil)
		st.repo.EXPECT().GetUserWalletWithTx(testCtx, st.tx, trf.SenderWalletID, trf.SenderID).Return(nil, testErrInternal)
		st.uow.EXPECT().Finish(testCtx, st.tx, testErrInternal)

		err := st.wallet.TransferBalance(testCtx, trf)

		assert.Error(t, err)
	})

	t.Run("sender is nil", func(t *testing.T) {
		st := createWalletTransfererSuite(ctrl)
		trf := createTestTransferWallet("01917a52-86af-7d6f-994f-771bcf2ffa8b", "01917a52-86af-73aa-817f-46baf900d0e8")
		rw := createTestWallet()
		st.uow.EXPECT().Begin(testCtx).Return(st.tx, nil)
		st.repo.EXPECT().GetUserWalletWithTx(testCtx, st.tx, trf.ReceiverWalletID, trf.ReceiverID).Return(rw, nil)
		st.repo.EXPECT().GetUserWalletWithTx(testCtx, st.tx, trf.SenderWalletID, trf.SenderID).Return(nil, nil)
		st.uow.EXPECT().Finish(testCtx, st.tx, entity.ErrInvalidUser())

		err := st.wallet.TransferBalance(testCtx, trf)

		assert.Error(t, err)
		assert.Equal(t, entity.ErrInvalidUser(), err)
	})

	t.Run("receiver is nil", func(t *testing.T) {
		st := createWalletTransfererSuite(ctrl)
		trf := createTestTransferWallet("01917a52-86af-7d6f-994f-771bcf2ffa8b", "01917a52-86af-73aa-817f-46baf900d0e8")
		sw := createTestWallet()
		st.uow.EXPECT().Begin(testCtx).Return(st.tx, nil)
		st.repo.EXPECT().GetUserWalletWithTx(testCtx, st.tx, trf.ReceiverWalletID, trf.ReceiverID).Return(nil, nil)
		st.repo.EXPECT().GetUserWalletWithTx(testCtx, st.tx, trf.SenderWalletID, trf.SenderID).Return(sw, nil)
		st.uow.EXPECT().Finish(testCtx, st.tx, entity.ErrInvalidUser())

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
		st.uow.EXPECT().Begin(testCtx).Return(st.tx, nil)
		st.repo.EXPECT().GetUserWalletWithTx(testCtx, st.tx, trf.SenderWalletID, trf.SenderID).Return(sw, nil)
		st.repo.EXPECT().GetUserWalletWithTx(testCtx, st.tx, trf.ReceiverWalletID, trf.ReceiverID).Return(rw, nil)
		st.uow.EXPECT().Finish(testCtx, st.tx, entity.ErrInsufficientBalance())

		err := st.wallet.TransferBalance(testCtx, trf)

		assert.Error(t, err)
		assert.Equal(t, entity.ErrInsufficientBalance(), err)
	})

	t.Run("add sender wallet returns error", func(t *testing.T) {
		st := createWalletTransfererSuite(ctrl)
		trf := createTestTransferWallet("01917a52-86af-73aa-817f-46baf900d0e8", "01917a52-86af-7d6f-994f-771bcf2ffa8b")
		sw := createTestWallet()
		rw := createTestWallet()
		st.uow.EXPECT().Begin(testCtx).Return(st.tx, nil)
		st.repo.EXPECT().GetUserWalletWithTx(testCtx, st.tx, trf.SenderWalletID, trf.SenderID).Return(sw, nil)
		st.repo.EXPECT().GetUserWalletWithTx(testCtx, st.tx, trf.ReceiverWalletID, trf.ReceiverID).Return(rw, nil)
		st.repo.EXPECT().AddWalletBalanceWithTx(testCtx, st.tx, trf.SenderWalletID, trf.Amount.Neg()).Return(testErrInternal)
		st.uow.EXPECT().Finish(testCtx, st.tx, testErrInternal)

		err := st.wallet.TransferBalance(testCtx, trf)

		assert.Error(t, err)
	})

	t.Run("add receiver wallet returns error", func(t *testing.T) {
		st := createWalletTransfererSuite(ctrl)
		trf := createTestTransferWallet("01917a52-86af-73aa-817f-46baf900d0e8", "01917a52-86af-7d6f-994f-771bcf2ffa8b")
		sw := createTestWallet()
		rw := createTestWallet()
		st.uow.EXPECT().Begin(testCtx).Return(st.tx, nil)
		st.repo.EXPECT().GetUserWalletWithTx(testCtx, st.tx, trf.SenderWalletID, trf.SenderID).Return(sw, nil)
		st.repo.EXPECT().GetUserWalletWithTx(testCtx, st.tx, trf.ReceiverWalletID, trf.ReceiverID).Return(rw, nil)
		st.repo.EXPECT().AddWalletBalanceWithTx(testCtx, st.tx, trf.SenderWalletID, trf.Amount.Neg()).Return(nil)
		st.repo.EXPECT().AddWalletBalanceWithTx(testCtx, st.tx, trf.ReceiverWalletID, trf.Amount).Return(testErrInternal)
		st.uow.EXPECT().Finish(testCtx, st.tx, testErrInternal)

		err := st.wallet.TransferBalance(testCtx, trf)

		assert.Error(t, err)
	})

	t.Run("success transfer balance", func(t *testing.T) {
		st := createWalletTransfererSuite(ctrl)
		trf := createTestTransferWallet("01917a52-86af-73aa-817f-46baf900d0e8", "01917a52-86af-7d6f-994f-771bcf2ffa8b")
		sw := createTestWallet()
		rw := createTestWallet()
		st.uow.EXPECT().Begin(testCtx).Return(st.tx, nil)
		st.repo.EXPECT().GetUserWalletWithTx(testCtx, st.tx, trf.SenderWalletID, trf.SenderID).Return(sw, nil)
		st.repo.EXPECT().GetUserWalletWithTx(testCtx, st.tx, trf.ReceiverWalletID, trf.ReceiverID).Return(rw, nil)
		st.repo.EXPECT().AddWalletBalanceWithTx(testCtx, st.tx, trf.SenderWalletID, trf.Amount.Neg()).Return(nil)
		st.repo.EXPECT().AddWalletBalanceWithTx(testCtx, st.tx, trf.ReceiverWalletID, trf.Amount).Return(nil)
		st.uow.EXPECT().Finish(testCtx, st.tx, nil)

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
	u := mock_uow.NewMockUnitOfWork(ctrl)
	t := mock_uow.NewMockTx(ctrl)
	w := service.NewWalletTransferer(r, u)
	return &WalletTransfererSuite{
		wallet: w,
		repo:   r,
		uow:    u,
		tx:     t,
	}
}
