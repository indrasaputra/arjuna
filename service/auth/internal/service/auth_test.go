package service_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"golang.org/x/crypto/bcrypt"

	sdklog "github.com/indrasaputra/arjuna/pkg/sdk/log"
	"github.com/indrasaputra/arjuna/service/auth/entity"
	"github.com/indrasaputra/arjuna/service/auth/internal/app"
	"github.com/indrasaputra/arjuna/service/auth/internal/service"
	mock_service "github.com/indrasaputra/arjuna/service/auth/test/mock/service"
)

var (
	testCtx        = context.Background()
	testEmail      = "email@email.com"
	testPassword   = "password"
	testEnv        = "development"
	testSigningKey = "key"
	testExpiry     = 5
)

type AuthSuite struct {
	auth *service.Auth
	repo *mock_service.MockAuthRepository
}

func TestNewAuth(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("successfully create an instance of Auth", func(t *testing.T) {
		st := createAuthSuite(ctrl)
		assert.NotNil(t, st.auth)
	})
}

func TestAuth_Register(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	app.Logger = sdklog.NewLogger(testEnv)

	t.Run("empty account is prohibited", func(t *testing.T) {
		st := createAuthSuite(ctrl)

		err := st.auth.Register(testCtx, nil)

		assert.Error(t, err)
		assert.Equal(t, entity.ErrEmptyAccount(), err)
	})

	t.Run("user id is invalid", func(t *testing.T) {
		st := createAuthSuite(ctrl)
		account := createTestAccount()
		account.UserID = ""

		err := st.auth.Register(testCtx, account)

		assert.Error(t, err)
	})

	t.Run("email is invalid", func(t *testing.T) {
		st := createAuthSuite(ctrl)
		emails := []string{
			"@domain",
			"@domain.com",
			"domain.com",
		}

		for _, email := range emails {
			account := createTestAccount()
			account.Email = email

			err := st.auth.Register(testCtx, account)

			assert.Error(t, err)
			assert.Equal(t, entity.ErrInvalidEmail(), err)
		}
	})

	t.Run("password is invalid", func(t *testing.T) {
		st := createAuthSuite(ctrl)
		account := createTestAccount()
		account.Password = ""

		err := st.auth.Register(testCtx, account)

		assert.Error(t, err)
	})

	t.Run("account repo insert returns not found", func(t *testing.T) {
		st := createAuthSuite(ctrl)
		account := createTestAccount()

		st.repo.EXPECT().Insert(testCtx, account).Return(entity.ErrNotFound())

		err := st.auth.Register(testCtx, account)

		assert.Error(t, err)
	})

	t.Run("account repo insert returns error", func(t *testing.T) {
		st := createAuthSuite(ctrl)
		account := createTestAccount()

		st.repo.EXPECT().Insert(testCtx, account).Return(entity.ErrInternal(""))

		err := st.auth.Register(testCtx, account)

		assert.Error(t, err)
	})

	t.Run("success register an account", func(t *testing.T) {
		st := createAuthSuite(ctrl)
		account := createTestAccount()

		st.repo.EXPECT().Insert(testCtx, account).Return(nil)

		err := st.auth.Register(testCtx, account)

		assert.NoError(t, err)
	})
}

func TestAuth_Login(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	app.Logger = sdklog.NewLogger(testEnv)

	t.Run("param is invalid", func(t *testing.T) {
		type testSuite struct {
			err      error
			email    string
			password string
		}

		tests := []testSuite{
			{email: "", password: "a", err: entity.ErrEmptyField("email")},
			{email: "a", password: "", err: entity.ErrEmptyField("password")},
		}

		st := createAuthSuite(ctrl)
		for _, test := range tests {
			token, err := st.auth.Login(testCtx, test.email, test.password)

			assert.Error(t, err)
			assert.Equal(t, test.err, err)
			assert.Nil(t, token)
		}
	})

	t.Run("repository returns error", func(t *testing.T) {
		st := createAuthSuite(ctrl)
		st.repo.EXPECT().GetByEmail(testCtx, testEmail).Return(nil, entity.ErrInternal("error"))

		token, err := st.auth.Login(testCtx, testEmail, testPassword)

		assert.Error(t, err)
		assert.Nil(t, token)
	})

	t.Run("password is invalid", func(t *testing.T) {
		st := createAuthSuite(ctrl)
		acc := createTestAccount()
		st.repo.EXPECT().GetByEmail(testCtx, testEmail).Return(acc, nil)

		token, err := st.auth.Login(testCtx, testEmail, "testPassword")

		assert.Error(t, err)
		assert.Nil(t, token)
	})

	t.Run("success login", func(t *testing.T) {
		st := createAuthSuite(ctrl)
		acc := createTestAccount()
		st.repo.EXPECT().GetByEmail(testCtx, testEmail).Return(acc, nil)

		token, err := st.auth.Login(testCtx, testEmail, testPassword)

		assert.NoError(t, err)
		assert.NotNil(t, token)
	})
}

func createAuthSuite(ctrl *gomock.Controller) *AuthSuite {
	r := mock_service.NewMockAuthRepository(ctrl)
	a := service.NewAuth(r, []byte(testSigningKey), testExpiry)
	return &AuthSuite{
		auth: a,
		repo: r,
	}
}

func createTestAccount() *entity.Account {
	hash, _ := bcrypt.GenerateFromPassword([]byte(testPassword), bcrypt.MinCost)
	return &entity.Account{
		ID:       "1",
		UserID:   "1",
		Email:    "first@account.com",
		Password: string(hash),
	}
}
