package integration

import (
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/tidwall/gjson"

	"github.com/indrasaputra/arjuna/service/user/entity"
)

func TestRegisterUser(t *testing.T) {
	t.Run("empty idemptency key", func(t *testing.T) {
		user := createUser()
		payload := map[string]any{"email": user.Email, "password": user.Password, "name": user.Name}

		status, resp := sendPost(httpURL+path, payload, "")

		assert.Equal(t, http.StatusBadRequest, status)
		assert.Equal(t, float64(3), gjson.GetBytes(resp, "code").Float())
		assert.Equal(t, "missing idempotency key", gjson.GetBytes(resp, "message").String())
		assert.Equal(t, "USER_ERROR_CODE_MISSING_IDEMPOTENCY_KEY", gjson.GetBytes(resp, "details.0.errorCode").String())
	})

	t.Run("idempotency key has been used", func(t *testing.T) {
		key := uuid.Must(uuid.NewV7()).String()
		user := createUser()
		payload := map[string]any{"email": "email", "password": user.Password, "name": user.Name}

		_, _ = sendPost(httpURL+path, payload, key)
		status, resp := sendPost(httpURL+path, payload, key)

		assert.Equal(t, http.StatusConflict, status)
		assert.Equal(t, float64(6), gjson.GetBytes(resp, "code").Float())
		assert.Equal(t, "USER_ERROR_CODE_ALREADY_EXISTS", gjson.GetBytes(resp, "details.0.errorCode").String())
	})

	t.Run("name contains character outside of alphabet", func(t *testing.T) {
		key := uuid.Must(uuid.NewV7()).String()
		user := createUser()
		payload := map[string]any{"email": user.Email, "password": user.Password, "name": "4rjun4"}

		status, resp := sendPost(httpURL+path, payload, key)

		assert.Equal(t, http.StatusBadRequest, status)
		assert.Equal(t, float64(3), gjson.GetBytes(resp, "code").Float())
		assert.Equal(t, "contain character outside of alphabet", gjson.GetBytes(resp, "details.0.fieldViolations.0.description").String())
		assert.Equal(t, "USER_ERROR_CODE_INVALID_NAME", gjson.GetBytes(resp, "details.1.errorCode").String())
	})

	t.Run("invalid email", func(t *testing.T) {
		key := uuid.Must(uuid.NewV7()).String()
		user := createUser()
		payload := map[string]any{"email": "invalid-email", "password": user.Password, "name": user.Name}

		status, resp := sendPost(httpURL+path, payload, key)

		assert.Equal(t, http.StatusBadRequest, status)
		assert.Equal(t, float64(3), gjson.GetBytes(resp, "code").Float())
		assert.Equal(t, "USER_ERROR_CODE_INVALID_EMAIL", gjson.GetBytes(resp, "details.1.errorCode").String())
	})

	t.Run("success register", func(t *testing.T) {
		key := uuid.Must(uuid.NewV7()).String()
		user := createUser()
		payload := map[string]any{"email": user.Email, "password": user.Password, "name": user.Name}

		status, resp := sendPost(httpURL+path, payload, key)

		assert.Equal(t, http.StatusOK, status)
		id := gjson.GetBytes(resp, "data.id").String()
		_, err := uuid.Parse(id)
		assert.NoError(t, err)
	})
}

func createUser() *entity.User {
	id, _ := uuid.NewV7()
	return &entity.User{
		ID:       id,
		Email:    email,
		Password: password,
		Name:     name,
	}
}
