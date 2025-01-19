package handler_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	sdklog "github.com/indrasaputra/arjuna/pkg/sdk/log"
	apiv1 "github.com/indrasaputra/arjuna/proto/api/v1"
	"github.com/indrasaputra/arjuna/service/user/entity"
	"github.com/indrasaputra/arjuna/service/user/internal/app"
	"github.com/indrasaputra/arjuna/service/user/internal/grpc/handler"
	mock_service "github.com/indrasaputra/arjuna/service/user/test/mock/service"
)

const (
	defaultLimit = uint(10)
)

var (
	testUser = &entity.User{}
	testEnv  = "development"
)

type UserQuerySuite struct {
	handler *handler.UserQuery
	getter  *mock_service.MockGetUser
}

func TestNewUserQuery(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("successful create an instance of UserQuery", func(t *testing.T) {
		st := createUserQuerySuite(ctrl)
		assert.NotNil(t, st.handler)
	})
}

func TestUserQuery_GetAllUsers(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	app.Logger = sdklog.NewLogger(testEnv)

	t.Run("nil request is prohibited", func(t *testing.T) {
		st := createUserQuerySuite(ctrl)

		res, err := st.handler.GetAllUsers(testCtx, nil)

		assert.Error(t, err)
		assert.Equal(t, entity.ErrEmptyUser(), err)
		assert.Nil(t, res)
	})

	t.Run("user service returns error", func(t *testing.T) {
		st := createUserQuerySuite(ctrl)
		st.getter.EXPECT().GetAll(testCtx, defaultLimit).Return([]*entity.User{}, assert.AnError)

		res, err := st.handler.GetAllUsers(testCtx, &apiv1.GetAllUsersRequest{Limit: uint32(defaultLimit)})

		assert.Error(t, err)
		assert.Nil(t, res)
	})

	t.Run("success get all users", func(t *testing.T) {
		st := createUserQuerySuite(ctrl)
		st.getter.EXPECT().GetAll(testCtx, defaultLimit).Return([]*entity.User{testUser, testUser}, nil)

		res, err := st.handler.GetAllUsers(testCtx, &apiv1.GetAllUsersRequest{Limit: uint32(defaultLimit)})

		assert.NoError(t, err)
		assert.NotEmpty(t, res.Data)
		assert.Equal(t, 2, len(res.Data))
	})
}

func createUserQuerySuite(ctrl *gomock.Controller) *UserQuerySuite {
	g := mock_service.NewMockGetUser(ctrl)
	h := handler.NewUserQuery(g)
	return &UserQuerySuite{
		handler: h,
		getter:  g,
	}
}
