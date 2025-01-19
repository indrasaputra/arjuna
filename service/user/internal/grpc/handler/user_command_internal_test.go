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

type UserCommandInternalSuite struct {
	handler *handler.UserCommandInternal
	deleter *mock_service.MockDeleteUser
}

func TestNewUserCommandInternal(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("successful create an instance of UserCommandInternal", func(t *testing.T) {
		st := createUserCommandInternalSuite(ctrl)
		assert.NotNil(t, st.handler)
	})
}

func TestUserCommandInternal_DeleteUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("nil request is prohibited", func(t *testing.T) {
		st := createUserCommandInternalSuite(ctrl)

		res, err := st.handler.DeleteUser(testCtx, nil)

		assert.Error(t, err)
		assert.Equal(t, entity.ErrEmptyUser(), err)
		assert.Nil(t, res)
	})

	t.Run("empty user is prohibited", func(t *testing.T) {
		st := createUserCommandInternalSuite(ctrl)

		res, err := st.handler.DeleteUser(testCtx, &apiv1.DeleteUserRequest{})

		assert.Error(t, err)
		assert.Equal(t, entity.ErrEmptyUser(), err)
		assert.Nil(t, res)
	})

	t.Run("user deleter service returns error", func(t *testing.T) {
		st := createUserCommandInternalSuite(ctrl)
		id := uuid.Must(uuid.NewV7())
		request := &apiv1.DeleteUserRequest{Id: id.String()}

		errors := []error{
			entity.ErrEmptyUser(),
			entity.ErrAlreadyExists(),
			entity.ErrInvalidName(),
			entity.ErrInvalidEmail(),
			entity.ErrInternal("error"),
		}
		for _, errRet := range errors {
			st.deleter.EXPECT().HardDelete(testCtx, id).Return(errRet)

			res, err := st.handler.DeleteUser(testCtx, request)

			assert.Error(t, err)
			assert.Equal(t, errRet, err)
			assert.Nil(t, res)
		}
	})

	t.Run("success delete user", func(t *testing.T) {
		id := uuid.Must(uuid.NewV7())
		request := &apiv1.DeleteUserRequest{Id: id.String()}
		st := createUserCommandInternalSuite(ctrl)
		st.deleter.EXPECT().HardDelete(testCtx, id).Return(nil)

		res, err := st.handler.DeleteUser(testCtx, request)

		assert.NoError(t, err)
		assert.NotNil(t, res)
	})
}

func createUserCommandInternalSuite(ctrl *gomock.Controller) *UserCommandInternalSuite {
	d := mock_service.NewMockDeleteUser(ctrl)
	h := handler.NewUserCommandInternal(d)
	return &UserCommandInternalSuite{
		handler: h,
		deleter: d,
	}
}
