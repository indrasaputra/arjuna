package service_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	sdklog "github.com/indrasaputra/arjuna/pkg/sdk/log"
	mock_uow "github.com/indrasaputra/arjuna/pkg/sdk/test/mock/uow"
	"github.com/indrasaputra/arjuna/service/user/entity"
	"github.com/indrasaputra/arjuna/service/user/internal/app"
	"github.com/indrasaputra/arjuna/service/user/internal/service"
	mock_service "github.com/indrasaputra/arjuna/service/user/test/mock/service"
)

type ctxKey string

var (
	testCtx            = context.Background()
	testCtxTx          = context.WithValue(testCtx, ctxKey("tx"), true)
	testEnv            = "development"
	testIdempotencyKey = "key"
)

type UserRegistrarSuite struct {
	registrar      *service.UserRegistrar
	txManager      *mock_uow.MockTxManager
	userRepo       *mock_service.MockRegisterUserRepository
	userOutboxRepo *mock_service.MockRegisterUserOutboxRepository
	keyRepo        *mock_service.MockIdempotencyKeyRepository
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

	t.Run("validate idempotency key returns error", func(t *testing.T) {
		st := createUserRegistrarSuite(ctrl)
		st.keyRepo.EXPECT().Exists(testCtx, testIdempotencyKey).Return(false, entity.ErrInternal("error"))

		id, err := st.registrar.Register(testCtx, nil, testIdempotencyKey)

		assert.Error(t, err)
		assert.Empty(t, id)
	})

	t.Run("idempotency key has been used", func(t *testing.T) {
		st := createUserRegistrarSuite(ctrl)
		st.keyRepo.EXPECT().Exists(testCtx, testIdempotencyKey).Return(true, nil)

		id, err := st.registrar.Register(testCtx, nil, testIdempotencyKey)

		assert.Error(t, err)
		assert.Empty(t, id)
	})

	t.Run("empty user is prohibited", func(t *testing.T) {
		st := createUserRegistrarSuite(ctrl)
		st.keyRepo.EXPECT().Exists(testCtx, testIdempotencyKey).Return(false, nil)

		id, err := st.registrar.Register(testCtx, nil, testIdempotencyKey)

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
			st.keyRepo.EXPECT().Exists(testCtx, testIdempotencyKey).Return(false, nil)

			id, err := st.registrar.Register(testCtx, user, testIdempotencyKey)

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

			st.keyRepo.EXPECT().Exists(testCtx, testIdempotencyKey).Return(false, nil)

			id, err := st.registrar.Register(testCtx, user, testIdempotencyKey)

			assert.Error(t, err)
			assert.Equal(t, entity.ErrInvalidEmail(), err)
			assert.Empty(t, id)
		}
	})

	t.Run("user repo insert returns error", func(t *testing.T) {
		st := createUserRegistrarSuite(ctrl)
		user := createTestUser()
		errReturn := entity.ErrInternal("")

		st.keyRepo.EXPECT().Exists(testCtx, testIdempotencyKey).Return(false, nil)
		st.userRepo.EXPECT().Insert(testCtxTx, user).Return(errReturn)
		st.txManager.EXPECT().Do(gomock.Any(), gomock.Any()).
			DoAndReturn(func(_ context.Context, fn func(context.Context) error) error {
				assert.Error(t, fn(testCtxTx))
				return errReturn
			})

		id, err := st.registrar.Register(testCtx, user, testIdempotencyKey)

		assert.Error(t, err)
		assert.Empty(t, id)
	})

	t.Run("user outbox repo insert with tx returns error", func(t *testing.T) {
		st := createUserRegistrarSuite(ctrl)
		user := createTestUser()
		errReturn := entity.ErrInternal("")

		st.keyRepo.EXPECT().Exists(testCtx, testIdempotencyKey).Return(false, nil)
		st.userRepo.EXPECT().Insert(testCtxTx, user).Return(nil)
		st.userOutboxRepo.EXPECT().Insert(testCtxTx, gomock.Any()).Return(errReturn)
		st.txManager.EXPECT().Do(gomock.Any(), gomock.Any()).
			DoAndReturn(func(_ context.Context, fn func(context.Context) error) error {
				assert.Error(t, fn(testCtxTx))
				return errReturn
			})

		id, err := st.registrar.Register(testCtx, user, testIdempotencyKey)

		assert.Error(t, err)
		assert.Empty(t, id)
	})

	t.Run("tx manager returns error", func(t *testing.T) {
		st := createUserRegistrarSuite(ctrl)
		user := createTestUser()
		errReturn := entity.ErrInternal("")

		st.keyRepo.EXPECT().Exists(testCtx, testIdempotencyKey).Return(false, nil)
		st.userRepo.EXPECT().Insert(testCtxTx, user).Return(nil)
		st.userOutboxRepo.EXPECT().Insert(testCtxTx, gomock.Any()).Return(nil)
		st.txManager.EXPECT().Do(gomock.Any(), gomock.Any()).
			DoAndReturn(func(_ context.Context, fn func(context.Context) error) error {
				assert.NoError(t, fn(testCtxTx))
				return errReturn
			})

		id, err := st.registrar.Register(testCtx, user, testIdempotencyKey)

		assert.Error(t, err)
		assert.Empty(t, id)
	})

	t.Run("success register user", func(t *testing.T) {
		st := createUserRegistrarSuite(ctrl)
		user := createTestUser()

		st.keyRepo.EXPECT().Exists(testCtx, testIdempotencyKey).Return(false, nil)
		st.userRepo.EXPECT().Insert(testCtxTx, user).Return(nil)
		st.userOutboxRepo.EXPECT().Insert(testCtxTx, gomock.Any()).Return(nil)
		st.txManager.EXPECT().Do(gomock.Any(), gomock.Any()).
			DoAndReturn(func(_ context.Context, fn func(context.Context) error) error {
				assert.NoError(t, fn(testCtxTx))
				return nil
			})

		id, err := st.registrar.Register(testCtx, user, testIdempotencyKey)

		assert.NoError(t, err)
		assert.NotEmpty(t, id)
	})
}

func createUserRegistrarSuite(ctrl *gomock.Controller) *UserRegistrarSuite {
	m := mock_uow.NewMockTxManager(ctrl)
	ur := mock_service.NewMockRegisterUserRepository(ctrl)
	uor := mock_service.NewMockRegisterUserOutboxRepository(ctrl)
	ik := mock_service.NewMockIdempotencyKeyRepository(ctrl)
	r := service.NewUserRegistrar(m, ur, uor, ik)
	return &UserRegistrarSuite{
		registrar:      r,
		txManager:      m,
		userRepo:       ur,
		userOutboxRepo: uor,
		keyRepo:        ik,
	}
}

func createTestUser() *entity.User {
	return &entity.User{
		ID:    uuid.Must(uuid.NewV7()),
		Name:  "First User",
		Email: "first@user.com",
	}
}
