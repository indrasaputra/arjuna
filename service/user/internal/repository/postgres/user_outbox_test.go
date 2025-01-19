package postgres_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/pashagolub/pgxmock/v4"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	sdkpostgres "github.com/indrasaputra/arjuna/pkg/sdk/database/postgres"
	mock_uow "github.com/indrasaputra/arjuna/pkg/sdk/test/mock/uow"
	"github.com/indrasaputra/arjuna/service/user/entity"
	"github.com/indrasaputra/arjuna/service/user/internal/repository/db"
	"github.com/indrasaputra/arjuna/service/user/internal/repository/postgres"
)

const (
	queryUpdateRecordStatus = `UPDATE users_outbox SET status = \$1, updated_at = NOW\(\) WHERE id = \$2`
)

type UserOutboxSuite struct {
	outbox *postgres.UserOutbox
	db     pgxmock.PgxPoolIface
	getter *mock_uow.MockTxGetter
}

func TestNewUserOutbox(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("successfully create an instance of User", func(t *testing.T) {
		st := createUserOutboxSuite(t, ctrl)
		assert.NotNil(t, st.outbox)
	})
}

func TestUserOutbox_Insert(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	query := `INSERT INTO
		users_outbox \(id, status, payload, created_at, updated_at, created_by, updated_by\)
		VALUES \(\$1, \$2, \$3, \$4, \$5, \$6, \$7\)`

	t.Run("nil outbox is prohibited", func(t *testing.T) {
		st := createUserOutboxSuite(t, ctrl)

		err := st.outbox.Insert(testCtx, nil)

		assert.Error(t, err)
		assert.Equal(t, entity.ErrEmptyUser(), err)
	})

	t.Run("nil payload is prohibited", func(t *testing.T) {
		st := createUserOutboxSuite(t, ctrl)
		out := createTestUserOutbox()
		out.Payload = nil

		err := st.outbox.Insert(testCtx, out)

		assert.Error(t, err)
		assert.Equal(t, entity.ErrEmptyUser(), err)
	})

	t.Run("insert duplicate outbox", func(t *testing.T) {
		st := createUserOutboxSuite(t, ctrl)
		out := createTestUserOutbox()
		st.getter.EXPECT().DefaultTrOrDB(testCtx, st.db).Return(st.db)
		st.db.ExpectExec(query).WithArgs(out.ID, db.UserOutboxStatus(out.Status), out.Payload, out.CreatedAt, out.UpdatedAt, out.CreatedBy, out.UpdatedBy).WillReturnError(&pgconn.PgError{Code: "23505"})

		err := st.outbox.Insert(testCtx, out)

		assert.Error(t, err)
		assert.Equal(t, entity.ErrAlreadyExists(), err)
	})

	t.Run("insert returns error", func(t *testing.T) {
		st := createUserOutboxSuite(t, ctrl)
		out := createTestUserOutbox()
		st.getter.EXPECT().DefaultTrOrDB(testCtx, st.db).Return(st.db)
		st.db.ExpectExec(query).WithArgs(out.ID, db.UserOutboxStatus(out.Status), out.Payload, out.CreatedAt, out.UpdatedAt, out.CreatedBy, out.UpdatedBy).WillReturnError(assert.AnError)

		err := st.outbox.Insert(testCtx, out)

		assert.Error(t, err)
	})

	t.Run("success insert", func(t *testing.T) {
		st := createUserOutboxSuite(t, ctrl)
		out := createTestUserOutbox()
		st.getter.EXPECT().DefaultTrOrDB(testCtx, st.db).Return(st.db)
		st.db.ExpectExec(query).WithArgs(out.ID, db.UserOutboxStatus(out.Status), out.Payload, out.CreatedAt, out.UpdatedAt, out.CreatedBy, out.UpdatedBy).WillReturnResult(pgxmock.NewResult("INSERT", 1))

		err := st.outbox.Insert(testCtx, out)

		assert.NoError(t, err)
	})
}

func TestUserOutbox_GetAllReady(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	query := `SELECT id, payload, status, created_at, updated_at, deleted_at, created_by, updated_by, deleted_by FROM users_outbox WHERE status = \$1 ORDER BY created_at ASC LIMIT \$2`
	limit := uint(10)

	t.Run("get all returns error", func(t *testing.T) {
		st := createUserOutboxSuite(t, ctrl)
		st.getter.EXPECT().DefaultTrOrDB(testCtx, st.db).Return(st.db)
		st.db.ExpectQuery(query).WithArgs(db.UserOutboxStatusREADY, int32(limit)).WillReturnError(assert.AnError)

		res, err := st.outbox.GetAllReady(testCtx, limit)

		assert.Error(t, err)
		assert.Empty(t, res)
	})

	t.Run("success get all", func(t *testing.T) {
		user := createTestUser()
		st := createUserOutboxSuite(t, ctrl)
		st.getter.EXPECT().DefaultTrOrDB(testCtx, st.db).Return(st.db)
		st.db.ExpectQuery(query).WithArgs(db.UserOutboxStatusREADY, int32(limit)).WillReturnRows(pgxmock.
			NewRows([]string{"id", "payload", "status", "created_at", "updated_at", "deleted_at", "created_by", "updated_by", "deleted_by"}).
			AddRow(uuid.Must(uuid.NewV7()), user, db.UserOutboxStatusREADY, user.CreatedAt, user.UpdatedAt, user.DeletedAt, user.CreatedBy, user.UpdatedBy, user.DeletedBy))

		res, err := st.outbox.GetAllReady(testCtx, limit)

		assert.NoError(t, err)
		assert.Len(t, res, 1)
	})
}

func TestUserOutbox_SetProcessed(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("set processed returns error", func(t *testing.T) {
		out := createTestUserOutbox()
		st := createUserOutboxSuite(t, ctrl)
		st.getter.EXPECT().DefaultTrOrDB(testCtx, st.db).Return(st.db)
		st.db.ExpectExec(queryUpdateRecordStatus).WithArgs(db.UserOutboxStatusPROCESSED, out.ID).WillReturnError(assert.AnError)

		err := st.outbox.SetProcessed(testCtx, out.ID)

		assert.Error(t, err)
	})

	t.Run("set processed success", func(t *testing.T) {
		out := createTestUser()
		st := createUserOutboxSuite(t, ctrl)
		st.getter.EXPECT().DefaultTrOrDB(testCtx, st.db).Return(st.db)
		st.db.ExpectExec(queryUpdateRecordStatus).WithArgs(db.UserOutboxStatusPROCESSED, out.ID).WillReturnResult(pgxmock.NewResult("UPDATE", 1))

		err := st.outbox.SetProcessed(testCtx, out.ID)

		assert.NoError(t, err)
	})
}

func TestUserOutbox_SetDelivered(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("set delivered returns error", func(t *testing.T) {
		out := createTestUserOutbox()
		st := createUserOutboxSuite(t, ctrl)
		st.getter.EXPECT().DefaultTrOrDB(testCtx, st.db).Return(st.db)
		st.db.ExpectExec(queryUpdateRecordStatus).WithArgs(db.UserOutboxStatusDELIVERED, out.ID).WillReturnError(assert.AnError)

		err := st.outbox.SetDelivered(testCtx, out.ID)

		assert.Error(t, err)
	})

	t.Run("set delivered success", func(t *testing.T) {
		out := createTestUser()
		st := createUserOutboxSuite(t, ctrl)
		st.getter.EXPECT().DefaultTrOrDB(testCtx, st.db).Return(st.db)
		st.db.ExpectExec(queryUpdateRecordStatus).WithArgs(db.UserOutboxStatusDELIVERED, out.ID).WillReturnResult(pgxmock.NewResult("UPDATE", 1))

		err := st.outbox.SetDelivered(testCtx, out.ID)

		assert.NoError(t, err)
	})
}

func TestUserOutbox_SetFailed(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("set failed returns error", func(t *testing.T) {
		out := createTestUserOutbox()
		st := createUserOutboxSuite(t, ctrl)
		st.getter.EXPECT().DefaultTrOrDB(testCtx, st.db).Return(st.db)
		st.db.ExpectExec(queryUpdateRecordStatus).WithArgs(db.UserOutboxStatusFAILED, out.ID).WillReturnError(assert.AnError)

		err := st.outbox.SetFailed(testCtx, out.ID)

		assert.Error(t, err)
	})

	t.Run("set failed success", func(t *testing.T) {
		out := createTestUser()
		st := createUserOutboxSuite(t, ctrl)
		st.getter.EXPECT().DefaultTrOrDB(testCtx, st.db).Return(st.db)
		st.db.ExpectExec(queryUpdateRecordStatus).WithArgs(db.UserOutboxStatusFAILED, out.ID).WillReturnResult(pgxmock.NewResult("UPDATE", 1))

		err := st.outbox.SetFailed(testCtx, out.ID)

		assert.NoError(t, err)
	})
}

func TestUserOutbox_SetRecordStatus(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("set status returns error", func(t *testing.T) {
		out := createTestUserOutbox()
		status := entity.UserOutboxStatusProcessed
		st := createUserOutboxSuite(t, ctrl)
		st.getter.EXPECT().DefaultTrOrDB(testCtx, st.db).Return(st.db)
		st.db.ExpectExec(queryUpdateRecordStatus).WithArgs(db.UserOutboxStatus(status), out.ID).WillReturnError(assert.AnError)

		err := st.outbox.SetRecordStatus(testCtx, out.ID, status)

		assert.Error(t, err)
	})

	t.Run("set status success", func(t *testing.T) {
		out := createTestUser()
		status := entity.UserOutboxStatusProcessed
		st := createUserOutboxSuite(t, ctrl)
		st.getter.EXPECT().DefaultTrOrDB(testCtx, st.db).Return(st.db)
		st.db.ExpectExec(queryUpdateRecordStatus).WithArgs(db.UserOutboxStatus(status), out.ID).WillReturnResult(pgxmock.NewResult("UPDATE", 1))

		err := st.outbox.SetRecordStatus(testCtx, out.ID, status)

		assert.NoError(t, err)
	})
}

func createTestUserOutbox() *entity.UserOutbox {
	user := createTestUser()
	return &entity.UserOutbox{
		ID:      uuid.Must(uuid.NewV7()),
		Status:  entity.UserOutboxStatusReady,
		Payload: user,
	}
}

func createUserOutboxSuite(t *testing.T, ctrl *gomock.Controller) *UserOutboxSuite {
	pool, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("error opening a stub database connection: %v\n", err)
	}
	g := mock_uow.NewMockTxGetter(ctrl)
	tx := sdkpostgres.NewTxDB(pool, g)
	q := db.New(tx)
	o := postgres.NewUserOutbox(q)
	return &UserOutboxSuite{
		outbox: o,
		db:     pool,
		getter: g,
	}
}
