package keycloak_test

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	kcsdk "github.com/indrasaputra/arjuna/pkg/sdk/keycloak"
	sdklog "github.com/indrasaputra/arjuna/pkg/sdk/log"
	mock_keycloak "github.com/indrasaputra/arjuna/pkg/sdk/test/mock/keycloak"
	"github.com/indrasaputra/arjuna/service/auth/internal/app"
	"github.com/indrasaputra/arjuna/service/auth/internal/repository/keycloak"
)

var (
	testCtx      = context.Background()
	testClientID = "client-id"
	testEmail    = "email@email.com"
	testPassword = "password"
	testEnv      = "development"
)

type AuthSuite struct {
	auth   *keycloak.Auth
	config *keycloak.Config
	client *mock_keycloak.MockKeycloak
}

func TestConfig_Validate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("client can't be nil", func(t *testing.T) {
		cfg := createKeycloakConfig(ctrl)
		cfg.Client = nil

		err := cfg.Validate()

		assert.Error(t, err)
	})

	t.Run("realm can't be empty", func(t *testing.T) {
		cfg := createKeycloakConfig(ctrl)
		cfg.Realm = "  "

		err := cfg.Validate()

		assert.Error(t, err)
	})
}

func TestNewAuth(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("invalid config is not accepted", func(t *testing.T) {
		cfg := createKeycloakConfig(ctrl)
		cfg.Client = nil

		user, err := keycloak.NewAuth(cfg)

		assert.Error(t, err)
		assert.Nil(t, user)
	})

	t.Run("success create a new auth repository with valid config", func(t *testing.T) {
		cfg := createKeycloakConfig(ctrl)

		user, err := keycloak.NewAuth(cfg)

		assert.NoError(t, err)
		assert.NotNil(t, user)
	})
}

func TestAuth_Login(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	app.Logger = sdklog.NewLogger(testEnv)

	t.Run("login returns error", func(t *testing.T) {
		errs := []error{
			errors.New("error"),
			kcsdk.NewError(http.StatusUnauthorized, ""),
			kcsdk.NewError(http.StatusBadRequest, ""),
			kcsdk.NewError(http.StatusInternalServerError, ""),
		}

		st := createAuthSuite(ctrl)
		for _, err := range errs {
			st.client.EXPECT().LoginUser(testCtx, st.config.Realm, testClientID, testEmail, testPassword).Return(nil, err)

			token, err := st.auth.Login(testCtx, testClientID, testEmail, testPassword)

			assert.Error(t, err)
			assert.Nil(t, token)
		}
	})

	t.Run("success login", func(t *testing.T) {
		st := createAuthSuite(ctrl)
		st.client.EXPECT().LoginUser(testCtx, st.config.Realm, testClientID, testEmail, testPassword).Return(&kcsdk.JWT{}, nil)

		token, err := st.auth.Login(testCtx, testClientID, testEmail, testPassword)

		assert.NoError(t, err)
		assert.NotNil(t, token)
	})
}

func createKeycloakConfig(ctrl *gomock.Controller) *keycloak.Config {
	return &keycloak.Config{
		Client: mock_keycloak.NewMockKeycloak(ctrl),
		Realm:  "realm",
	}
}

func createAuthSuite(ctrl *gomock.Controller) *AuthSuite {
	client := mock_keycloak.NewMockKeycloak(ctrl)
	config := createKeycloakConfig(ctrl)
	config.Client = client
	auth, _ := keycloak.NewAuth(config)

	return &AuthSuite{
		auth:   auth,
		config: config,
		client: client,
	}
}
