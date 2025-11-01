//go:build integration
// +build integration

package integration

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tidwall/gjson"

	apiv1 "github.com/indrasaputra/arjuna/proto/api/v1"
)

func TestLogin(t *testing.T) {
	// register account for login
	account := createAccount()
	req := &apiv1.RegisterAccountRequest{Account: account}

	_, _ = grpcClient.RegisterAccount(testCtxBasic, req)

	t.Run("empty email", func(t *testing.T) {
		payload := map[string]any{"email": "", "password": account.Password}

		status, resp := sendPost(httpURL+"/v1/auth/login", payload)

		assert.Equal(t, http.StatusBadRequest, status)
		assert.Equal(t, float64(3), gjson.GetBytes(resp, "code").Float())
		assert.Equal(t, "email", gjson.GetBytes(resp, "details.0.fieldViolations.0.field").String())
		assert.Equal(t, "empty or nil", gjson.GetBytes(resp, "details.0.fieldViolations.0.description").String())
		assert.Equal(t, "AUTH_ERROR_CODE_EMPTY_FIELD", gjson.GetBytes(resp, "details.1.errorCode").String())
	})

	t.Run("empty password", func(t *testing.T) {
		payload := map[string]any{"email": account.Email, "password": ""}

		status, resp := sendPost(httpURL+"/v1/auth/login", payload)

		assert.Equal(t, http.StatusBadRequest, status)
		assert.Equal(t, float64(3), gjson.GetBytes(resp, "code").Float())
		assert.Equal(t, "password", gjson.GetBytes(resp, "details.0.fieldViolations.0.field").String())
		assert.Equal(t, "empty or nil", gjson.GetBytes(resp, "details.0.fieldViolations.0.description").String())
		assert.Equal(t, "AUTH_ERROR_CODE_EMPTY_FIELD", gjson.GetBytes(resp, "details.1.errorCode").String())
	})

	t.Run("invalid email", func(t *testing.T) {
		payload := map[string]any{"email": "invalid-email", "password": "password"}

		status, resp := sendPost(httpURL+"/v1/auth/login", payload)

		assert.Equal(t, http.StatusBadRequest, status)
		assert.Equal(t, float64(3), gjson.GetBytes(resp, "code").Float())
		assert.Equal(t, "credential is invalid", gjson.GetBytes(resp, "message").String())
		assert.Equal(t, "AUTH_ERROR_CODE_INVALID_CREDENTIAL", gjson.GetBytes(resp, "details.0.errorCode").String())
	})

	t.Run("success login", func(t *testing.T) {
		payload := map[string]any{"email": account.Email, "password": account.Password}

		status, resp := sendPost(httpURL+"/v1/auth/login", payload)

		assert.Equal(t, http.StatusOK, status)
		assert.NotEmpty(t, gjson.GetBytes(resp, "data.access_token").String())
		assert.NotEmpty(t, gjson.GetBytes(resp, "data.access_token_expires_in").String())
	})

	t.Run("invalid password", func(t *testing.T) {
		payload := map[string]any{"email": account.Email, "password": "not-the-right-password"}

		status, resp := sendPost(httpURL+"/v1/auth/login", payload)

		assert.Equal(t, http.StatusBadRequest, status)
		assert.Equal(t, float64(3), gjson.GetBytes(resp, "code").Float())
		assert.Equal(t, "credential is invalid", gjson.GetBytes(resp, "message").String())
		assert.Equal(t, "AUTH_ERROR_CODE_INVALID_CREDENTIAL", gjson.GetBytes(resp, "details.0.errorCode").String())
	})
}
