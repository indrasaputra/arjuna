package service_test

import (
	"context"
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	sdklog "github.com/indrasaputra/arjuna/pkg/sdk/log"
	mock_uow "github.com/indrasaputra/arjuna/pkg/sdk/test/mock/uow"
	"github.com/indrasaputra/arjuna/service/transaction/entity"
	"github.com/indrasaputra/arjuna/service/transaction/internal/app"
	"github.com/indrasaputra/arjuna/service/transaction/internal/service"
	mock_service "github.com/indrasaputra/arjuna/service/transaction/test/mock/service"
)

var (
	testCtx              = context.Background()
	testSenderID         = "1"
	testReceiverID       = "2"
	testSenderWalletID   = "3"
	testReceiverWalletID = "4"
	testEnv              = "development"
	testAmount, _        = decimal.NewFromString("10.23")
	testIdempotencyKey   = "key"
	testErrInternal      = entity.ErrInternal("")
)

type TransactionCreatorSuite struct {
	trx           *service.TransactionCreator
	trxRepo       *mock_service.MockCreateTransactionRepository
	trxOutboxRepo *mock_service.MockCreateTransactionOutboxRepository
	keyRepo       *mock_service.MockIdempotencyKeyRepository
	unit          *mock_uow.MockUnitOfWork
	tx            *mock_uow.MockTx
}

func TestNewTransactionCreator(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("successfully create an instance of Transaction", func(t *testing.T) {
		st := createTransactionCreatorSuite(ctrl)
		assert.NotNil(t, st.trx)
	})
}

func TestTransactionCreator_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	app.Logger = sdklog.NewLogger(testEnv)

	t.Run("validate idempotency key returns error", func(t *testing.T) {
		st := createTransactionCreatorSuite(ctrl)
		st.keyRepo.EXPECT().Exists(testCtx, testIdempotencyKey).Return(false, testErrInternal)

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
		trx.SenderID = ""
		st.keyRepo.EXPECT().Exists(testCtx, testIdempotencyKey).Return(false, nil)

		id, err := st.trx.Create(testCtx, trx, testIdempotencyKey)

		assert.Error(t, err)
		assert.Empty(t, id)
	})

	t.Run("receiver id is invalid", func(t *testing.T) {
		st := createTransactionCreatorSuite(ctrl)
		trx := createTestTransaction()
		trx.ReceiverID = ""
		st.keyRepo.EXPECT().Exists(testCtx, testIdempotencyKey).Return(false, nil)

		id, err := st.trx.Create(testCtx, trx, testIdempotencyKey)

		assert.Error(t, err)
		assert.Empty(t, id)
	})

	t.Run("sender wallet id is invalid", func(t *testing.T) {
		st := createTransactionCreatorSuite(ctrl)
		trx := createTestTransaction()
		trx.SenderWalletID = ""
		st.keyRepo.EXPECT().Exists(testCtx, testIdempotencyKey).Return(false, nil)

		id, err := st.trx.Create(testCtx, trx, testIdempotencyKey)

		assert.Error(t, err)
		assert.Empty(t, id)
	})

	t.Run("receiver wallet id is invalid", func(t *testing.T) {
		st := createTransactionCreatorSuite(ctrl)
		trx := createTestTransaction()
		trx.ReceiverWalletID = ""
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

	t.Run("tx begin returns error", func(t *testing.T) {
		st := createTransactionCreatorSuite(ctrl)
		trx := createTestTransaction()
		st.keyRepo.EXPECT().Exists(testCtx, testIdempotencyKey).Return(false, nil)
		st.unit.EXPECT().Begin(testCtx).Return(nil, testErrInternal)

		id, err := st.trx.Create(testCtx, trx, testIdempotencyKey)

		assert.Error(t, err)
		assert.Empty(t, id)
	})

	t.Run("trx repo returns error", func(t *testing.T) {
		st := createTransactionCreatorSuite(ctrl)
		trx := createTestTransaction()
		st.keyRepo.EXPECT().Exists(testCtx, testIdempotencyKey).Return(false, nil)
		st.unit.EXPECT().Begin(testCtx).Return(st.tx, nil)
		st.trxRepo.EXPECT().InsertWithTx(testCtx, st.tx, trx).Return(testErrInternal)
		st.unit.EXPECT().Finish(testCtx, st.tx, testErrInternal).Return(nil)

		id, err := st.trx.Create(testCtx, trx, testIdempotencyKey)

		assert.Error(t, err)
		assert.Empty(t, id)
	})

	t.Run("trx outbox repo returns error", func(t *testing.T) {
		st := createTransactionCreatorSuite(ctrl)
		trx := createTestTransaction()
		st.keyRepo.EXPECT().Exists(testCtx, testIdempotencyKey).Return(false, nil)
		st.unit.EXPECT().Begin(testCtx).Return(st.tx, nil)
		st.trxRepo.EXPECT().InsertWithTx(testCtx, st.tx, trx).Return(nil)
		st.trxOutboxRepo.EXPECT().InsertWithTx(testCtx, st.tx, gomock.Any()).Return(testErrInternal)
		st.unit.EXPECT().Finish(testCtx, st.tx, testErrInternal).Return(nil)

		id, err := st.trx.Create(testCtx, trx, testIdempotencyKey)

		assert.Error(t, err)
		assert.Empty(t, id)
	})

	t.Run("unit of work finish returns error", func(t *testing.T) {
		st := createTransactionCreatorSuite(ctrl)
		trx := createTestTransaction()
		st.keyRepo.EXPECT().Exists(testCtx, testIdempotencyKey).Return(false, nil)
		st.unit.EXPECT().Begin(testCtx).Return(st.tx, nil)
		st.trxRepo.EXPECT().InsertWithTx(testCtx, st.tx, trx).Return(nil)
		st.trxOutboxRepo.EXPECT().InsertWithTx(testCtx, st.tx, gomock.Any()).Return(nil)
		st.unit.EXPECT().Finish(testCtx, st.tx, nil).Return(testErrInternal)

		id, err := st.trx.Create(testCtx, trx, testIdempotencyKey)

		assert.Error(t, err)
		assert.Empty(t, id)
	})

	t.Run("success create transaction error", func(t *testing.T) {
		st := createTransactionCreatorSuite(ctrl)
		trx := createTestTransaction()
		st.keyRepo.EXPECT().Exists(testCtx, testIdempotencyKey).Return(false, nil)
		st.unit.EXPECT().Begin(testCtx).Return(st.tx, nil)
		st.trxRepo.EXPECT().InsertWithTx(testCtx, st.tx, trx).Return(nil)
		st.trxOutboxRepo.EXPECT().InsertWithTx(testCtx, st.tx, gomock.Any()).Return(nil)
		st.unit.EXPECT().Finish(testCtx, st.tx, nil).Return(nil)

		id, err := st.trx.Create(testCtx, trx, testIdempotencyKey)

		assert.NoError(t, err)
		assert.NotEmpty(t, id)
	})
}

func createTestTransaction() *entity.Transaction {
	return &entity.Transaction{
		ID:               "1",
		SenderID:         testSenderID,
		ReceiverID:       testReceiverID,
		SenderWalletID:   testSenderWalletID,
		ReceiverWalletID: testReceiverWalletID,
		Amount:           testAmount,
	}
}

func createTransactionCreatorSuite(ctrl *gomock.Controller) *TransactionCreatorSuite {
	r := mock_service.NewMockCreateTransactionRepository(ctrl)
	o := mock_service.NewMockCreateTransactionOutboxRepository(ctrl)
	i := mock_service.NewMockIdempotencyKeyRepository(ctrl)
	u := mock_uow.NewMockUnitOfWork(ctrl)
	x := mock_uow.NewMockTx(ctrl)
	t := service.NewTransactionCreator(r, o, i, u)
	return &TransactionCreatorSuite{
		trx:           t,
		trxRepo:       r,
		trxOutboxRepo: o,
		keyRepo:       i,
		unit:          u,
		tx:            x,
	}
}
