package handler_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	apiv1 "github.com/indrasaputra/arjuna/proto/api/v1"
	"github.com/indrasaputra/arjuna/service/transaction/entity"
	"github.com/indrasaputra/arjuna/service/transaction/internal/grpc/handler"
	mock_service "github.com/indrasaputra/arjuna/service/transaction/test/mock/service"
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

	t.Run("nil request is prohibited", func(t *testing.T) {
		st := createTransactionCommandSuite(ctrl)

		res, err := st.handler.CreateTransaction(testCtx, nil)

		assert.Error(t, err)
		assert.Equal(t, entity.ErrEmptyTransaction(), err)
		assert.Nil(t, res)
	})

	t.Run("empty transaction is prohibited", func(t *testing.T) {
		st := createTransactionCommandSuite(ctrl)

		res, err := st.handler.CreateTransaction(testCtx, &apiv1.CreateTransactionRequest{})

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
			st.creator.EXPECT().Create(testCtx, gomock.Any()).Return(uuid.Must(uuid.NewV7()), errRet)

			res, err := st.handler.CreateTransaction(testCtx, request)

			assert.Error(t, err)
			assert.Equal(t, errRet, err)
			assert.Nil(t, res)
		}
	})

	t.Run("success create transaction", func(t *testing.T) {
		st := createTransactionCommandSuite(ctrl)
		id := uuid.Must(uuid.NewV7())
		st.creator.EXPECT().Create(testCtx, gomock.Any()).Return(id, nil)
		request := &apiv1.CreateTransactionRequest{
			Transaction: &apiv1.Transaction{
				SenderId:   uuid.Must(uuid.NewV7()).String(),
				ReceiverId: uuid.Must(uuid.NewV7()).String(),
				Amount:     "10.23",
			},
		}

		res, err := st.handler.CreateTransaction(testCtx, request)

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
