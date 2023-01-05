package uow_test

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	mock_uow "github.com/indrasaputra/arjuna/pkg/sdk/test/mock/uow"
	"github.com/indrasaputra/arjuna/pkg/sdk/uow"
)

var (
	testCtx   = context.Background()
	errReturn = errors.New("error")
)

type UnitWorkerExecutor struct {
	db   *mock_uow.MockDB
	tx   *mock_uow.MockTx
	unit *uow.UnitWorker
}

func TestNewUnitWorker(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("successfully create an instance of UnitWorker", func(t *testing.T) {
		exec := createUnitWorkerExecutor(ctrl)
		assert.NotNil(t, exec.unit)
	})
}

func TestUnitWorker_Begin(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("begin returns tx", func(t *testing.T) {
		exec := createUnitWorkerExecutor(ctrl)
		exec.db.EXPECT().Begin(testCtx).Return(exec.tx, nil)
		res, err := exec.unit.Begin(testCtx)

		assert.NoError(t, err)
		assert.NotNil(t, res)
	})
}

func TestUnitWorker_Finish(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("error exists then rollback return errors", func(t *testing.T) {
		exec := createUnitWorkerExecutor(ctrl)
		exec.tx.EXPECT().Rollback(testCtx).Return(errReturn)
		err := exec.unit.Finish(testCtx, exec.tx, errReturn)

		assert.Error(t, err)
	})

	t.Run("error exists then success rollback", func(t *testing.T) {
		exec := createUnitWorkerExecutor(ctrl)
		exec.tx.EXPECT().Rollback(testCtx).Return(nil)
		err := exec.unit.Finish(testCtx, exec.tx, errReturn)

		assert.NoError(t, err)
	})

	t.Run("commit returns error", func(t *testing.T) {
		exec := createUnitWorkerExecutor(ctrl)
		exec.tx.EXPECT().Commit(testCtx).Return(errReturn)
		err := exec.unit.Finish(testCtx, exec.tx, nil)

		assert.Error(t, err)
	})

	t.Run("success commit", func(t *testing.T) {
		exec := createUnitWorkerExecutor(ctrl)
		exec.tx.EXPECT().Commit(testCtx).Return(nil)
		err := exec.unit.Finish(testCtx, exec.tx, nil)

		assert.NoError(t, err)
	})
}

func createUnitWorkerExecutor(ctrl *gomock.Controller) *UnitWorkerExecutor {
	db := mock_uow.NewMockDB(ctrl)
	tx := mock_uow.NewMockTx(ctrl)
	u := uow.NewUnitWorker(db)
	return &UnitWorkerExecutor{
		unit: u,
		tx:   tx,
		db:   db,
	}
}
