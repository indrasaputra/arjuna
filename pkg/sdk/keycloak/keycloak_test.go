package keycloak_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/indrasaputra/arjuna/pkg/sdk/keycloak"
	mock_keycloak "github.com/indrasaputra/arjuna/pkg/sdk/test/mock/keycloak"
)

const (
	baseURL     = "http://localhost:8080/"
	token       = "token"
	realmArjuna = "arjuna"
)

var (
	testCtx    = context.Background()
	errGeneric = errors.New("something wrong")
)

type ClientExecutor struct {
	client *keycloak.Client
	doer   *mock_keycloak.MockDoer
}

func TestNewError(t *testing.T) {
	t.Run("success create an instance of Error", func(t *testing.T) {
		err := keycloak.NewError(http.StatusInternalServerError, "internal server error")
		assert.NotNil(t, err)
	})
}

func TestError_Error(t *testing.T) {
	t.Run("success implement error interface", func(t *testing.T) {
		err := keycloak.NewError(http.StatusInternalServerError, "internal server error")

		_, ok := interface{}(err).(error)

		assert.True(t, ok)
		assert.Equal(t, "internal server error", err.Error())
	})
}

func TestNewClient(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("successfully create an instance of Client", func(t *testing.T) {
		exec := createClientExecutor(ctrl)
		assert.NotNil(t, exec.client)
	})
}

func TestClient_LoginAdmin(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	username := "username"
	password := "password"

	t.Run("doer returns error", func(t *testing.T) {
		exec := createClientExecutor(ctrl)
		exec.doer.EXPECT().Do(gomock.Any()).Return(nil, errGeneric)

		res, err := exec.client.LoginAdmin(testCtx, username, password)

		assert.Error(t, err)
		assert.Nil(t, res)
	})

	t.Run("keycloak returns undecoded error", func(t *testing.T) {
		body := io.NopCloser(strings.NewReader("error"))

		exec := createClientExecutor(ctrl)
		exec.doer.EXPECT().Do(gomock.Any()).Return(&http.Response{Body: body, StatusCode: http.StatusBadRequest}, nil)

		res, err := exec.client.LoginAdmin(testCtx, username, password)

		assert.Error(t, err)
		assert.Nil(t, res)
	})

	t.Run("keycloak returns bad request code", func(t *testing.T) {
		body := io.NopCloser(strings.NewReader(`{"error": "error", "error_description": "desc"}`))

		exec := createClientExecutor(ctrl)
		exec.doer.EXPECT().Do(gomock.Any()).Return(&http.Response{Body: body, StatusCode: http.StatusBadRequest}, nil)

		res, err := exec.client.LoginAdmin(testCtx, username, password)

		assert.Error(t, err)
		assert.Nil(t, res)
	})

	t.Run("fail decode response", func(t *testing.T) {
		body := io.NopCloser(strings.NewReader("something"))

		exec := createClientExecutor(ctrl)
		exec.doer.EXPECT().Do(gomock.Any()).Return(&http.Response{Body: body, StatusCode: http.StatusOK}, nil)

		res, err := exec.client.LoginAdmin(testCtx, username, password)

		assert.Error(t, err)
		assert.Nil(t, res)
	})

	t.Run("success login admin", func(t *testing.T) {
		jwt, _ := json.Marshal(&keycloak.JWT{})
		body := io.NopCloser(bytes.NewReader(jwt))

		exec := createClientExecutor(ctrl)
		exec.doer.EXPECT().Do(gomock.Any()).Return(&http.Response{Body: body, StatusCode: http.StatusOK}, nil)

		res, err := exec.client.LoginAdmin(testCtx, username, password)

		assert.NoError(t, err)
		assert.NotNil(t, res)
	})
}

func TestClient_LoginUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	clientID := "clientID"
	username := "username"
	password := "password"

	t.Run("doer returns error", func(t *testing.T) {
		exec := createClientExecutor(ctrl)
		exec.doer.EXPECT().Do(gomock.Any()).Return(nil, errGeneric)

		res, err := exec.client.LoginUser(testCtx, realmArjuna, clientID, username, password)

		assert.Error(t, err)
		assert.Nil(t, res)
	})

	t.Run("keycloak returns undecoded error", func(t *testing.T) {
		body := io.NopCloser(strings.NewReader("error"))

		exec := createClientExecutor(ctrl)
		exec.doer.EXPECT().Do(gomock.Any()).Return(&http.Response{Body: body, StatusCode: http.StatusBadRequest}, nil)

		res, err := exec.client.LoginUser(testCtx, realmArjuna, clientID, username, password)

		assert.Error(t, err)
		assert.Nil(t, res)
	})

	t.Run("keycloak returns bad request code", func(t *testing.T) {
		body := io.NopCloser(strings.NewReader(`{"error": "error", "error_description": "desc"}`))

		exec := createClientExecutor(ctrl)
		exec.doer.EXPECT().Do(gomock.Any()).Return(&http.Response{Body: body, StatusCode: http.StatusBadRequest}, nil)

		res, err := exec.client.LoginUser(testCtx, realmArjuna, clientID, username, password)

		assert.Error(t, err)
		assert.Nil(t, res)
	})

	t.Run("fail decode response", func(t *testing.T) {
		body := io.NopCloser(strings.NewReader("something"))

		exec := createClientExecutor(ctrl)
		exec.doer.EXPECT().Do(gomock.Any()).Return(&http.Response{Body: body, StatusCode: http.StatusOK}, nil)

		res, err := exec.client.LoginUser(testCtx, realmArjuna, clientID, username, password)

		assert.Error(t, err)
		assert.Nil(t, res)
	})

	t.Run("success login user", func(t *testing.T) {
		jwt, _ := json.Marshal(&keycloak.JWT{})
		body := io.NopCloser(bytes.NewReader(jwt))

		exec := createClientExecutor(ctrl)
		exec.doer.EXPECT().Do(gomock.Any()).Return(&http.Response{Body: body, StatusCode: http.StatusOK}, nil)

		res, err := exec.client.LoginUser(testCtx, realmArjuna, clientID, username, password)

		assert.NoError(t, err)
		assert.NotNil(t, res)
	})
}

func TestClient_CreateRealm(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	realm := &keycloak.RealmRepresentation{}
	body := io.NopCloser(bytes.NewReader([]byte(`{}`)))

	t.Run("doer returns error", func(t *testing.T) {
		exec := createClientExecutor(ctrl)
		exec.doer.EXPECT().Do(gomock.Any()).Return(nil, errGeneric)

		err := exec.client.CreateRealm(testCtx, token, realm)

		assert.Error(t, err)
	})

	t.Run("unauthorized request", func(t *testing.T) {
		exec := createClientExecutor(ctrl)
		exec.doer.EXPECT().Do(gomock.Any()).Return(&http.Response{StatusCode: http.StatusUnauthorized, Body: body}, nil)

		err := exec.client.CreateRealm(testCtx, token, realm)

		assert.Error(t, err)
	})

	t.Run("keycloak doesn't respond with 201 status code", func(t *testing.T) {
		exec := createClientExecutor(ctrl)
		exec.doer.EXPECT().Do(gomock.Any()).Return(&http.Response{StatusCode: http.StatusInternalServerError, Body: body}, nil)

		err := exec.client.CreateRealm(testCtx, token, realm)

		assert.Error(t, err)
	})

	t.Run("successfully create a new realm", func(t *testing.T) {
		exec := createClientExecutor(ctrl)
		exec.doer.EXPECT().Do(gomock.Any()).Return(&http.Response{StatusCode: http.StatusCreated, Body: body}, nil)

		err := exec.client.CreateRealm(testCtx, token, realm)

		assert.NoError(t, err)
	})
}

func TestClient_CreateClient(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	client := &keycloak.ClientRepresentation{}
	body := io.NopCloser(bytes.NewReader([]byte(`{}`)))

	t.Run("doer returns error", func(t *testing.T) {
		exec := createClientExecutor(ctrl)
		exec.doer.EXPECT().Do(gomock.Any()).Return(nil, errGeneric)

		err := exec.client.CreateClient(testCtx, token, realmArjuna, client)

		assert.Error(t, err)
	})

	t.Run("keycloak doesn't respond with 201 status code", func(t *testing.T) {
		exec := createClientExecutor(ctrl)
		exec.doer.EXPECT().Do(gomock.Any()).Return(&http.Response{StatusCode: http.StatusInternalServerError, Body: body}, nil)

		err := exec.client.CreateClient(testCtx, token, realmArjuna, client)

		assert.Error(t, err)
	})

	t.Run("successfully create a new client", func(t *testing.T) {
		exec := createClientExecutor(ctrl)
		exec.doer.EXPECT().Do(gomock.Any()).Return(&http.Response{StatusCode: http.StatusCreated, Body: body}, nil)

		err := exec.client.CreateClient(testCtx, token, realmArjuna, client)

		assert.NoError(t, err)
	})
}

func TestClient_CreateUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	user := &keycloak.UserRepresentation{}
	body := io.NopCloser(bytes.NewReader([]byte(`{}`)))

	t.Run("doer returns error", func(t *testing.T) {
		exec := createClientExecutor(ctrl)
		exec.doer.EXPECT().Do(gomock.Any()).Return(nil, errGeneric)

		err := exec.client.CreateUser(testCtx, token, realmArjuna, user)

		assert.Error(t, err)
	})

	t.Run("keycloak with 409 status code", func(t *testing.T) {
		exec := createClientExecutor(ctrl)
		exec.doer.EXPECT().Do(gomock.Any()).Return(&http.Response{StatusCode: http.StatusConflict, Body: body}, nil)

		err := exec.client.CreateUser(testCtx, token, realmArjuna, user)

		assert.Error(t, err)
	})

	t.Run("keycloak doesn't respond with 201 status code", func(t *testing.T) {
		exec := createClientExecutor(ctrl)
		exec.doer.EXPECT().Do(gomock.Any()).Return(&http.Response{StatusCode: http.StatusInternalServerError, Body: body}, nil)

		err := exec.client.CreateUser(testCtx, token, realmArjuna, user)

		assert.Error(t, err)
	})

	t.Run("successfully create a new user", func(t *testing.T) {
		exec := createClientExecutor(ctrl)
		exec.doer.EXPECT().Do(gomock.Any()).Return(&http.Response{StatusCode: http.StatusCreated, Body: body}, nil)

		err := exec.client.CreateUser(testCtx, token, realmArjuna, user)

		assert.NoError(t, err)
	})
}

func TestClient_DeleteUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	id := "1"

	t.Run("delete api returns error", func(t *testing.T) {
		exec := createClientExecutor(ctrl)
		exec.doer.EXPECT().Do(gomock.Any()).Return(nil, errors.New("error"))

		err := exec.client.DeleteUser(testCtx, token, realmArjuna, id)

		assert.Error(t, err)
	})

	t.Run("success delete user", func(t *testing.T) {
		body := io.NopCloser(bytes.NewReader([]byte(`{}`)))
		exec := createClientExecutor(ctrl)
		exec.doer.EXPECT().Do(gomock.Any()).Return(&http.Response{StatusCode: http.StatusNoContent, Body: body}, nil)

		err := exec.client.DeleteUser(testCtx, token, realmArjuna, id)

		assert.NoError(t, err)
	})
}

func TestClient_GetUserByEmail(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	email := "arjuna@arjuna.com"

	t.Run("doer returns error", func(t *testing.T) {
		exec := createClientExecutor(ctrl)
		exec.doer.EXPECT().Do(gomock.Any()).Return(nil, errGeneric)

		user, err := exec.client.GetUserByEmail(testCtx, token, realmArjuna, email)

		assert.Error(t, err)
		assert.Nil(t, user)
	})

	t.Run("user not found", func(t *testing.T) {
		body := io.NopCloser(bytes.NewReader([]byte(`{}`)))
		exec := createClientExecutor(ctrl)
		exec.doer.EXPECT().Do(gomock.Any()).Return(&http.Response{StatusCode: http.StatusOK, Body: body}, nil)

		user, err := exec.client.GetUserByEmail(testCtx, token, realmArjuna, email)

		assert.Error(t, err)
		assert.Nil(t, user)
	})

	t.Run("success find user", func(t *testing.T) {
		body := io.NopCloser(bytes.NewReader([]byte(`[{"id": "abc", "email": "admin@arjuna.com"}]`)))
		exec := createClientExecutor(ctrl)
		exec.doer.EXPECT().Do(gomock.Any()).Return(&http.Response{StatusCode: http.StatusOK, Body: body}, nil)

		user, err := exec.client.GetUserByEmail(testCtx, token, realmArjuna, email)

		assert.NoError(t, err)
		assert.NotNil(t, user)
	})
}

func TestClient_GetAllUsers(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("get users returns error", func(t *testing.T) {
		exec := createClientExecutor(ctrl)
		exec.doer.EXPECT().Do(gomock.Any()).Return(nil, errGeneric)

		users, err := exec.client.GetAllUsers(testCtx, token, realmArjuna)

		assert.Error(t, err)
		assert.Nil(t, users)
	})

	t.Run("success get all users", func(t *testing.T) {
		body := io.NopCloser(bytes.NewReader([]byte(`[{"id": "5c44f049-8ab2-4d0f-b41d-7b08f467e817"}]`)))
		exec := createClientExecutor(ctrl)
		exec.doer.EXPECT().Do(gomock.Any()).Return(&http.Response{StatusCode: http.StatusOK, Body: body}, nil)

		users, err := exec.client.GetAllUsers(testCtx, token, realmArjuna)

		assert.NoError(t, err)
		assert.NotEmpty(t, users)
	})
}

func createClientExecutor(ctrl *gomock.Controller) *ClientExecutor {
	doer := mock_keycloak.NewMockDoer(ctrl)
	client := keycloak.NewClient(doer, baseURL)
	return &ClientExecutor{
		client: client,
		doer:   doer,
	}
}
