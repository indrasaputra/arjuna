package uow_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	mock_uow "github.com/indrasaputra/arjuna/pkg/sdk/test/mock/uow"
	"github.com/indrasaputra/arjuna/pkg/sdk/uow"
)

var (
	testCtx   = context.Background()
	errReturn = errors.New("error")
)

type UnitWorkerSuite struct {
	db   *mock_uow.MockDB
	tx   *mock_uow.MockTx
	unit *uow.UnitWorker
}

func TestNewUnitWorker(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("successfully create an instance of UnitWorker", func(t *testing.T) {
		st := createUnitWorkerSuite(ctrl)
		assert.NotNil(t, st.unit)
	})
}

func TestUnitWorker_Begin(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("begin returns tx", func(t *testing.T) {
		st := createUnitWorkerSuite(ctrl)
		st.db.EXPECT().Begin(testCtx).Return(st.tx, nil)
		res, err := st.unit.Begin(testCtx)

		assert.NoError(t, err)
		assert.NotNil(t, res)
	})
}

func TestUnitWorker_Finish(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("error exists then rollback return errors", func(t *testing.T) {
		st := createUnitWorkerSuite(ctrl)
		st.tx.EXPECT().Rollback(testCtx).Return(errReturn)
		err := st.unit.Finish(testCtx, st.tx, errReturn)

		assert.Error(t, err)
	})

	t.Run("error exists then success rollback", func(t *testing.T) {
		st := createUnitWorkerSuite(ctrl)
		st.tx.EXPECT().Rollback(testCtx).Return(nil)
		err := st.unit.Finish(testCtx, st.tx, errReturn)

		assert.NoError(t, err)
	})

	t.Run("commit returns error", func(t *testing.T) {
		st := createUnitWorkerSuite(ctrl)
		st.tx.EXPECT().Commit(testCtx).Return(errReturn)
		err := st.unit.Finish(testCtx, st.tx, nil)

		assert.Error(t, err)
	})

	t.Run("success commit", func(t *testing.T) {
		st := createUnitWorkerSuite(ctrl)
		st.tx.EXPECT().Commit(testCtx).Return(nil)
		err := st.unit.Finish(testCtx, st.tx, nil)

		assert.NoError(t, err)
	})
}

func createUnitWorkerSuite(ctrl *gomock.Controller) *UnitWorkerSuite {
	db := mock_uow.NewMockDB(ctrl)
	tx := mock_uow.NewMockTx(ctrl)
	u := uow.NewUnitWorker(db)
	return &UnitWorkerSuite{
		unit: u,
		tx:   tx,
		db:   db,
	}
}
