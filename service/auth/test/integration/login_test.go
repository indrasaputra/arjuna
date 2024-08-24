package integration

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	apiv1 "github.com/indrasaputra/arjuna/proto/api/v1"
)

func TestLogin(t *testing.T) {
	loginEmail := "auth-login+1@arjuna.com"
	registerAccount(loginEmail)

	t.Run("invalid email", func(t *testing.T) {
		payload := map[string]any{"email": "", "password": password}

		status, resp := sendPost(httpURL+"/v1/auth/login", payload, "")

		assert.Equal(t, http.StatusBadRequest, status)
		assert.Equal(t, float64(3), resp["code"])
	})

	t.Run("invalid password", func(t *testing.T) {
		payload := map[string]any{"email": loginEmail, "password": ""}

		status, resp := sendPost(httpURL+"/v1/auth/login", payload, "")

		assert.Equal(t, http.StatusBadRequest, status)
		assert.Equal(t, float64(3), resp["code"])
	})

	t.Run("success login", func(t *testing.T) {
		payload := map[string]any{"email": loginEmail, "password": password}

		status, resp := sendPost(httpURL+"/v1/auth/login", payload, "")

		assert.Equal(t, http.StatusBadRequest, status)
		assert.NotEmpty(t, resp["access_token"])
		assert.NotEmpty(t, resp["access_token_expires_in"])
	})
}

func sendPost(url string, payload map[string]any, token string) (int, map[string]any) {
	return sendHTTPRequest(http.MethodPost, url, payload, token)
}

func sendHTTPRequest(method, url string, payload map[string]any, token string) (int, map[string]any) {
	var body io.Reader
	if payload != nil {
		p, err := json.Marshal(payload)
		if err != nil {
			log.Fatal(err)
		}
		body = bytes.NewReader(p)
	}
	req, err := http.NewRequestWithContext(testCtx, method, url, body)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Add("Content-Type", "application/json")
	if token != "" {
		req.Header.Add("Authorization", "Bearer "+token)
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	b, err := io.ReadAll(resp.Body)
	_ = resp.Body.Close()
	if err != nil {
		log.Fatal(err)
	}

	var v map[string]any
	_ = json.Unmarshal(b, &v)

	return resp.StatusCode, v
}

func registerAccount(email string) {
	account := createAccount()
	account.Email = email
	req := &apiv1.RegisterAccountRequest{Account: account}

	_, _ = grpcClient.RegisterAccount(testCtxBasic, req)
}
