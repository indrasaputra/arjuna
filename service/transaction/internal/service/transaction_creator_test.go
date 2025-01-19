package service_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/indrasaputra/arjuna/service/transaction/entity"
	"github.com/indrasaputra/arjuna/service/transaction/internal/service"
	mock_service "github.com/indrasaputra/arjuna/service/transaction/test/mock/service"
)

var (
	testCtx            = context.Background()
	testSenderID       = uuid.Must(uuid.NewV7())
	testReceiverID     = uuid.Must(uuid.NewV7())
	testAmount, _      = decimal.NewFromString("10.23")
	testIdempotencyKey = "key"
)

type TransactionCreatorSuite struct {
	trx     *service.TransactionCreator
	trxRepo *mock_service.MockCreateTransactionRepository
	keyRepo *mock_service.MockIdempotencyKeyRepository
}

func TestNewTransactionCreator(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("successfully create an instance of TransactionCreator", func(t *testing.T) {
		st := createTransactionCreatorSuite(ctrl)
		assert.NotNil(t, st.trx)
	})
}

func TestTransactionCreator_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("validate idempotency key returns error", func(t *testing.T) {
		st := createTransactionCreatorSuite(ctrl)
		st.keyRepo.EXPECT().Exists(testCtx, testIdempotencyKey).Return(false, assert.AnError)

		id, err := st.trx.Create(testCtx, nil, testIdempotencyKey)

		assert.Error(t, err)
		assert.Empty(t, id)
	})

	t.Run("idempotency key has been used", func(t *testing.T) {
		st := createTransactionCreatorSuite(ctrl)
		st.keyRepo.EXPECT().Exists(testCtx, testIdempotencyKey).Return(true, nil)

		id, err := st.trx.Create(testCtx, nil, testIdempotencyKey)

		assert.Error(t, err)
		assert.Empty(t, id)
	})

	t.Run("empty transaction is prohibited", func(t *testing.T) {
		st := createTransactionCreatorSuite(ctrl)
		st.keyRepo.EXPECT().Exists(testCtx, testIdempotencyKey).Return(false, nil)

		id, err := st.trx.Create(testCtx, nil, testIdempotencyKey)

		assert.Error(t, err)
		assert.Equal(t, entity.ErrEmptyTransaction(), err)
		assert.Empty(t, id)
	})

	t.Run("sender id is invalid", func(t *testing.T) {
		st := createTransactionCreatorSuite(ctrl)
		trx := createTestTransaction()
		trx.SenderID = uuid.Nil
		st.keyRepo.EXPECT().Exists(testCtx, testIdempotencyKey).Return(false, nil)

		id, err := st.trx.Create(testCtx, trx, testIdempotencyKey)

		assert.Error(t, err)
		assert.Empty(t, id)
	})

	t.Run("receiver id is invalid", func(t *testing.T) {
		st := createTransactionCreatorSuite(ctrl)
		trx := createTestTransaction()
		trx.ReceiverID = uuid.Nil
		st.keyRepo.EXPECT().Exists(testCtx, testIdempotencyKey).Return(false, nil)

		id, err := st.trx.Create(testCtx, trx, testIdempotencyKey)

		assert.Error(t, err)
		assert.Empty(t, id)
	})

	t.Run("amount is invalid", func(t *testing.T) {
		st := createTransactionCreatorSuite(ctrl)
		trx := createTestTransaction()
		trx.Amount = decimal.Zero
		st.keyRepo.EXPECT().Exists(testCtx, testIdempotencyKey).Return(false, nil)

		id, err := st.trx.Create(testCtx, trx, testIdempotencyKey)

		assert.Error(t, err)
		assert.Empty(t, id)
	})

	t.Run("trx repo insert returns error", func(t *testing.T) {
		st := createTransactionCreatorSuite(ctrl)
		trx := createTestTransaction()
		st.keyRepo.EXPECT().Exists(testCtx, testIdempotencyKey).Return(false, nil)

		st.trxRepo.EXPECT().Insert(testCtx, trx).Return(assert.AnError)

		id, err := st.trx.Create(testCtx, trx, testIdempotencyKey)

		assert.Error(t, err)
		assert.Empty(t, id)
	})

	t.Run("success create a transaction", func(t *testing.T) {
		st := createTransactionCreatorSuite(ctrl)
		trx := createTestTransaction()
		st.keyRepo.EXPECT().Exists(testCtx, testIdempotencyKey).Return(false, nil)

		st.trxRepo.EXPECT().Insert(testCtx, trx).Return(nil)

		id, err := st.trx.Create(testCtx, trx, testIdempotencyKey)

		assert.NoError(t, err)
		assert.NotEmpty(t, id)
	})
}

func createTransactionCreatorSuite(ctrl *gomock.Controller) *TransactionCreatorSuite {
	r := mock_service.NewMockCreateTransactionRepository(ctrl)
	i := mock_service.NewMockIdempotencyKeyRepository(ctrl)
	t := service.NewTransactionCreator(r, i)
	return &TransactionCreatorSuite{
		trx:     t,
		trxRepo: r,
		keyRepo: i,
	}
}

func createTestTransaction() *entity.Transaction {
	return &entity.Transaction{
		ID:         uuid.Must(uuid.NewV7()),
		SenderID:   testSenderID,
		ReceiverID: testReceiverID,
		Amount:     testAmount,
	}
}
