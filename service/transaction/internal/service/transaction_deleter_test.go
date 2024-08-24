package service_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	sdklog "github.com/indrasaputra/arjuna/pkg/sdk/log"
	"github.com/indrasaputra/arjuna/service/transaction/entity"
	"github.com/indrasaputra/arjuna/service/transaction/internal/app"
	"github.com/indrasaputra/arjuna/service/transaction/internal/service"
	mock_service "github.com/indrasaputra/arjuna/service/transaction/test/mock/service"
)

type TransactionDeleterSuite struct {
	trx  *service.TransactionDeleter
	repo *mock_service.MockDeleteTransactionRepository
}

func TestNewTransactionDeleter(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("successfully create an instance of TransactionDeleter", func(t *testing.T) {
		st := createTransactionDeleterSuite(ctrl)
		assert.NotNil(t, st.trx)
	})
}

func TestTransactionDeleter_DeleteAllTransactions(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	app.Logger = sdklog.NewLogger(testEnv)

	t.Run("delete all transactions returns error", func(t *testing.T) {
		st := createTransactionDeleterSuite(ctrl)
		st.repo.EXPECT().DeleteAll(testCtx).Return(entity.ErrInternal(""))

		err := st.trx.DeleteAllTransactions(testCtx)

		assert.Error(t, err)
	})

	t.Run("delete all transactions returns success", func(t *testing.T) {
		st := createTransactionDeleterSuite(ctrl)
		st.repo.EXPECT().DeleteAll(testCtx).Return(nil)

		err := st.trx.DeleteAllTransactions(testCtx)

		assert.NoError(t, err)
	})
}

func createTransactionDeleterSuite(ctrl *gomock.Controller) *TransactionDeleterSuite {
	r := mock_service.NewMockDeleteTransactionRepository(ctrl)
	t := service.NewTransactionDeleter(r)
	return &TransactionDeleterSuite{
		trx:  t,
		repo: r,
	}
}
