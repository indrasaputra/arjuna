package handler_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	sdklog "github.com/indrasaputra/arjuna/pkg/sdk/log"
	apiv1 "github.com/indrasaputra/arjuna/proto/api/v1"
	"github.com/indrasaputra/arjuna/service/wallet/entity"
	"github.com/indrasaputra/arjuna/service/wallet/internal/app"
	"github.com/indrasaputra/arjuna/service/wallet/internal/grpc/handler"
	mock_service "github.com/indrasaputra/arjuna/service/wallet/test/mock/service"
)

const (
	testIdempotencyKey = "key"
	testEnv            = "development"
)

type WalletCommandSuite struct {
	handler *handler.WalletCommand
	creator *mock_service.MockCreateWallet
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
	app.Logger = sdklog.NewLogger(testEnv)

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
				UserId:  "1",
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
				UserId:  "1",
				Balance: "10.23",
			},
		}

		res, err := st.handler.CreateWallet(testCtx, request)

		assert.NoError(t, err)
		assert.NotNil(t, res)
	})
}

func createWalletCommandSuite(ctrl *gomock.Controller) *WalletCommandSuite {
	w := mock_service.NewMockCreateWallet(ctrl)
	h := handler.NewWalletCommand(w)
	return &WalletCommandSuite{
		handler: h,
		creator: w,
	}
}
