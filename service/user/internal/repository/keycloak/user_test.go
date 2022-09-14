package keycloak_test

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/twilio/twilio-go/client/jwt"

	kcsdk "github.com/indrasaputra/arjuna/pkg/sdk/keycloak"
	mock_keycloak "github.com/indrasaputra/arjuna/pkg/sdk/test/mock/keycloak"
	"github.com/indrasaputra/arjuna/service/user/entity"
	"github.com/indrasaputra/arjuna/service/user/internal/repository/keycloak"
)

var (
	testCtx = context.Background()
	testJWT = &kcsdk.JWT{}
)

type UserExecutor struct {
	user   *keycloak.User
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

	t.Run("admin username can't be empty", func(t *testing.T) {
		cfg := createKeycloakConfig(ctrl)
		cfg.AdminUsername = "  "

		err := cfg.Validate()

		assert.Error(t, err)
	})

	t.Run("admin password can't be empty", func(t *testing.T) {
		cfg := createKeycloakConfig(ctrl)
		cfg.AdminPassword = "  "

		err := cfg.Validate()

		assert.Error(t, err)
	})
}

func TestNewUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("invalid config is not accepted", func(t *testing.T) {
		cfg := createKeycloakConfig(ctrl)
		cfg.Client = nil

		user, err := keycloak.NewUser(cfg)

		assert.Error(t, err)
		assert.Nil(t, user)
	})

	t.Run("success create a new user with valid config", func(t *testing.T) {
		cfg := createKeycloakConfig(ctrl)

		user, err := keycloak.NewUser(cfg)

		assert.NoError(t, err)
		assert.NotNil(t, user)
	})
}

func TestUser_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	user := createUserEntity()

	t.Run("unable to login as admin to Keycloak", func(t *testing.T) {
		exec := createUserExecutor(ctrl)
		exec.client.EXPECT().LoginAdmin(testCtx, exec.config.AdminUsername, exec.config.AdminPassword).Return(nil, errors.New("error"))

		id, err := exec.user.Create(testCtx, user)

		assert.Error(t, err)
		assert.Empty(t, id)
	})

	t.Run("user already exists", func(t *testing.T) {
		exec := createUserExecutor(ctrl)
		exec.client.EXPECT().LoginAdmin(testCtx, exec.config.AdminUsername, exec.config.AdminPassword).Return(testJWT, nil)
		exec.client.EXPECT().CreateUser(testCtx, testJWT.AccessToken, exec.config.Realm, gomock.Any()).Return(kcsdk.ErrConflict)

		id, err := exec.user.Create(testCtx, user)

		assert.Error(t, err)
		assert.Empty(t, id)
	})

	t.Run("create user returns error", func(t *testing.T) {
		exec := createUserExecutor(ctrl)
		exec.client.EXPECT().LoginAdmin(testCtx, exec.config.AdminUsername, exec.config.AdminPassword).Return(testJWT, nil)
		exec.client.EXPECT().CreateUser(testCtx, testJWT.AccessToken, exec.config.Realm, gomock.Any()).Return(errors.New("error"))

		id, err := exec.user.Create(testCtx, user)

		assert.Error(t, err)
		assert.Empty(t, id)
	})

	t.Run("get user id returns error", func(t *testing.T) {
		exec := createUserExecutor(ctrl)
		exec.client.EXPECT().LoginAdmin(testCtx, exec.config.AdminUsername, exec.config.AdminPassword).Return(jwt, nil)
		exec.client.EXPECT().CreateUser(testCtx, jwt.AccessToken, exec.config.Realm, gomock.Any()).Return(nil)
		exec.client.EXPECT().GetUserByEmail(testCtx, jwt.AccessToken, exec.config.Realm, gomock.Any()).Return(nil, errors.New("error"))

		id, err := exec.user.Create(testCtx, user)

		assert.Error(t, err)
		assert.Empty(t, id)
	})

	t.Run("get user id returns user not found", func(t *testing.T) {
		exec := createUserExecutor(ctrl)
		exec.client.EXPECT().LoginAdmin(testCtx, exec.config.AdminUsername, exec.config.AdminPassword).Return(jwt, nil)
		exec.client.EXPECT().CreateUser(testCtx, jwt.AccessToken, exec.config.Realm, gomock.Any()).Return(nil)
		exec.client.EXPECT().GetUserByEmail(testCtx, jwt.AccessToken, exec.config.Realm, gomock.Any()).Return(nil, kcsdk.ErrUserNotFound)

		id, err := exec.user.Create(testCtx, user)

		assert.Error(t, err)
		assert.Empty(t, id)
	})

	t.Run("success create a new user", func(t *testing.T) {
		exec := createUserExecutor(ctrl)
		exec.client.EXPECT().LoginAdmin(testCtx, exec.config.AdminUsername, exec.config.AdminPassword).Return(jwt, nil)
		exec.client.EXPECT().CreateUser(testCtx, jwt.AccessToken, exec.config.Realm, gomock.Any()).Return(nil)
		exec.client.EXPECT().GetUserByEmail(testCtx, jwt.AccessToken, exec.config.Realm, gomock.Any()).Return(&kcsdk.UserRepresentation{ID: "id"}, nil)

		id, err := exec.user.Create(testCtx, user)

		assert.NoError(t, err)
		assert.NotEmpty(t, id)
	})
}

func TestUser_GetAll(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("unable to login as admin to Keycloak", func(t *testing.T) {
		exec := createUserExecutor(ctrl)
		exec.client.EXPECT().LoginAdmin(testCtx, exec.config.AdminUsername, exec.config.AdminPassword).Return(nil, errors.New("error"))

		users, err := exec.user.GetAll(testCtx)

		assert.Error(t, err)
		assert.Empty(t, users)
	})

	t.Run("get all users from keycloak returns error", func(t *testing.T) {
		exec := createUserExecutor(ctrl)
		exec.client.EXPECT().LoginAdmin(testCtx, exec.config.AdminUsername, exec.config.AdminPassword).Return(testJWT, nil)
		exec.client.EXPECT().GetAllUsers(testCtx, testJWT.AccessToken, exec.config.Realm).Return(nil, errors.New("error"))

		users, err := exec.user.GetAll(testCtx)

		assert.Error(t, err)
		assert.Empty(t, users)
	})

	t.Run("success get all users", func(t *testing.T) {
		exec := createUserExecutor(ctrl)
		exec.client.EXPECT().LoginAdmin(testCtx, exec.config.AdminUsername, exec.config.AdminPassword).Return(testJWT, nil)
		exec.client.EXPECT().GetAllUsers(testCtx, testJWT.AccessToken, exec.config.Realm).Return([]*kcsdk.UserRepresentation{{}}, nil)

		users, err := exec.user.GetAll(testCtx)

		assert.NoError(t, err)
		assert.NotEmpty(t, users)
	})
}

func createKeycloakConfig(ctrl *gomock.Controller) *keycloak.Config {
	return &keycloak.Config{
		Client:        mock_keycloak.NewMockKeycloak(ctrl),
		Realm:         "realm",
		AdminUsername: "admin",
		AdminPassword: "password",
	}
}

func createUserEntity() *entity.User {
	return &entity.User{
		Name:     "Zlatan Ibrahimovic",
		Email:    "zlatan@ibrahimovic.com",
		Password: "strongeststrikerintheworld!",
	}
}

func createUserExecutor(ctrl *gomock.Controller) *UserExecutor {
	client := mock_keycloak.NewMockKeycloak(ctrl)
	config := createKeycloakConfig(ctrl)
	config.Client = client
	user, _ := keycloak.NewUser(config)

	return &UserExecutor{
		user:   user,
		config: config,
		client: client,
	}
}
