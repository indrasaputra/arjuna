package handler_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc/metadata"

	sdklog "github.com/indrasaputra/arjuna/pkg/sdk/log"
	apiv1 "github.com/indrasaputra/arjuna/proto/api/v1"
	"github.com/indrasaputra/arjuna/service/user/entity"
	"github.com/indrasaputra/arjuna/service/user/internal/app"
	"github.com/indrasaputra/arjuna/service/user/internal/grpc/handler"
	mock_service "github.com/indrasaputra/arjuna/service/user/test/mock/service"
)

const (
	testIdempotencyKey = "key"
)

var (
	testCtxWithValidKey   = metadata.NewIncomingContext(testCtx, metadata.Pairs("X-Idempotency-Key", testIdempotencyKey))
	testCtxWithInvalidKey = metadata.NewIncomingContext(testCtx, metadata.Pairs("another-key", ""))
)

type UserCommandSuite struct {
	handler   *handler.UserCommand
	registrar *mock_service.MockRegisterUser
}

func TestNewUserCommand(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("successful create an instance of UserCommand", func(t *testing.T) {
		st := createUserCommandSuite(ctrl)
		assert.NotNil(t, st.handler)
	})
}

func TestUserCommand_RegisterUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	app.Logger = sdklog.NewLogger(testEnv)

	t.Run("idempotency key is missing", func(t *testing.T) {
		st := createUserCommandSuite(ctrl)

		res, err := st.handler.RegisterUser(testCtxWithInvalidKey, nil)

		assert.Error(t, err)
		assert.Equal(t, entity.ErrMissingIdempotencyKey(), err)
		assert.Nil(t, res)
	})

	t.Run("nil request is prohibited", func(t *testing.T) {
		st := createUserCommandSuite(ctrl)

		res, err := st.handler.RegisterUser(testCtxWithValidKey, nil)

		assert.Error(t, err)
		assert.Equal(t, entity.ErrEmptyUser(), err)
		assert.Nil(t, res)
	})

	t.Run("empty user is prohibited", func(t *testing.T) {
		st := createUserCommandSuite(ctrl)

		res, err := st.handler.RegisterUser(testCtxWithValidKey, &apiv1.RegisterUserRequest{})

		assert.Error(t, err)
		assert.Equal(t, entity.ErrEmptyUser(), err)
		assert.Nil(t, res)
	})

	t.Run("user service returns error", func(t *testing.T) {
		st := createUserCommandSuite(ctrl)
		request := &apiv1.RegisterUserRequest{
			User: &apiv1.User{
				Name:     "First User",
				Email:    "first@user.com",
				Password: "BestPlayer",
			},
		}

		errors := []error{
			entity.ErrEmptyUser(),
			entity.ErrAlreadyExists(),
			entity.ErrInvalidName(),
			entity.ErrInvalidEmail(),
			entity.ErrInternal("error"),
		}
		for _, errRet := range errors {
			st.registrar.EXPECT().Register(testCtxWithValidKey, gomock.Any(), testIdempotencyKey).Return("", errRet)

			res, err := st.handler.RegisterUser(testCtxWithValidKey, request)

			assert.Error(t, err)
			assert.Equal(t, errRet, err)
			assert.Nil(t, res)
		}
	})

	t.Run("success register user", func(t *testing.T) {
		st := createUserCommandSuite(ctrl)
		st.registrar.EXPECT().Register(testCtxWithValidKey, gomock.Any(), testIdempotencyKey).Return("id", nil)
		request := &apiv1.RegisterUserRequest{
			User: &apiv1.User{
				Name:     "First User",
				Email:    "first@user.com",
				Password: "BestPlayer",
			},
		}

		res, err := st.handler.RegisterUser(testCtxWithValidKey, request)

		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, "id", res.Data.GetId())
	})
}

func createUserCommandSuite(ctrl *gomock.Controller) *UserCommandSuite {
	r := mock_service.NewMockRegisterUser(ctrl)
	h := handler.NewUserCommand(r)
	return &UserCommandSuite{
		handler:   h,
		registrar: r,
	}
}
