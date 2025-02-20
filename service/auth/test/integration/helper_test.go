//go:build integration
// +build integration

package integration

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

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

	email         = "user+1@arjuna.com" // from test/fixture/accounts.json
	password      = "password"          // from test/fixture/accounts.json
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

func sendPost(url string, payload map[string]any) (int, []byte) {
	return sendHTTPRequest(http.MethodPost, url, payload)
}

func sendHTTPRequest(method, url string, payload map[string]any) (int, []byte) {
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

	resp, err := httpClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	b, err := io.ReadAll(resp.Body)
	_ = resp.Body.Close()
	if err != nil {
		log.Fatal(err)
	}
	return resp.StatusCode, b
}
