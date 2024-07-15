package service_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	sdklog "github.com/indrasaputra/arjuna/pkg/sdk/log"
	mock_uow "github.com/indrasaputra/arjuna/pkg/sdk/test/mock/uow"
	"github.com/indrasaputra/arjuna/service/user/entity"
	"github.com/indrasaputra/arjuna/service/user/internal/app"
	"github.com/indrasaputra/arjuna/service/user/internal/service"
	mock_service "github.com/indrasaputra/arjuna/service/user/test/mock/service"
)

var (
	testCtx = context.Background()
	testEnv = "development"
)

type UserRegistrarSuite struct {
	registrar      *service.UserRegistrar
	userRepo       *mock_service.MockRegisterUserRepository
	userOutboxRepo *mock_service.MockRegisterUserOutboxRepository
	unit           *mock_uow.MockUnitOfWork
	tx             *mock_uow.MockTx
}

func TestNewUserRegistrar(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("successfully create an instance of UserRegistrar", func(t *testing.T) {
		st := createUserRegistrarSuite(ctrl)
		assert.NotNil(t, st.registrar)
	})
}

func TestUserRegistrar_Register(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	app.Logger = sdklog.NewLogger(testEnv)

	t.Run("empty user is prohibited", func(t *testing.T) {
		st := createUserRegistrarSuite(ctrl)

		id, err := st.registrar.Register(testCtx, nil)

		assert.Error(t, err)
		assert.Equal(t, entity.ErrEmptyUser(), err)
		assert.Empty(t, id)
	})

	t.Run("name contains character other than alphabet", func(t *testing.T) {
		st := createUserRegistrarSuite(ctrl)
		names := []string{
			"123",
			"First Us3r",
			"First User !!!",
			"!@#$%^&*()",
		}

		for _, name := range names {
			user := &entity.User{Name: name}

			id, err := st.registrar.Register(testCtx, user)

			assert.Error(t, err)
			assert.Equal(t, entity.ErrInvalidName(), err)
			assert.Empty(t, id)
		}
	})

	t.Run("email is invalid", func(t *testing.T) {
		st := createUserRegistrarSuite(ctrl)
		emails := []string{
			"@domain",
			"@domain.com",
			"domain.com",
		}

		for _, email := range emails {
			user := createTestUser()
			user.Email = email

			id, err := st.registrar.Register(testCtx, user)

			assert.Error(t, err)
			assert.Equal(t, entity.ErrInvalidEmail(), err)
			assert.Empty(t, id)
		}
	})

	t.Run("unit of work begin returns error", func(t *testing.T) {
		st := createUserRegistrarSuite(ctrl)
		user := createTestUser()
		errReturn := entity.ErrInternal("")

		st.unit.EXPECT().Begin(testCtx).Return(nil, errReturn)

		id, err := st.registrar.Register(testCtx, user)

		assert.Error(t, err)
		assert.Empty(t, id)
	})

	t.Run("user repo insert with tx returns error", func(t *testing.T) {
		st := createUserRegistrarSuite(ctrl)
		user := createTestUser()
		errReturn := entity.ErrInternal("")

		st.unit.EXPECT().Begin(testCtx).Return(st.tx, nil)
		st.userRepo.EXPECT().InsertWithTx(testCtx, st.tx, user).Return(errReturn)
		st.unit.EXPECT().Finish(testCtx, st.tx, errReturn).Return(nil)

		id, err := st.registrar.Register(testCtx, user)

		assert.Error(t, err)
		assert.Empty(t, id)
	})

	t.Run("user outbox repo insert with tx returns error", func(t *testing.T) {
		st := createUserRegistrarSuite(ctrl)
		user := createTestUser()
		errReturn := entity.ErrInternal("")

		st.unit.EXPECT().Begin(testCtx).Return(st.tx, nil)
		st.userRepo.EXPECT().InsertWithTx(testCtx, st.tx, user).Return(nil)
		st.userOutboxRepo.EXPECT().InsertWithTx(testCtx, st.tx, gomock.Any()).Return(errReturn)
		st.unit.EXPECT().Finish(testCtx, st.tx, errReturn).Return(nil)

		id, err := st.registrar.Register(testCtx, user)

		assert.Error(t, err)
		assert.Empty(t, id)
	})

	t.Run("unit of work finish returns error", func(t *testing.T) {
		st := createUserRegistrarSuite(ctrl)
		user := createTestUser()
		errReturn := entity.ErrInternal("")

		st.unit.EXPECT().Begin(testCtx).Return(st.tx, nil)
		st.userRepo.EXPECT().InsertWithTx(testCtx, st.tx, user).Return(nil)
		st.userOutboxRepo.EXPECT().InsertWithTx(testCtx, st.tx, gomock.Any()).Return(nil)
		st.unit.EXPECT().Finish(testCtx, st.tx, nil).Return(errReturn)

		id, err := st.registrar.Register(testCtx, user)

		assert.Error(t, err)
		assert.Empty(t, id)
	})

	t.Run("success register user", func(t *testing.T) {
		st := createUserRegistrarSuite(ctrl)
		user := createTestUser()

		st.unit.EXPECT().Begin(testCtx).Return(st.tx, nil)
		st.userRepo.EXPECT().InsertWithTx(testCtx, st.tx, user).Return(nil)
		st.userOutboxRepo.EXPECT().InsertWithTx(testCtx, st.tx, gomock.Any()).Return(nil)
		st.unit.EXPECT().Finish(testCtx, st.tx, nil).Return(nil)

		id, err := st.registrar.Register(testCtx, user)

		assert.NoError(t, err)
		assert.NotEmpty(t, id)
	})
}

func createUserRegistrarSuite(ctrl *gomock.Controller) *UserRegistrarSuite {
	ur := mock_service.NewMockRegisterUserRepository(ctrl)
	uor := mock_service.NewMockRegisterUserOutboxRepository(ctrl)
	u := mock_uow.NewMockUnitOfWork(ctrl)
	tx := mock_uow.NewMockTx(ctrl)
	r := service.NewUserRegistrar(ur, uor, u)
	return &UserRegistrarSuite{
		registrar:      r,
		userRepo:       ur,
		userOutboxRepo: uor,
		unit:           u,
		tx:             tx,
	}
}

func createTestUser() *entity.User {
	return &entity.User{
		ID:    "1",
		Name:  "First User",
		Email: "first@user.com",
	}
}
