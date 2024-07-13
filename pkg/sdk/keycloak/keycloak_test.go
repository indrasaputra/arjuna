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

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

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

type ClientSuite struct {
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
		st := createClientSuite(ctrl)
		assert.NotNil(t, st.client)
	})
}

func TestClient_LoginAdmin(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	username := "username"
	password := "password"

	t.Run("doer returns error", func(t *testing.T) {
		st := createClientSuite(ctrl)
		st.doer.EXPECT().Do(gomock.Any()).Return(nil, errGeneric)

		res, err := st.client.LoginAdmin(testCtx, username, password)

		assert.Error(t, err)
		assert.Nil(t, res)
	})

	t.Run("keycloak returns undecoded error", func(t *testing.T) {
		body := io.NopCloser(strings.NewReader("error"))

		st := createClientSuite(ctrl)
		st.doer.EXPECT().Do(gomock.Any()).Return(&http.Response{Body: body, StatusCode: http.StatusBadRequest}, nil)

		res, err := st.client.LoginAdmin(testCtx, username, password)

		assert.Error(t, err)
		assert.Nil(t, res)
	})

	t.Run("keycloak returns bad request code", func(t *testing.T) {
		body := io.NopCloser(strings.NewReader(`{"error": "error", "error_description": "desc"}`))

		st := createClientSuite(ctrl)
		st.doer.EXPECT().Do(gomock.Any()).Return(&http.Response{Body: body, StatusCode: http.StatusBadRequest}, nil)

		res, err := st.client.LoginAdmin(testCtx, username, password)

		assert.Error(t, err)
		assert.Nil(t, res)
	})

	t.Run("fail decode response", func(t *testing.T) {
		body := io.NopCloser(strings.NewReader("something"))

		st := createClientSuite(ctrl)
		st.doer.EXPECT().Do(gomock.Any()).Return(&http.Response{Body: body, StatusCode: http.StatusOK}, nil)

		res, err := st.client.LoginAdmin(testCtx, username, password)

		assert.Error(t, err)
		assert.Nil(t, res)
	})

	t.Run("success login admin", func(t *testing.T) {
		jwt, _ := json.Marshal(&keycloak.JWT{})
		body := io.NopCloser(bytes.NewReader(jwt))

		st := createClientSuite(ctrl)
		st.doer.EXPECT().Do(gomock.Any()).Return(&http.Response{Body: body, StatusCode: http.StatusOK}, nil)

		res, err := st.client.LoginAdmin(testCtx, username, password)

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
		st := createClientSuite(ctrl)
		st.doer.EXPECT().Do(gomock.Any()).Return(nil, errGeneric)

		res, err := st.client.LoginUser(testCtx, realmArjuna, clientID, username, password)

		assert.Error(t, err)
		assert.Nil(t, res)
	})

	t.Run("keycloak returns undecoded error", func(t *testing.T) {
		body := io.NopCloser(strings.NewReader("error"))

		st := createClientSuite(ctrl)
		st.doer.EXPECT().Do(gomock.Any()).Return(&http.Response{Body: body, StatusCode: http.StatusBadRequest}, nil)

		res, err := st.client.LoginUser(testCtx, realmArjuna, clientID, username, password)

		assert.Error(t, err)
		assert.Nil(t, res)
	})

	t.Run("keycloak returns bad request code", func(t *testing.T) {
		body := io.NopCloser(strings.NewReader(`{"error": "error", "error_description": "desc"}`))

		st := createClientSuite(ctrl)
		st.doer.EXPECT().Do(gomock.Any()).Return(&http.Response{Body: body, StatusCode: http.StatusBadRequest}, nil)

		res, err := st.client.LoginUser(testCtx, realmArjuna, clientID, username, password)

		assert.Error(t, err)
		assert.Nil(t, res)
	})

	t.Run("fail decode response", func(t *testing.T) {
		body := io.NopCloser(strings.NewReader("something"))

		st := createClientSuite(ctrl)
		st.doer.EXPECT().Do(gomock.Any()).Return(&http.Response{Body: body, StatusCode: http.StatusOK}, nil)

		res, err := st.client.LoginUser(testCtx, realmArjuna, clientID, username, password)

		assert.Error(t, err)
		assert.Nil(t, res)
	})

	t.Run("success login user", func(t *testing.T) {
		jwt, _ := json.Marshal(&keycloak.JWT{})
		body := io.NopCloser(bytes.NewReader(jwt))

		st := createClientSuite(ctrl)
		st.doer.EXPECT().Do(gomock.Any()).Return(&http.Response{Body: body, StatusCode: http.StatusOK}, nil)

		res, err := st.client.LoginUser(testCtx, realmArjuna, clientID, username, password)

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
		st := createClientSuite(ctrl)
		st.doer.EXPECT().Do(gomock.Any()).Return(nil, errGeneric)

		err := st.client.CreateRealm(testCtx, token, realm)

		assert.Error(t, err)
	})

	t.Run("unauthorized request", func(t *testing.T) {
		st := createClientSuite(ctrl)
		st.doer.EXPECT().Do(gomock.Any()).Return(&http.Response{StatusCode: http.StatusUnauthorized, Body: body}, nil)

		err := st.client.CreateRealm(testCtx, token, realm)

		assert.Error(t, err)
	})

	t.Run("keycloak doesn't respond with 201 status code", func(t *testing.T) {
		st := createClientSuite(ctrl)
		st.doer.EXPECT().Do(gomock.Any()).Return(&http.Response{StatusCode: http.StatusInternalServerError, Body: body}, nil)

		err := st.client.CreateRealm(testCtx, token, realm)

		assert.Error(t, err)
	})

	t.Run("successfully create a new realm", func(t *testing.T) {
		st := createClientSuite(ctrl)
		st.doer.EXPECT().Do(gomock.Any()).Return(&http.Response{StatusCode: http.StatusCreated, Body: body}, nil)

		err := st.client.CreateRealm(testCtx, token, realm)

		assert.NoError(t, err)
	})
}

func TestClient_CreateClient(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	client := &keycloak.ClientRepresentation{}
	body := io.NopCloser(bytes.NewReader([]byte(`{}`)))

	t.Run("doer returns error", func(t *testing.T) {
		st := createClientSuite(ctrl)
		st.doer.EXPECT().Do(gomock.Any()).Return(nil, errGeneric)

		err := st.client.CreateClient(testCtx, token, realmArjuna, client)

		assert.Error(t, err)
	})

	t.Run("keycloak doesn't respond with 201 status code", func(t *testing.T) {
		st := createClientSuite(ctrl)
		st.doer.EXPECT().Do(gomock.Any()).Return(&http.Response{StatusCode: http.StatusInternalServerError, Body: body}, nil)

		err := st.client.CreateClient(testCtx, token, realmArjuna, client)

		assert.Error(t, err)
	})

	t.Run("successfully create a new client", func(t *testing.T) {
		st := createClientSuite(ctrl)
		st.doer.EXPECT().Do(gomock.Any()).Return(&http.Response{StatusCode: http.StatusCreated, Body: body}, nil)

		err := st.client.CreateClient(testCtx, token, realmArjuna, client)

		assert.NoError(t, err)
	})
}

func TestClient_CreateUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	user := &keycloak.UserRepresentation{}
	body := io.NopCloser(bytes.NewReader([]byte(`{}`)))

	t.Run("doer returns error", func(t *testing.T) {
		st := createClientSuite(ctrl)
		st.doer.EXPECT().Do(gomock.Any()).Return(nil, errGeneric)

		err := st.client.CreateUser(testCtx, token, realmArjuna, user)

		assert.Error(t, err)
	})

	t.Run("keycloak with 409 status code", func(t *testing.T) {
		st := createClientSuite(ctrl)
		st.doer.EXPECT().Do(gomock.Any()).Return(&http.Response{StatusCode: http.StatusConflict, Body: body}, nil)

		err := st.client.CreateUser(testCtx, token, realmArjuna, user)

		assert.Error(t, err)
	})

	t.Run("keycloak doesn't respond with 201 status code", func(t *testing.T) {
		st := createClientSuite(ctrl)
		st.doer.EXPECT().Do(gomock.Any()).Return(&http.Response{StatusCode: http.StatusInternalServerError, Body: body}, nil)

		err := st.client.CreateUser(testCtx, token, realmArjuna, user)

		assert.Error(t, err)
	})

	t.Run("successfully create a new user", func(t *testing.T) {
		st := createClientSuite(ctrl)
		st.doer.EXPECT().Do(gomock.Any()).Return(&http.Response{StatusCode: http.StatusCreated, Body: body}, nil)

		err := st.client.CreateUser(testCtx, token, realmArjuna, user)

		assert.NoError(t, err)
	})
}

func TestClient_DeleteUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	id := "1"

	t.Run("delete api returns error", func(t *testing.T) {
		st := createClientSuite(ctrl)
		st.doer.EXPECT().Do(gomock.Any()).Return(nil, errors.New("error"))

		err := st.client.DeleteUser(testCtx, token, realmArjuna, id)

		assert.Error(t, err)
	})

	t.Run("success delete user", func(t *testing.T) {
		body := io.NopCloser(bytes.NewReader([]byte(`{}`)))
		st := createClientSuite(ctrl)
		st.doer.EXPECT().Do(gomock.Any()).Return(&http.Response{StatusCode: http.StatusNoContent, Body: body}, nil)

		err := st.client.DeleteUser(testCtx, token, realmArjuna, id)

		assert.NoError(t, err)
	})
}

func TestClient_GetUserByEmail(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	email := "arjuna@arjuna.com"

	t.Run("doer returns error", func(t *testing.T) {
		st := createClientSuite(ctrl)
		st.doer.EXPECT().Do(gomock.Any()).Return(nil, errGeneric)

		user, err := st.client.GetUserByEmail(testCtx, token, realmArjuna, email)

		assert.Error(t, err)
		assert.Nil(t, user)
	})

	t.Run("user not found", func(t *testing.T) {
		body := io.NopCloser(bytes.NewReader([]byte(`{}`)))
		st := createClientSuite(ctrl)
		st.doer.EXPECT().Do(gomock.Any()).Return(&http.Response{StatusCode: http.StatusOK, Body: body}, nil)

		user, err := st.client.GetUserByEmail(testCtx, token, realmArjuna, email)

		assert.Error(t, err)
		assert.Nil(t, user)
	})

	t.Run("success find user", func(t *testing.T) {
		body := io.NopCloser(bytes.NewReader([]byte(`[{"id": "abc", "email": "admin@arjuna.com"}]`)))
		st := createClientSuite(ctrl)
		st.doer.EXPECT().Do(gomock.Any()).Return(&http.Response{StatusCode: http.StatusOK, Body: body}, nil)

		user, err := st.client.GetUserByEmail(testCtx, token, realmArjuna, email)

		assert.NoError(t, err)
		assert.NotNil(t, user)
	})
}

func TestClient_GetAllUsers(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("get users returns error", func(t *testing.T) {
		st := createClientSuite(ctrl)
		st.doer.EXPECT().Do(gomock.Any()).Return(nil, errGeneric)

		users, err := st.client.GetAllUsers(testCtx, token, realmArjuna)

		assert.Error(t, err)
		assert.Nil(t, users)
	})

	t.Run("success get all users", func(t *testing.T) {
		body := io.NopCloser(bytes.NewReader([]byte(`[{"id": "5c44f049-8ab2-4d0f-b41d-7b08f467e817"}]`)))
		st := createClientSuite(ctrl)
		st.doer.EXPECT().Do(gomock.Any()).Return(&http.Response{StatusCode: http.StatusOK, Body: body}, nil)

		users, err := st.client.GetAllUsers(testCtx, token, realmArjuna)

		assert.NoError(t, err)
		assert.NotEmpty(t, users)
	})
}

func createClientSuite(ctrl *gomock.Controller) *ClientSuite {
	doer := mock_keycloak.NewMockDoer(ctrl)
	client := keycloak.NewClient(doer, baseURL)
	return &ClientSuite{
		client: client,
		doer:   doer,
	}
}
