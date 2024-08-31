package integration

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	apiv1 "github.com/indrasaputra/arjuna/proto/api/v1"
)

func TestRegister(t *testing.T) {
	t.Run("unauthenticated request", func(t *testing.T) {
		account := createAccount()
		req := &apiv1.RegisterAccountRequest{Account: account}

		res, err := grpcClient.RegisterAccount(testCtx, req)

		assert.Error(t, err)
		assert.Equal(t, codes.Unauthenticated, status.Code(err))
		assert.Empty(t, res)
	})

	t.Run("invalid user id", func(t *testing.T) {
		account := createAccount()
		account.UserId = ""
		req := &apiv1.RegisterAccountRequest{Account: account}

		res, err := grpcClient.RegisterAccount(testCtxBasic, req)

		assert.Error(t, err)
		assert.Equal(t, codes.InvalidArgument, status.Code(err))
		assert.Empty(t, res)
	})

	t.Run("invalid email", func(t *testing.T) {
		account := createAccount()
		account.Email = "invalid-email"
		req := &apiv1.RegisterAccountRequest{Account: account}

		res, err := grpcClient.RegisterAccount(testCtxBasic, req)

		assert.Error(t, err)
		assert.Equal(t, codes.InvalidArgument, status.Code(err))
		assert.Empty(t, res)
	})

	t.Run("invalid password", func(t *testing.T) {
		account := createAccount()
		account.Password = ""
		req := &apiv1.RegisterAccountRequest{Account: account}

		res, err := grpcClient.RegisterAccount(testCtxBasic, req)

		assert.Error(t, err)
		assert.Equal(t, codes.InvalidArgument, status.Code(err))
		assert.Empty(t, res)
	})

	t.Run("success register", func(t *testing.T) {
		account := createAccount()
		req := &apiv1.RegisterAccountRequest{Account: account}

		res, err := grpcClient.RegisterAccount(testCtxBasic, req)

		assert.NoError(t, err)
		assert.NotEmpty(t, res)
	})
}

func createAccount() *apiv1.Account {
	return &apiv1.Account{
		UserId:   uuid.Must(uuid.NewV7()).String(),
		Email:    "auth-register+1@arjuna.com",
		Password: "password",
	}
}
