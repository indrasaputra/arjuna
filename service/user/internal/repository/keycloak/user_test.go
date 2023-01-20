package keycloak_test

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	kcsdk "github.com/indrasaputra/arjuna/pkg/sdk/keycloak"
	mock_keycloak "github.com/indrasaputra/arjuna/pkg/sdk/test/mock/keycloak"
	"github.com/indrasaputra/arjuna/service/user/entity"
	"github.com/indrasaputra/arjuna/service/user/internal/repository/keycloak"
)

var (
	testCtx = context.Background()
	testJWT = &kcsdk.JWT{}
)

type UserSuite struct {
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
		st := createUserSuite(ctrl)
		st.client.EXPECT().LoginAdmin(testCtx, st.config.AdminUsername, st.config.AdminPassword).Return(nil, errors.New("error"))

		id, err := st.user.Create(testCtx, user)

		assert.Error(t, err)
		assert.Empty(t, id)
	})

	t.Run("user already exists", func(t *testing.T) {
		st := createUserSuite(ctrl)
		st.client.EXPECT().LoginAdmin(testCtx, st.config.AdminUsername, st.config.AdminPassword).Return(testJWT, nil)
		st.client.EXPECT().CreateUser(testCtx, testJWT.AccessToken, st.config.Realm, gomock.Any()).Return(kcsdk.NewError(http.StatusConflict, ""))

		id, err := st.user.Create(testCtx, user)

		assert.Error(t, err)
		assert.Empty(t, id)
	})

	t.Run("create user returns error", func(t *testing.T) {
		st := createUserSuite(ctrl)
		st.client.EXPECT().LoginAdmin(testCtx, st.config.AdminUsername, st.config.AdminPassword).Return(testJWT, nil)
		st.client.EXPECT().CreateUser(testCtx, testJWT.AccessToken, st.config.Realm, gomock.Any()).Return(errors.New("error"))

		id, err := st.user.Create(testCtx, user)

		assert.Error(t, err)
		assert.Empty(t, id)
	})

	t.Run("get user id returns error", func(t *testing.T) {
		st := createUserSuite(ctrl)
		st.client.EXPECT().LoginAdmin(testCtx, st.config.AdminUsername, st.config.AdminPassword).Return(testJWT, nil)
		st.client.EXPECT().CreateUser(testCtx, testJWT.AccessToken, st.config.Realm, gomock.Any()).Return(nil)
		st.client.EXPECT().GetUserByEmail(testCtx, testJWT.AccessToken, st.config.Realm, gomock.Any()).Return(nil, errors.New("error"))

		id, err := st.user.Create(testCtx, user)

		assert.Error(t, err)
		assert.Empty(t, id)
	})

	t.Run("get user id returns user not found", func(t *testing.T) {
		st := createUserSuite(ctrl)
		st.client.EXPECT().LoginAdmin(testCtx, st.config.AdminUsername, st.config.AdminPassword).Return(testJWT, nil)
		st.client.EXPECT().CreateUser(testCtx, testJWT.AccessToken, st.config.Realm, gomock.Any()).Return(nil)
		st.client.EXPECT().GetUserByEmail(testCtx, testJWT.AccessToken, st.config.Realm, gomock.Any()).Return(nil, kcsdk.NewError(http.StatusNotFound, ""))

		id, err := st.user.Create(testCtx, user)

		assert.Error(t, err)
		assert.Empty(t, id)
	})

	t.Run("success create a new user", func(t *testing.T) {
		st := createUserSuite(ctrl)
		st.client.EXPECT().LoginAdmin(testCtx, st.config.AdminUsername, st.config.AdminPassword).Return(testJWT, nil)
		st.client.EXPECT().CreateUser(testCtx, testJWT.AccessToken, st.config.Realm, gomock.Any()).Return(nil)
		st.client.EXPECT().GetUserByEmail(testCtx, testJWT.AccessToken, st.config.Realm, gomock.Any()).Return(&kcsdk.UserRepresentation{ID: "id"}, nil)

		id, err := st.user.Create(testCtx, user)

		assert.NoError(t, err)
		assert.NotEmpty(t, id)
	})
}

func TestUser_HardDelete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	id := "1"

	t.Run("unable to login as admin to Keycloak", func(t *testing.T) {
		st := createUserSuite(ctrl)
		st.client.EXPECT().LoginAdmin(testCtx, st.config.AdminUsername, st.config.AdminPassword).Return(nil, errors.New("error"))

		err := st.user.HardDelete(testCtx, id)

		assert.Error(t, err)
	})

	t.Run("delete user returns error", func(t *testing.T) {
		st := createUserSuite(ctrl)
		st.client.EXPECT().LoginAdmin(testCtx, st.config.AdminUsername, st.config.AdminPassword).Return(testJWT, nil)
		st.client.EXPECT().DeleteUser(testCtx, testJWT.AccessToken, st.config.Realm, id).Return(kcsdk.NewError(http.StatusInternalServerError, ""))

		err := st.user.HardDelete(testCtx, id)

		assert.Error(t, err)
	})

	t.Run("success delete user", func(t *testing.T) {
		st := createUserSuite(ctrl)
		st.client.EXPECT().LoginAdmin(testCtx, st.config.AdminUsername, st.config.AdminPassword).Return(testJWT, nil)
		st.client.EXPECT().DeleteUser(testCtx, testJWT.AccessToken, st.config.Realm, id).Return(nil)

		err := st.user.HardDelete(testCtx, id)

		assert.NoError(t, err)
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
		Name:     "First User",
		Email:    "first@user.com",
		Password: "strongeststrikerintheworld!",
	}
}

func createUserSuite(ctrl *gomock.Controller) *UserSuite {
	client := mock_keycloak.NewMockKeycloak(ctrl)
	config := createKeycloakConfig(ctrl)
	config.Client = client
	user, _ := keycloak.NewUser(config)

	return &UserSuite{
		user:   user,
		config: config,
		client: client,
	}
}
