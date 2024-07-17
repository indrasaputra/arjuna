package service_test

import (
	"context"
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	sdklog "github.com/indrasaputra/arjuna/pkg/sdk/log"
	"github.com/indrasaputra/arjuna/service/transaction/entity"
	"github.com/indrasaputra/arjuna/service/transaction/internal/app"
	"github.com/indrasaputra/arjuna/service/transaction/internal/service"
	mock_service "github.com/indrasaputra/arjuna/service/transaction/test/mock/service"
)

var (
	testCtx            = context.Background()
	testSenderID       = "1"
	testReceiverID     = "2"
	testEnv            = "development"
	testAmount, _      = decimal.NewFromString("10.23")
	testIdempotencyKey = "key"
)

type TransactionSuite struct {
	trx     *service.TransactionCreator
	trxRepo *mock_service.MockCreateTransactionRepository
	keyRepo *mock_service.MockIdempotencyKeyRepository
}

func TestNewTransaction(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("successfully create an instance of Transaction", func(t *testing.T) {
		st := createTransactionSuite(ctrl)
		assert.NotNil(t, st.trx)
	})
}

func TestTransactionCreator_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	app.Logger = sdklog.NewLogger(testEnv)

	t.Run("validate idempotency key returns error", func(t *testing.T) {
		st := createTransactionSuite(ctrl)
		st.keyRepo.EXPECT().Exists(testCtx, testIdempotencyKey).Return(false, entity.ErrInternal("error"))

		id, err := st.trx.Create(testCtx, nil, testIdempotencyKey)

		assert.Error(t, err)
		assert.Empty(t, id)
	})

	t.Run("idempotency key has been used", func(t *testing.T) {
		st := createTransactionSuite(ctrl)
		st.keyRepo.EXPECT().Exists(testCtx, testIdempotencyKey).Return(true, nil)

		id, err := st.trx.Create(testCtx, nil, testIdempotencyKey)

		assert.Error(t, err)
		assert.Empty(t, id)
	})

	t.Run("empty transaction is prohibited", func(t *testing.T) {
		st := createTransactionSuite(ctrl)
		st.keyRepo.EXPECT().Exists(testCtx, testIdempotencyKey).Return(false, nil)

		id, err := st.trx.Create(testCtx, nil, testIdempotencyKey)

		assert.Error(t, err)
		assert.Equal(t, entity.ErrEmptyTransaction(), err)
		assert.Empty(t, id)
	})

	t.Run("sender id is invalid", func(t *testing.T) {
		st := createTransactionSuite(ctrl)
		trx := createTestTransaction()
		trx.SenderID = ""
		st.keyRepo.EXPECT().Exists(testCtx, testIdempotencyKey).Return(false, nil)

		id, err := st.trx.Create(testCtx, trx, testIdempotencyKey)

		assert.Error(t, err)
		assert.Empty(t, id)
	})

	t.Run("receiver id is invalid", func(t *testing.T) {
		st := createTransactionSuite(ctrl)
		trx := createTestTransaction()
		trx.ReceiverID = ""
		st.keyRepo.EXPECT().Exists(testCtx, testIdempotencyKey).Return(false, nil)

		id, err := st.trx.Create(testCtx, trx, testIdempotencyKey)

		assert.Error(t, err)
		assert.Empty(t, id)
	})

	t.Run("amount is invalid", func(t *testing.T) {
		st := createTransactionSuite(ctrl)
		trx := createTestTransaction()
		trx.Amount = decimal.Zero
		st.keyRepo.EXPECT().Exists(testCtx, testIdempotencyKey).Return(false, nil)

		id, err := st.trx.Create(testCtx, trx, testIdempotencyKey)

		assert.Error(t, err)
		assert.Empty(t, id)
	})

	t.Run("trx repo insert returns error", func(t *testing.T) {
		st := createTransactionSuite(ctrl)
		trx := createTestTransaction()
		st.keyRepo.EXPECT().Exists(testCtx, testIdempotencyKey).Return(false, nil)

		st.trxRepo.EXPECT().Insert(testCtx, trx).Return(entity.ErrInternal(""))

		id, err := st.trx.Create(testCtx, trx, testIdempotencyKey)

		assert.Error(t, err)
		assert.Empty(t, id)
	})

	t.Run("success create a transaction", func(t *testing.T) {
		st := createTransactionSuite(ctrl)
		trx := createTestTransaction()
		st.keyRepo.EXPECT().Exists(testCtx, testIdempotencyKey).Return(false, nil)

		st.trxRepo.EXPECT().Insert(testCtx, trx).Return(nil)

		id, err := st.trx.Create(testCtx, trx, testIdempotencyKey)

		assert.NoError(t, err)
		assert.NotEmpty(t, id)
	})
}

func createTransactionSuite(ctrl *gomock.Controller) *TransactionSuite {
	r := mock_service.NewMockCreateTransactionRepository(ctrl)
	i := mock_service.NewMockIdempotencyKeyRepository(ctrl)
	t := service.NewTransactionCreator(r, i)
	return &TransactionSuite{
		trx:     t,
		trxRepo: r,
		keyRepo: i,
	}
}

func createTestTransaction() *entity.Transaction {
	return &entity.Transaction{
		ID:         "1",
		SenderID:   testSenderID,
		ReceiverID: testReceiverID,
		Amount:     testAmount,
	}
}
