package postgres_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/pashagolub/pgxmock/v4"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	sdklog "github.com/indrasaputra/arjuna/pkg/sdk/log"
	mock_uow "github.com/indrasaputra/arjuna/pkg/sdk/test/mock/uow"
	"github.com/indrasaputra/arjuna/pkg/sdk/uow"
	"github.com/indrasaputra/arjuna/service/wallet/entity"
	"github.com/indrasaputra/arjuna/service/wallet/internal/app"
	"github.com/indrasaputra/arjuna/service/wallet/internal/repository/db"
	"github.com/indrasaputra/arjuna/service/wallet/internal/repository/postgres"
)

var (
	testCtx = context.Background()
	testEnv = "development"
)

type WalletSuite struct {
	wallet *postgres.Wallet
	db     pgxmock.PgxPoolIface
	getter *mock_uow.MockTxGetter
}

func TestNewWallet(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("successfully create an instance of Wallet", func(t *testing.T) {
		st := createWalletSuite(t, ctrl)
		assert.NotNil(t, st.wallet)
	})
}

func TestWallet_Insert(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	app.Logger = sdklog.NewLogger(testEnv)
	query := `INSERT INTO wallets \(id, user_id, balance, created_at, updated_at, created_by, updated_by\)
				VALUES \(\$1, \$2, \$3, \$4, \$5, \$6, \$7\)`

	t.Run("nil wallets is prohibited", func(t *testing.T) {
		st := createWalletSuite(t, ctrl)

		err := st.wallet.Insert(testCtx, nil)

		assert.Error(t, err)
		assert.Equal(t, entity.ErrEmptyWallet(), err)
	})

	t.Run("insert duplicate wallets", func(t *testing.T) {
		wallet := createTestWallet()
		st := createWalletSuite(t, ctrl)
		st.getter.EXPECT().DefaultTrOrDB(testCtx, st.db).Return(st.db)
		st.db.ExpectExec(query).
			WithArgs(wallet.ID, wallet.UserID, wallet.Balance, wallet.CreatedAt, wallet.UpdatedAt, wallet.CreatedBy, wallet.UpdatedBy).
			WillReturnError(&pgconn.PgError{Code: "23505"})

		err := st.wallet.Insert(testCtx, wallet)

		assert.Error(t, err)
		assert.Equal(t, entity.ErrAlreadyExists(), err)
	})

	t.Run("insert returns error", func(t *testing.T) {
		wallet := createTestWallet()
		st := createWalletSuite(t, ctrl)
		st.getter.EXPECT().DefaultTrOrDB(testCtx, st.db).Return(st.db)
		st.db.ExpectExec(query).
			WithArgs(wallet.ID, wallet.UserID, wallet.Balance, wallet.CreatedAt, wallet.UpdatedAt, wallet.CreatedBy, wallet.UpdatedBy).
			WillReturnError(assert.AnError)

		err := st.wallet.Insert(testCtx, wallet)

		assert.Error(t, err)
	})

	t.Run("success insert wallets", func(t *testing.T) {
		wallet := createTestWallet()
		st := createWalletSuite(t, ctrl)
		st.getter.EXPECT().DefaultTrOrDB(testCtx, st.db).Return(st.db)
		st.db.ExpectExec(query).
			WithArgs(wallet.ID, wallet.UserID, wallet.Balance, wallet.CreatedAt, wallet.UpdatedAt, wallet.CreatedBy, wallet.UpdatedBy).
			WillReturnResult(pgxmock.NewResult("INSERT", 1))

		err := st.wallet.Insert(testCtx, wallet)

		assert.NoError(t, err)
	})
}

func TestWallet_AddWalletBalance(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	app.Logger = sdklog.NewLogger(testEnv)

	query := `UPDATE wallets SET balance = balance \+ \$2 WHERE id = \$1`

	t.Run("add account balance returns internal error", func(t *testing.T) {
		st := createWalletSuite(t, ctrl)
		id := uuid.Must(uuid.NewV7())
		amount, _ := decimal.NewFromString("4.56")
		st.getter.EXPECT().DefaultTrOrDB(testCtx, st.db).Return(st.db)
		st.db.ExpectExec(query).
			WithArgs(id, amount).
			WillReturnError(assert.AnError)

		err := st.wallet.AddWalletBalance(testCtx, id, amount)

		assert.Error(t, err)
	})

	t.Run("add account balance returns success", func(t *testing.T) {
		st := createWalletSuite(t, ctrl)
		id := uuid.Must(uuid.NewV7())
		amount, _ := decimal.NewFromString("4.56")
		st.getter.EXPECT().DefaultTrOrDB(testCtx, st.db).Return(st.db)
		st.db.ExpectExec(query).
			WithArgs(id, amount).
			WillReturnResult(pgxmock.NewResult("UPDATE", 1))

		err := st.wallet.AddWalletBalance(testCtx, id, amount)

		assert.NoError(t, err)
	})
}

func TestWallet_GetUserWalletForUpdate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	app.Logger = sdklog.NewLogger(testEnv)

	query := `SELECT id, user_id, balance, created_at, updated_at, deleted_at, created_by, updated_by, deleted_by FROM wallets WHERE id = \$1 AND user_id = \$2 LIMIT 1 FOR NO KEY UPDATE`

	t.Run("wallet not found", func(t *testing.T) {
		st := createWalletSuite(t, ctrl)
		id := uuid.Must(uuid.NewV7())
		userID := uuid.Must(uuid.NewV7())
		st.getter.EXPECT().DefaultTrOrDB(testCtx, st.db).Return(st.db)
		st.db.ExpectQuery(query).WithArgs(id, userID).WillReturnError(pgx.ErrNoRows)

		res, err := st.wallet.GetUserWalletForUpdate(testCtx, id, userID)

		assert.Error(t, err)
		assert.Equal(t, entity.ErrEmptyWallet(), err)
		assert.Nil(t, res)
	})

	t.Run("select returns error", func(t *testing.T) {
		st := createWalletSuite(t, ctrl)
		id := uuid.Must(uuid.NewV7())
		userID := uuid.Must(uuid.NewV7())
		st.getter.EXPECT().DefaultTrOrDB(testCtx, st.db).Return(st.db)
		st.db.ExpectQuery(query).WithArgs(id, userID).WillReturnError(assert.AnError)
		// st.db.ExpectQuery(query).WithArgs(id, userID).WillReturnRows(pgxmock.
		// 	NewRows([]string{"id", "user_id", "balance", "created_at", "updated_at", "deleted_at", "created_by", "updated_by", "deleted_by"}).
		// 	AddRow(wallet.ID, wallet.UserID, wallet.Balance, wallet.CreatedAt, wallet.UpdatedAt, wallet.DeletedAt, wallet.CreatedBy, wallet.UpdatedBy, wallet.DeletedBy))

		res, err := st.wallet.GetUserWalletForUpdate(testCtx, id, userID)

		assert.Error(t, err)
		assert.Nil(t, res)
	})

	t.Run("select returns error", func(t *testing.T) {
		wallet := createTestWallet()
		st := createWalletSuite(t, ctrl)
		id := uuid.Must(uuid.NewV7())
		userID := uuid.Must(uuid.NewV7())
		st.getter.EXPECT().DefaultTrOrDB(testCtx, st.db).Return(st.db)
		st.db.ExpectQuery(query).WithArgs(id, userID).WillReturnRows(pgxmock.
			NewRows([]string{"id", "user_id", "balance", "created_at", "updated_at", "deleted_at", "created_by", "updated_by", "deleted_by"}).
			AddRow(wallet.ID, wallet.UserID, wallet.Balance, wallet.CreatedAt, wallet.UpdatedAt, wallet.DeletedAt, wallet.CreatedBy, wallet.UpdatedBy, wallet.DeletedBy))

		res, err := st.wallet.GetUserWalletForUpdate(testCtx, id, userID)

		assert.NoError(t, err)
		assert.NotNil(t, res)
	})
}

func createTestWallet() *entity.Wallet {
	b, _ := decimal.NewFromString("10.23")
	return &entity.Wallet{
		ID:      uuid.Must(uuid.NewV7()),
		UserID:  uuid.Must(uuid.NewV7()),
		Balance: b,
	}
}

func createWalletSuite(t *testing.T, ctrl *gomock.Controller) *WalletSuite {
	pool, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("error opening a stub database connection: %v\n", err)
	}
	defer pool.Close()
	g := mock_uow.NewMockTxGetter(ctrl)
	tx := uow.NewTxDB(pool, g)
	q := db.New(tx)
	w := postgres.NewWallet(q)
	return &WalletSuite{
		wallet: w,
		db:     pool,
		getter: g,
	}
}
