package service_test

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/indrasaputra/arjuna/service/auth/entity"
	"github.com/indrasaputra/arjuna/service/auth/internal/service"
	mock_service "github.com/indrasaputra/arjuna/service/auth/test/mock/service"
)

var (
	testCtx      = context.Background()
	testClientID = "client-id"
	testEmail    = "email@email.com"
	testPassword = "password"
)

type AuthExecutor struct {
	auth *service.Auth
	repo *mock_service.MockAuthRepository
}

func TestNewAuth(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("successfully create an instance of Auth", func(t *testing.T) {
		exec := createAuthExecutor(ctrl)
		assert.NotNil(t, exec.auth)
	})
}

func TestAuth_Login(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("param is invalid", func(t *testing.T) {
		type testSuite struct {
			clientID string
			email    string
			password string
			err      error
		}

		tests := []testSuite{
			{clientID: "", email: "a", password: "a", err: entity.ErrEmptyField("clientId")},
			{clientID: "a", email: "", password: "a", err: entity.ErrEmptyField("email")},
			{clientID: "a", email: "a", password: "", err: entity.ErrEmptyField("password")},
		}

		exec := createAuthExecutor(ctrl)
		for _, test := range tests {
			token, err := exec.auth.Login(testCtx, test.clientID, test.email, test.password)

			assert.Error(t, err)
			assert.Equal(t, test.err, err)
			assert.Nil(t, token)
		}
	})

	t.Run("repository returns error", func(t *testing.T) {
		exec := createAuthExecutor(ctrl)
		exec.repo.EXPECT().Login(testCtx, testClientID, testEmail, testPassword).Return(nil, entity.ErrInternal("error"))

		token, err := exec.auth.Login(testCtx, testClientID, testEmail, testPassword)

		assert.Error(t, err)
		assert.Nil(t, token)
	})

	t.Run("success login", func(t *testing.T) {
		exec := createAuthExecutor(ctrl)
		exec.repo.EXPECT().Login(testCtx, testClientID, testEmail, testPassword).Return(&entity.Token{}, nil)

		token, err := exec.auth.Login(testCtx, testClientID, testEmail, testPassword)

		assert.NoError(t, err)
		assert.NotNil(t, token)
	})
}

func createAuthExecutor(ctrl *gomock.Controller) *AuthExecutor {
	r := mock_service.NewMockAuthRepository(ctrl)
	a := service.NewAuth(r)
	return &AuthExecutor{
		auth: a,
		repo: r,
	}
}
