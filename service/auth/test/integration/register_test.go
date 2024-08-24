//go:build integration
// +build integration

package integration

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	apiv1 "github.com/indrasaputra/arjuna/proto/api/v1"
)

var (
	testCtx      = context.Background()
	token        = base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", basicUsername, basicPassword)))
	testCtxToken = metadata.NewOutgoingContext(testCtx, metadata.Pairs("authorization", fmt.Sprintf("basic %s", token)))

	grpcURL    = "localhost:8002"
	grpcClient apiv1.AuthServiceClient
	httpClient *http.Client

	userID        = uuid.Must(uuid.NewV7())
	email         = "auth-register+1@arjuna.com"
	password      = "WeakPassword123"
	basicUsername = "auth-user"
	basicPassword = "auth-password"
)

func init() {
	setupClients()
}

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

		res, err := grpcClient.RegisterAccount(testCtxToken, req)

		assert.Error(t, err)
		assert.Equal(t, codes.InvalidArgument, status.Code(err))
		assert.Empty(t, res)
	})

	t.Run("invalid email", func(t *testing.T) {
		account := createAccount()
		account.Email = "invalid-email"
		req := &apiv1.RegisterAccountRequest{Account: account}

		res, err := grpcClient.RegisterAccount(testCtxToken, req)

		assert.Error(t, err)
		assert.Equal(t, codes.InvalidArgument, status.Code(err))
		assert.Empty(t, res)
	})

	t.Run("invalid password", func(t *testing.T) {
		account := createAccount()
		account.Password = ""
		req := &apiv1.RegisterAccountRequest{Account: account}

		res, err := grpcClient.RegisterAccount(testCtxToken, req)

		assert.Error(t, err)
		assert.Equal(t, codes.InvalidArgument, status.Code(err))
		assert.Empty(t, res)
	})

	t.Run("success register", func(t *testing.T) {
		account := createAccount()
		req := &apiv1.RegisterAccountRequest{Account: account}

		res, err := grpcClient.RegisterAccount(testCtxToken, req)

		assert.NoError(t, err)
		assert.NotEmpty(t, res)
	})
}

func setupClients() {
	conn, _ := grpc.NewClient(grpcURL, grpc.WithTransportCredentials(insecure.NewCredentials()))
	grpcClient = apiv1.NewAuthServiceClient(conn)

	httpClient = http.DefaultClient
}

func createAccount() *apiv1.Account {
	return &apiv1.Account{
		UserId:   userID.String(),
		Email:    email,
		Password: password,
	}
}
