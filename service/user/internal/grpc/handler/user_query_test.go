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

const (
	defaultLimit = uint(10)
)

var (
	testUser = &entity.User{}
)

type UserQueryExecutor struct {
	handler *handler.UserQuery
	getter  *mock_service.MockGetUser
}

func TestNewUserQuery(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("successful create an instance of UserQuery", func(t *testing.T) {
		exec := createUserQueryExecutor(ctrl)
		assert.NotNil(t, exec.handler)
	})
}

func TestUserQuery_GetAllUsers(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("nil request is prohibited", func(t *testing.T) {
		exec := createUserQueryExecutor(ctrl)

		res, err := exec.handler.GetAllUsers(testCtx, nil)

		assert.Error(t, err)
		assert.Equal(t, entity.ErrEmptyUser(), err)
		assert.Nil(t, res)
	})

	t.Run("user service returns error", func(t *testing.T) {
		exec := createUserQueryExecutor(ctrl)
		exec.getter.EXPECT().GetAll(testCtx, defaultLimit).Return([]*entity.User{}, entity.ErrInternal(""))

		res, err := exec.handler.GetAllUsers(testCtx, &apiv1.GetAllUsersRequest{Limit: uint32(defaultLimit)})

		assert.Error(t, err)
		assert.Equal(t, entity.ErrInternal(""), err)
		assert.Nil(t, res)
	})

	t.Run("success get all users", func(t *testing.T) {
		exec := createUserQueryExecutor(ctrl)
		exec.getter.EXPECT().GetAll(testCtx, defaultLimit).Return([]*entity.User{testUser, testUser}, nil)

		res, err := exec.handler.GetAllUsers(testCtx, &apiv1.GetAllUsersRequest{Limit: uint32(defaultLimit)})

		assert.NoError(t, err)
		assert.NotEmpty(t, res.Data)
		assert.Equal(t, 2, len(res.Data))
	})
}

func createUserQueryExecutor(ctrl *gomock.Controller) *UserQueryExecutor {
	g := mock_service.NewMockGetUser(ctrl)
	h := handler.NewUserQuery(g)
	return &UserQueryExecutor{
		handler: h,
		getter:  g,
	}
}