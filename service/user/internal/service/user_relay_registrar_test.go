package service_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	sdklog "github.com/indrasaputra/arjuna/pkg/sdk/log"
	"github.com/indrasaputra/arjuna/service/user/entity"
	"github.com/indrasaputra/arjuna/service/user/internal/app"
	"github.com/indrasaputra/arjuna/service/user/internal/service"
	mock_service "github.com/indrasaputra/arjuna/service/user/test/mock/service"
)

const (
	limitGetAllReady = uint(1)
)

type UserRelayRegistrarSuite struct {
	relayer        *service.UserRelayRegistrar
	userOutboxRepo *mock_service.MockRelayRegisterUserOutboxRepository
	orchestration  *mock_service.MockRelayRegisterUserOrchestration
}

func TestNewUserRelayRegistrar(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("successfully create an instance of UserRelayRegistrar", func(t *testing.T) {
		st := createUserRelayRegistrarSuite(ctrl)
		assert.NotNil(t, st.relayer)
	})
}

func TestUserRelayRegistrar_Register(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	app.Logger = sdklog.NewLogger(testEnv)

	t.Run("user outbox get all ready record returns error", func(t *testing.T) {
		st := createUserRelayRegistrarSuite(ctrl)
		errReturn := entity.ErrInternal("")

		st.userOutboxRepo.EXPECT().GetAllReady(testCtx, limitGetAllReady).Return(nil, errReturn)

		err := st.relayer.Register(testCtx)

		assert.Error(t, err)
	})

	t.Run("set record processed returns error", func(t *testing.T) {
		st := createUserRelayRegistrarSuite(ctrl)
		rc := createTestUserOutbox()
		records := []*entity.UserOutbox{rc}
		errReturn := entity.ErrInternal("")

		st.userOutboxRepo.EXPECT().GetAllReady(testCtx, limitGetAllReady).Return(records, nil)
		st.userOutboxRepo.EXPECT().SetProcessed(testCtx, rc.ID).Return(errReturn)

		err := st.relayer.Register(testCtx)

		assert.NoError(t, err)
	})

	t.Run("set record processed returns error", func(t *testing.T) {
		st := createUserRelayRegistrarSuite(ctrl)
		rc := createTestUserOutbox()
		records := []*entity.UserOutbox{rc}
		errReturn := entity.ErrInternal("")

		st.userOutboxRepo.EXPECT().GetAllReady(testCtx, limitGetAllReady).Return(records, nil)
		st.userOutboxRepo.EXPECT().SetProcessed(testCtx, rc.ID).Return(errReturn)

		err := st.relayer.Register(testCtx)

		assert.NoError(t, err)
	})

	t.Run("enqueue to orchestrator returns error", func(t *testing.T) {
		st := createUserRelayRegistrarSuite(ctrl)
		rc := createTestUserOutbox()
		records := []*entity.UserOutbox{rc}
		errReturn := entity.ErrInternal("")

		st.userOutboxRepo.EXPECT().GetAllReady(testCtx, limitGetAllReady).Return(records, nil)
		st.userOutboxRepo.EXPECT().SetProcessed(testCtx, rc.ID).Return(nil)
		st.orchestration.EXPECT().RegisterUser(testCtx, gomock.Any()).Return(nil, errReturn)
		st.userOutboxRepo.EXPECT().SetFailed(testCtx, rc.ID).Return(nil)

		err := st.relayer.Register(testCtx)

		assert.NoError(t, err)
	})

	t.Run("set failed returns error", func(t *testing.T) {
		st := createUserRelayRegistrarSuite(ctrl)
		rc := createTestUserOutbox()
		records := []*entity.UserOutbox{rc}
		errReturn := entity.ErrInternal("")

		st.userOutboxRepo.EXPECT().GetAllReady(testCtx, limitGetAllReady).Return(records, nil)
		st.userOutboxRepo.EXPECT().SetProcessed(testCtx, rc.ID).Return(nil)
		st.orchestration.EXPECT().RegisterUser(testCtx, gomock.Any()).Return(nil, errReturn)
		st.userOutboxRepo.EXPECT().SetFailed(testCtx, rc.ID).Return(errReturn)

		err := st.relayer.Register(testCtx)

		assert.NoError(t, err)
	})

	t.Run("set delivered returns error", func(t *testing.T) {
		st := createUserRelayRegistrarSuite(ctrl)
		rc := createTestUserOutbox()
		records := []*entity.UserOutbox{rc}
		errReturn := entity.ErrInternal("")

		st.userOutboxRepo.EXPECT().GetAllReady(testCtx, limitGetAllReady).Return(records, nil)
		st.userOutboxRepo.EXPECT().SetProcessed(testCtx, rc.ID).Return(nil)
		st.orchestration.EXPECT().RegisterUser(testCtx, gomock.Any()).Return(nil, nil)
		st.userOutboxRepo.EXPECT().SetDelivered(testCtx, rc.ID).Return(errReturn)

		err := st.relayer.Register(testCtx)

		assert.NoError(t, err)
	})

	t.Run("set delivered success", func(t *testing.T) {
		st := createUserRelayRegistrarSuite(ctrl)
		rc := createTestUserOutbox()
		records := []*entity.UserOutbox{rc}

		st.userOutboxRepo.EXPECT().GetAllReady(testCtx, limitGetAllReady).Return(records, nil)
		st.userOutboxRepo.EXPECT().SetProcessed(testCtx, rc.ID).Return(nil)
		st.orchestration.EXPECT().RegisterUser(testCtx, gomock.Any()).Return(nil, nil)
		st.userOutboxRepo.EXPECT().SetDelivered(testCtx, rc.ID).Return(nil)

		err := st.relayer.Register(testCtx)

		assert.NoError(t, err)
	})
}

func createUserRelayRegistrarSuite(ctrl *gomock.Controller) *UserRelayRegistrarSuite {
	u := mock_service.NewMockRelayRegisterUserOutboxRepository(ctrl)
	o := mock_service.NewMockRelayRegisterUserOrchestration(ctrl)
	r := service.NewUserRelayRegistrar(u, o)
	return &UserRelayRegistrarSuite{
		relayer:        r,
		userOutboxRepo: u,
		orchestration:  o,
	}
}

func createTestUserOutbox() *entity.UserOutbox {
	user := createTestUser()
	return &entity.UserOutbox{
		ID:      "1",
		Status:  entity.UserOutboxStatusReady,
		Payload: user,
	}
}
