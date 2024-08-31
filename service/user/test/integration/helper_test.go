package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	apiv1 "github.com/indrasaputra/arjuna/proto/api/v1"
)

var (
	testCtx = context.Background()

	httpURL    = "http://localhost:8000"
	grpcURL    = "localhost:8001"
	grpcClient apiv1.UserCommandServiceClient
	httpClient *http.Client
	path       = "/v1/users/register"

	email    = "user-register+1@arjuna.com"
	password = "password"
	name     = "User Register First"
)

func init() {
	setupClients()
}

func setupClients() {
	conn, _ := grpc.NewClient(grpcURL, grpc.WithTransportCredentials(insecure.NewCredentials()))
	grpcClient = apiv1.NewUserCommandServiceClient(conn)

	httpClient = http.DefaultClient
}

func sendPost(url string, payload map[string]any, key string) (int, []byte) {
	return sendHTTPRequest(http.MethodPost, url, payload, key)
}

func sendHTTPRequest(method, url string, payload map[string]any, key string) (int, []byte) {
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
	if key != "" {
		req.Header.Add("X-Idempotency-Key", key)
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
	return resp.StatusCode, b
}
