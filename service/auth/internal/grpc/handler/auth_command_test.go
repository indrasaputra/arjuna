package handler_test

import (
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	apiv1 "github.com/indrasaputra/arjuna/proto/api/v1"
	"github.com/indrasaputra/arjuna/service/auth/entity"
	"github.com/indrasaputra/arjuna/service/auth/internal/grpc/handler"
	mock_service "github.com/indrasaputra/arjuna/service/auth/test/mock/service"
)

var (
	testClientID = "client-id"
	testEmail    = "email@email.com"
	testPassword = "password"
)

type AuthExecutor struct {
	handler *handler.Auth
	auth    *mock_service.MockAuthentication
}

func TestNewAuth(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("successful create an instance of Auth", func(t *testing.T) {
		exec := createAuthExecutor(ctrl)
		assert.NotNil(t, exec.handler)
	})
}

func TestAuth_Login(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("request is invalid", func(t *testing.T) {
		type testSuite struct {
			request *apiv1.LoginRequest
			err     error
		}

		tests := []testSuite{
			{request: nil, err: entity.ErrEmptyField("request body")},
			{request: &apiv1.LoginRequest{Credential: nil}, err: entity.ErrEmptyField("request body")},
			{request: &apiv1.LoginRequest{Credential: &apiv1.Credential{ClientId: "", Email: "a", Password: "a"}}, err: entity.ErrEmptyField("client id")},
			{request: &apiv1.LoginRequest{Credential: &apiv1.Credential{ClientId: "a", Email: "", Password: "a"}}, err: entity.ErrEmptyField("email")},
			{request: &apiv1.LoginRequest{Credential: &apiv1.Credential{ClientId: "a", Email: "a", Password: ""}}, err: entity.ErrEmptyField("password")},
		}

		exec := createAuthExecutor(ctrl)
		for _, test := range tests {
			res, err := exec.handler.Login(testCtx, test.request)

			assert.Error(t, err)
			assert.Equal(t, test.err, err)
			assert.Nil(t, res)
		}
	})

	t.Run("auth service returns error", func(t *testing.T) {
		exec := createAuthExecutor(ctrl)
		exec.auth.EXPECT().Login(testCtx, testClientID, testEmail, testPassword).Return(nil, errors.New("error"))

		req := &apiv1.LoginRequest{Credential: &apiv1.Credential{ClientId: testClientID, Email: testEmail, Password: testPassword}}
		res, err := exec.handler.Login(testCtx, req)

		assert.Error(t, err)
		assert.Nil(t, res)
	})

	t.Run("success login", func(t *testing.T) {
		exec := createAuthExecutor(ctrl)
		exec.auth.EXPECT().Login(testCtx, testClientID, testEmail, testPassword).Return(&entity.Token{}, nil)

		req := &apiv1.LoginRequest{Credential: &apiv1.Credential{ClientId: testClientID, Email: testEmail, Password: testPassword}}
		res, err := exec.handler.Login(testCtx, req)

		assert.NoError(t, err)
		assert.NotNil(t, res)
	})
}

func createAuthExecutor(ctrl *gomock.Controller) *AuthExecutor {
	r := mock_service.NewMockAuthentication(ctrl)
	h := handler.NewAuth(r)
	return &AuthExecutor{
		handler: h,
		auth:    r,
	}
}
