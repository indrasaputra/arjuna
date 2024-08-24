package integration

// import (
// 	"bytes"
// 	"context"
// 	"encoding/base64"
// 	"encoding/json"
// 	"fmt"
// 	"io"
// 	"log"
// 	"net/http"

// 	"github.com/google/uuid"
// 	"google.golang.org/grpc"
// 	"google.golang.org/grpc/credentials/insecure"
// 	"google.golang.org/grpc/metadata"

// 	apiv1 "github.com/indrasaputra/arjuna/proto/api/v1"
// )

// var (
// 	testCtx      = context.Background()
// 	token        = base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", basicUsername, basicPassword)))
// 	testCtxBasic = metadata.NewOutgoingContext(testCtx, metadata.Pairs("authorization", fmt.Sprintf("basic %s", token)))

// 	httpURL    = "http://localhost:8000"
// 	grpcURL    = "localhost:8003"
// 	grpcClient apiv1.TransactionCommandServiceClient
// 	httpClient *http.Client

// 	userID        = uuid.Must(uuid.NewV7())
// 	email         = "transaction-register+1@arjuna.com"
// 	password      = "WeakPassword123"
// 	basicUsername = "transaction-user"
// 	basicPassword = "transaction-password"
// )

// func init() {
// 	setupClients()
// }

// func setupClients() {
// 	conn, _ := grpc.NewClient(grpcURL, grpc.WithTransportCredentials(insecure.NewCredentials()))
// 	grpcClient = apiv1.NewTransactionCommandServiceClient(conn)

// 	httpClient = http.DefaultClient
// }

// func sendPost(url string, payload map[string]any, token string) (int, []byte) {
// 	return sendHTTPRequest(http.MethodPost, url, payload, token)
// }

// func sendHTTPRequest(method, url string, payload map[string]any, token string) (int, []byte) {
// 	var body io.Reader
// 	if payload != nil {
// 		p, err := json.Marshal(payload)
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 		body = bytes.NewReader(p)
// 	}
// 	req, err := http.NewRequestWithContext(testCtx, method, url, body)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	req.Header.Add("Content-Type", "application/json")
// 	if token != "" {
// 		req.Header.Add("Authorization", "Bearer "+token)
// 	}

// 	resp, err := httpClient.Do(req)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	b, err := io.ReadAll(resp.Body)
// 	_ = resp.Body.Close()
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	return resp.StatusCode, b
// }

// func deleteAllTransactions() {
// 	_, err := grpcClient.DeleteAllTransactions(testCtxBasic, &apiv1.DeleteAllTransactionsRequest{})
// 	if err != nil {
// 		panic(err)
// 	}
// }
