package handler_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc/metadata"

	sdklog "github.com/indrasaputra/arjuna/pkg/sdk/log"
	apiv1 "github.com/indrasaputra/arjuna/proto/api/v1"
	"github.com/indrasaputra/arjuna/service/transaction/entity"
	"github.com/indrasaputra/arjuna/service/transaction/internal/app"
	"github.com/indrasaputra/arjuna/service/transaction/internal/grpc/handler"
	mock_service "github.com/indrasaputra/arjuna/service/transaction/test/mock/service"
)

const (
	testIdempotencyKey = "key"
	testEnv            = "development"
)

var (
	testCtxWithValidKey   = metadata.NewIncomingContext(testCtx, metadata.Pairs("X-Idempotency-Key", testIdempotencyKey))
	testCtxWithInvalidKey = metadata.NewIncomingContext(testCtx, metadata.Pairs("another-key", ""))
)

type TransactionCommandSuite struct {
	handler *handler.TransactionCommand
	creator *mock_service.MockCreateTransaction
}

func TestNewTransactionCommand(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("successful create an instance of TransactionCommand", func(t *testing.T) {
		st := createTransactionCommandSuite(ctrl)
		assert.NotNil(t, st.handler)
	})
}

func TestTransactionCommand_CreateTransaction(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	app.Logger = sdklog.NewLogger(testEnv)

	t.Run("metadata not found", func(t *testing.T) {
		st := createTransactionCommandSuite(ctrl)

		res, err := st.handler.CreateTransaction(context.Background(), nil)

		assert.Error(t, err)
		assert.Equal(t, entity.ErrInternal("metadata not found from incoming context"), err)
		assert.Nil(t, res)
	})

	t.Run("idempotency key is missing", func(t *testing.T) {
		st := createTransactionCommandSuite(ctrl)

		res, err := st.handler.CreateTransaction(testCtxWithInvalidKey, nil)

		assert.Error(t, err)
		assert.Equal(t, entity.ErrMissingIdempotencyKey(), err)
		assert.Nil(t, res)
	})

	t.Run("nil request is prohibited", func(t *testing.T) {
		st := createTransactionCommandSuite(ctrl)

		res, err := st.handler.CreateTransaction(testCtxWithValidKey, nil)

		assert.Error(t, err)
		assert.Equal(t, entity.ErrEmptyTransaction(), err)
		assert.Nil(t, res)
	})

	t.Run("empty transaction is prohibited", func(t *testing.T) {
		st := createTransactionCommandSuite(ctrl)

		res, err := st.handler.CreateTransaction(testCtxWithValidKey, &apiv1.CreateTransactionRequest{})

		assert.Error(t, err)
		assert.Equal(t, entity.ErrEmptyTransaction(), err)
		assert.Nil(t, res)
	})

	t.Run("transaction service returns error", func(t *testing.T) {
		st := createTransactionCommandSuite(ctrl)
		request := &apiv1.CreateTransactionRequest{
			Transaction: &apiv1.Transaction{
				SenderId:   uuid.Must(uuid.NewV7()).String(),
				ReceiverId: uuid.Must(uuid.NewV7()).String(),
				Amount:     "10.23",
			},
		}

		errors := []error{
			entity.ErrEmptyTransaction(),
			entity.ErrAlreadyExists(),
			entity.ErrInvalidSender(),
			entity.ErrInvalidReceiver(),
			entity.ErrInvalidAmount(),
			assert.AnError,
		}
		for _, errRet := range errors {
			st.creator.EXPECT().Create(testCtxWithValidKey, gomock.Any(), testIdempotencyKey).Return(uuid.Must(uuid.NewV7()), errRet)

			res, err := st.handler.CreateTransaction(testCtxWithValidKey, request)

			assert.Error(t, err)
			assert.Equal(t, errRet, err)
			assert.Nil(t, res)
		}
	})

	t.Run("success create transaction", func(t *testing.T) {
		st := createTransactionCommandSuite(ctrl)
		id := uuid.Must(uuid.NewV7())
		st.creator.EXPECT().Create(testCtxWithValidKey, gomock.Any(), testIdempotencyKey).Return(id, nil)
		request := &apiv1.CreateTransactionRequest{
			Transaction: &apiv1.Transaction{
				SenderId:   uuid.Must(uuid.NewV7()).String(),
				ReceiverId: uuid.Must(uuid.NewV7()).String(),
				Amount:     "10.23",
			},
		}

		res, err := st.handler.CreateTransaction(testCtxWithValidKey, request)

		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, id.String(), res.Data.GetId())
	})
}

func createTransactionCommandSuite(ctrl *gomock.Controller) *TransactionCommandSuite {
	c := mock_service.NewMockCreateTransaction(ctrl)
	h := handler.NewTransactionCommand(c)
	return &TransactionCommandSuite{
		handler: h,
		creator: c,
	}
}
