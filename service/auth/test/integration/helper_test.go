//go:build integration
// +build integration

package integration

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"

	apiv1 "github.com/indrasaputra/arjuna/proto/api/v1"
)

var (
	testCtx      = context.Background()
	token        = base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", basicUsername, basicPassword)))
	testCtxBasic = metadata.NewOutgoingContext(testCtx, metadata.Pairs("authorization", fmt.Sprintf("basic %s", token)))

	httpURL    = "http://localhost:8000"
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

func setupClients() {
	conn, _ := grpc.NewClient(grpcURL, grpc.WithTransportCredentials(insecure.NewCredentials()))
	grpcClient = apiv1.NewAuthServiceClient(conn)

	httpClient = http.DefaultClient
}
