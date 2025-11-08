package handler_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/indrasaputra/arjuna/pkg/sdk/grpc/interceptor"
	apiv1 "github.com/indrasaputra/arjuna/proto/api/v1"
	"github.com/indrasaputra/arjuna/service/wallet/entity"
	"github.com/indrasaputra/arjuna/service/wallet/internal/grpc/handler"
	mock_service "github.com/indrasaputra/arjuna/service/wallet/test/mock/service"
)

var (
	testCtxWithAuth = context.WithValue(testCtx, interceptor.HeaderKeyUserID, uuid.Must(uuid.NewV7()))
)

type WalletCommandSuite struct {
	handler  *handler.WalletCommand
	creator  *mock_service.MockCreateWallet
	topup    *mock_service.MockTopupWallet
	transfer *mock_service.MockTransferWallet
}

func TestNewWalletCommand(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("successful create an instance of WalletCommand", func(t *testing.T) {
		st := createWalletCommandSuite(ctrl)
		assert.NotNil(t, st.handler)
	})
}

func TestWalletCommand_CreateWallet(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("nil request is prohibited", func(t *testing.T) {
		st := createWalletCommandSuite(ctrl)

		res, err := st.handler.CreateWallet(testCtx, nil)

		assert.Error(t, err)
		assert.Equal(t, entity.ErrEmptyWallet(), err)
		assert.Nil(t, res)
	})

	t.Run("empty wallet is prohibited", func(t *testing.T) {
		st := createWalletCommandSuite(ctrl)

		res, err := st.handler.CreateWallet(testCtx, &apiv1.CreateWalletRequest{})

		assert.Error(t, err)
		assert.Equal(t, entity.ErrEmptyWallet(), err)
		assert.Nil(t, res)
	})

	t.Run("wallet service returns error", func(t *testing.T) {
		st := createWalletCommandSuite(ctrl)
		request := &apiv1.CreateWalletRequest{
			Wallet: &apiv1.Wallet{
				UserId:  uuid.Must(uuid.NewV7()).String(),
				Balance: "10.23",
			},
		}

		errors := []error{
			entity.ErrEmptyWallet(),
			entity.ErrAlreadyExists(),
			entity.ErrInvalidUser(),
			entity.ErrInvalidBalance(),
			entity.ErrInternal("error"),
		}
		for _, errRet := range errors {
			st.creator.EXPECT().Create(testCtx, gomock.Any()).Return(errRet)

			res, err := st.handler.CreateWallet(testCtx, request)

			assert.Error(t, err)
			assert.Equal(t, errRet, err)
			assert.Nil(t, res)
		}
	})

	t.Run("success create wallet", func(t *testing.T) {
		st := createWalletCommandSuite(ctrl)
		st.creator.EXPECT().Create(testCtx, gomock.Any()).Return(nil)
		request := &apiv1.CreateWalletRequest{
			Wallet: &apiv1.Wallet{
				UserId:  uuid.Must(uuid.NewV7()).String(),
				Balance: "10.23",
			},
		}

		res, err := st.handler.CreateWallet(testCtx, request)

		assert.NoError(t, err)
		assert.NotNil(t, res)
	})
}

func TestWalletCommand_TopupWallet(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("nil request is prohibited", func(t *testing.T) {
		st := createWalletCommandSuite(ctrl)

		res, err := st.handler.TopupWallet(testCtxWithAuth, nil)

		assert.Error(t, err)
		assert.Equal(t, entity.ErrEmptyWallet(), err)
		assert.Nil(t, res)
	})

	t.Run("empty topup is prohibited", func(t *testing.T) {
		st := createWalletCommandSuite(ctrl)

		res, err := st.handler.TopupWallet(testCtxWithAuth, &apiv1.TopupWalletRequest{})

		assert.Error(t, err)
		assert.Equal(t, entity.ErrEmptyWallet(), err)
		assert.Nil(t, res)
	})

	t.Run("topup service returns error", func(t *testing.T) {
		st := createWalletCommandSuite(ctrl)
		request := &apiv1.TopupWalletRequest{
			Topup: &apiv1.Topup{
				WalletId: uuid.Must(uuid.NewV7()).String(),
				Amount:   "10.23",
			},
		}

		errors := []error{
			entity.ErrInvalidUser(),
			entity.ErrEmptyWallet(),
			entity.ErrInvalidAmount(),
			entity.ErrInternal("error"),
		}
		for _, errRet := range errors {
			st.topup.EXPECT().Topup(testCtxWithAuth, gomock.Any()).Return(nil, errRet)

			res, err := st.handler.TopupWallet(testCtxWithAuth, request)

			assert.Error(t, err)
			assert.Equal(t, errRet, err)
			assert.Nil(t, res)
		}
	})

	t.Run("success topup wallet", func(t *testing.T) {
		walletID := uuid.Must(uuid.NewV7())
		userID := uuid.Must(uuid.NewV7())

		st := createWalletCommandSuite(ctrl)
		st.topup.EXPECT().Topup(testCtxWithAuth, gomock.Any()).Return(&entity.Wallet{
			ID:      walletID,
			UserID:  userID,
			Balance: decimal.NewFromFloat(10.23),
		}, nil)
		request := &apiv1.TopupWalletRequest{
			Topup: &apiv1.Topup{
				WalletId: walletID.String(),
				Amount:   "10.23",
			},
		}

		res, err := st.handler.TopupWallet(testCtxWithAuth, request)

		assert.NoError(t, err)
		assert.NotNil(t, res)
	})
}

func TestWalletCommand_TransferBalance(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("nil request is prohibited", func(t *testing.T) {
		st := createWalletCommandSuite(ctrl)

		res, err := st.handler.TransferBalance(testCtx, nil)

		assert.Error(t, err)
		assert.Equal(t, entity.ErrEmptyWallet(), err)
		assert.Nil(t, res)
	})

	t.Run("empty transfer is prohibited", func(t *testing.T) {
		st := createWalletCommandSuite(ctrl)

		res, err := st.handler.TransferBalance(testCtx, &apiv1.TransferBalanceRequest{})

		assert.Error(t, err)
		assert.Equal(t, entity.ErrEmptyWallet(), err)
		assert.Nil(t, res)
	})

	t.Run("wallet service returns error", func(t *testing.T) {
		st := createWalletCommandSuite(ctrl)
		request := &apiv1.TransferBalanceRequest{
			Transfer: &apiv1.Transfer{
				Amount:           "10.23",
				SenderId:         uuid.Must(uuid.NewV7()).String(),
				ReceiverId:       uuid.Must(uuid.NewV7()).String(),
				SenderWalletId:   uuid.Must(uuid.NewV7()).String(),
				ReceiverWalletId: uuid.Must(uuid.NewV7()).String(),
			},
		}

		errors := []error{
			entity.ErrEmptyWallet(),
			entity.ErrInvalidUser(),
			entity.ErrInvalidAmount(),
			entity.ErrInternal("error"),
		}
		for _, errRet := range errors {
			st.transfer.EXPECT().TransferBalance(testCtx, gomock.Any()).Return(errRet)

			res, err := st.handler.TransferBalance(testCtx, request)

			assert.Error(t, err)
			assert.Equal(t, errRet, err)
			assert.Nil(t, res)
		}
	})

	t.Run("success create wallet", func(t *testing.T) {
		st := createWalletCommandSuite(ctrl)
		st.transfer.EXPECT().TransferBalance(testCtx, gomock.Any()).Return(nil)
		request := &apiv1.TransferBalanceRequest{
			Transfer: &apiv1.Transfer{
				Amount:           "10.23",
				SenderId:         uuid.Must(uuid.NewV7()).String(),
				ReceiverId:       uuid.Must(uuid.NewV7()).String(),
				SenderWalletId:   uuid.Must(uuid.NewV7()).String(),
				ReceiverWalletId: uuid.Must(uuid.NewV7()).String(),
			},
		}

		res, err := st.handler.TransferBalance(testCtx, request)

		assert.NoError(t, err)
		assert.NotNil(t, res)
	})
}

func createWalletCommandSuite(ctrl *gomock.Controller) *WalletCommandSuite {
	c := mock_service.NewMockCreateWallet(ctrl)
	t := mock_service.NewMockTopupWallet(ctrl)
	tf := mock_service.NewMockTransferWallet(ctrl)
	h := handler.NewWalletCommand(c, t, tf)
	return &WalletCommandSuite{
		handler:  h,
		creator:  c,
		topup:    t,
		transfer: tf,
	}
}
