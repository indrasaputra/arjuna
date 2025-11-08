package handler_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	apiv1 "github.com/indrasaputra/arjuna/proto/api/v1"
	"github.com/indrasaputra/arjuna/service/user/entity"
	"github.com/indrasaputra/arjuna/service/user/internal/grpc/handler"
	mock_service "github.com/indrasaputra/arjuna/service/user/test/mock/service"
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

	t.Run("nil request is prohibited", func(t *testing.T) {
		st := createUserCommandSuite(ctrl)

		res, err := st.handler.RegisterUser(testCtx, nil)

		assert.Error(t, err)
		assert.Equal(t, entity.ErrEmptyUser(), err)
		assert.Nil(t, res)
	})

	t.Run("empty user is prohibited", func(t *testing.T) {
		st := createUserCommandSuite(ctrl)

		res, err := st.handler.RegisterUser(testCtx, &apiv1.RegisterUserRequest{})

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
			st.registrar.EXPECT().Register(testCtx, gomock.Any()).Return(uuid.Must(uuid.NewV7()), errRet)

			res, err := st.handler.RegisterUser(testCtx, request)

			assert.Error(t, err)
			assert.Equal(t, errRet, err)
			assert.Nil(t, res)
		}
	})

	t.Run("success register user", func(t *testing.T) {
		st := createUserCommandSuite(ctrl)
		id := uuid.Must(uuid.NewV7())
		st.registrar.EXPECT().Register(testCtx, gomock.Any()).Return(id, nil)
		request := &apiv1.RegisterUserRequest{
			User: &apiv1.User{
				Name:     "First User",
				Email:    "first@user.com",
				Password: "BestPlayer",
			},
		}

		res, err := st.handler.RegisterUser(testCtx, request)

		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, id.String(), res.Data.GetId())
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
