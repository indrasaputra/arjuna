package handler_test

import (
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	sdklog "github.com/indrasaputra/arjuna/pkg/sdk/log"
	apiv1 "github.com/indrasaputra/arjuna/proto/api/v1"
	"github.com/indrasaputra/arjuna/service/auth/entity"
	"github.com/indrasaputra/arjuna/service/auth/internal/app"
	"github.com/indrasaputra/arjuna/service/auth/internal/grpc/handler"
	mock_service "github.com/indrasaputra/arjuna/service/auth/test/mock/service"
)

var (
	testClientID = "client-id"
	testEmail    = "email@email.com"
	testPassword = "password"
	testEnv      = "development"
)

type AuthSuite struct {
	handler *handler.Auth
	auth    *mock_service.MockAuthentication
}

func TestNewAuth(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("successful create an instance of Auth", func(t *testing.T) {
		st := createAuthSuite(ctrl)
		assert.NotNil(t, st.handler)
	})
}

func TestAuth_Login(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	app.Logger = sdklog.NewLogger(testEnv)

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

		st := createAuthSuite(ctrl)
		for _, test := range tests {
			res, err := st.handler.Login(testCtx, test.request)

			assert.Error(t, err)
			assert.Equal(t, test.err, err)
			assert.Nil(t, res)
		}
	})

	t.Run("auth service returns error", func(t *testing.T) {
		st := createAuthSuite(ctrl)
		st.auth.EXPECT().Login(testCtx, testClientID, testEmail, testPassword).Return(nil, errors.New("error"))

		req := &apiv1.LoginRequest{Credential: &apiv1.Credential{ClientId: testClientID, Email: testEmail, Password: testPassword}}
		res, err := st.handler.Login(testCtx, req)

		assert.Error(t, err)
		assert.Nil(t, res)
	})

	t.Run("success login", func(t *testing.T) {
		st := createAuthSuite(ctrl)
		st.auth.EXPECT().Login(testCtx, testClientID, testEmail, testPassword).Return(&entity.Token{}, nil)

		req := &apiv1.LoginRequest{Credential: &apiv1.Credential{ClientId: testClientID, Email: testEmail, Password: testPassword}}
		res, err := st.handler.Login(testCtx, req)

		assert.NoError(t, err)
		assert.NotNil(t, res)
	})
}

func createAuthSuite(ctrl *gomock.Controller) *AuthSuite {
	r := mock_service.NewMockAuthentication(ctrl)
	h := handler.NewAuth(r)
	return &AuthSuite{
		handler: h,
		auth:    r,
	}
}
