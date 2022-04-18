package handler_test

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	apiv1 "github.com/indrasaputra/arjuna/proto/api/v1"
	"github.com/indrasaputra/arjuna/service/user/entity"
	"github.com/indrasaputra/arjuna/service/user/internal/grpc/handler"
	mock_service "github.com/indrasaputra/arjuna/service/user/test/mock/service"
)

type UserCommandExecutor struct {
	handler     *handler.UserCommand
	registrator *mock_service.MockRegisterUser
}

func TestNewUserCommand(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("successful create an instance of UserCommand", func(t *testing.T) {
		exec := createUserCommandExecutor(ctrl)
		assert.NotNil(t, exec.handler)
	})
}

func TestUserCommand_RegisterUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("nil request is prohibited", func(t *testing.T) {
		exec := createUserCommandExecutor(ctrl)

		res, err := exec.handler.RegisterUser(testCtx, nil)

		assert.Error(t, err)
		assert.Equal(t, entity.ErrEmptyUser(), err)
		assert.Nil(t, res)
	})

	t.Run("empty user is prohibited", func(t *testing.T) {
		exec := createUserCommandExecutor(ctrl)

		res, err := exec.handler.RegisterUser(testCtx, &apiv1.RegisterUserRequest{})

		assert.Error(t, err)
		assert.Equal(t, entity.ErrEmptyUser(), err)
		assert.Nil(t, res)
	})

	t.Run("user service returns error", func(t *testing.T) {
		exec := createUserCommandExecutor(ctrl)
		request := &apiv1.RegisterUserRequest{
			User: &apiv1.User{
				Name:     "Zlatan Ibrahimovic",
				Email:    "zlatan@ibrahimovic.com",
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
			exec.registrator.EXPECT().Register(testCtx, gomock.Any()).Return("", errRet)

			res, err := exec.handler.RegisterUser(testCtx, request)

			assert.Error(t, err)
			assert.Equal(t, errRet, err)
			assert.Nil(t, res)
		}
	})

	t.Run("success register user", func(t *testing.T) {
		exec := createUserCommandExecutor(ctrl)
		exec.registrator.EXPECT().Register(testCtx, gomock.Any()).Return("id", nil)
		request := &apiv1.RegisterUserRequest{
			User: &apiv1.User{
				Name:     "Zlatan Ibrahimovic",
				Email:    "zlatan@ibrahimovic.com",
				Password: "BestPlayer",
			},
		}

		res, err := exec.handler.RegisterUser(testCtx, request)

		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, "id", res.Data.GetId())
	})
}

func createUserCommandExecutor(ctrl *gomock.Controller) *UserCommandExecutor {
	r := mock_service.NewMockRegisterUser(ctrl)
	h := handler.NewUserCommand(r)
	return &UserCommandExecutor{
		handler:     h,
		registrator: r,
	}
}
