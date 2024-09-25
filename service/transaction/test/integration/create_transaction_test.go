//go:build integration
// +build integration

package integration

// import (
// 	"net/http"
// 	"testing"

// 	"github.com/stretchr/testify/assert"
// 	"github.com/tidwall/gjson"

// 	apiv1 "github.com/indrasaputra/arjuna/proto/api/v1"
// )

// func TestCreateTransaction(t *testing.T) {
// 	deleteAllTransactions()
// 	registerAccount(email)

// 	t.Run("invalid email", func(t *testing.T) {
// 		payload := map[string]any{"email": "", "password": password}

// 		status, resp := sendPost(httpURL+"/v1/auth/login", payload, "")

// 		assert.Equal(t, http.StatusBadRequest, status)
// 		assert.Equal(t, float64(3), gjson.GetBytes(resp, "code").Float())
// 	})

// 	t.Run("invalid password", func(t *testing.T) {
// 		payload := map[string]any{"email": loginEmail, "password": ""}

// 		status, resp := sendPost(httpURL+"/v1/auth/login", payload, "")

// 		assert.Equal(t, http.StatusBadRequest, status)
// 		assert.Equal(t, float64(3), gjson.GetBytes(resp, "code").Float())
// 	})

// 	t.Run("success login", func(t *testing.T) {
// 		payload := map[string]any{"email": loginEmail, "password": password}

// 		status, resp := sendPost(httpURL+"/v1/auth/login", payload, "")

// 		assert.Equal(t, http.StatusOK, status)
// 		assert.NotEmpty(t, gjson.GetBytes(resp, "data.access_token").String())
// 		assert.NotEmpty(t, gjson.GetBytes(resp, "data.access_token_expires_in").String())
// 	})
// }

// func registerAccount(email string) {
// 	account := createAccount()
// 	account.Email = email
// 	req := &apiv1.RegisterAccountRequest{Account: account}

// 	_, _ = grpcClient.RegisterAccount(testCtxBasic, req)
// }
