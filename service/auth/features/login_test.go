package user_test

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/cucumber/godog"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	apiv1 "github.com/indrasaputra/arjuna/proto/api/v1"
)

var (
	testCtx     = context.Background()
	authRestURL = "http://localhost:8000/v1/auth"
	userRestURL = "http://localhost:8000/v1/users"
	userGrpcURL = "localhost:8001"

	userGrpcClient apiv1.UserCommandInternalServiceClient
	httpClient     *http.Client
	httpStatus     int
	httpBody       []byte
)

type User struct {
	ID    string `json:"id"`
	Email string `json:"email"`
}

type GetAllUsersResponse struct {
	Data []*User `json:"data"`
}

func TestMain(_ *testing.M) {
	status := godog.TestSuite{
		Name:                "auth v1",
		ScenarioInitializer: InitializeScenario,
	}.Run()

	os.Exit(status)
}

func setupClients() {
	_ = godotenv.Load()
	url := os.Getenv("HTTP_AUTH_SERVER_URL")
	if url != "" {
		authRestURL = url
	}
	url = os.Getenv("HTTP_USER_SERVER_URL")
	if url != "" {
		userRestURL = url
	}
	httpClient = &http.Client{}

	url = os.Getenv("GRPC_USER_SERVER_URL")
	if url != "" {
		userGrpcURL = url
	}
	conn, _ := grpc.NewClient(userGrpcURL, grpc.WithTransportCredentials(insecure.NewCredentials()))
	userGrpcClient = apiv1.NewUserCommandInternalServiceClient(conn)
}

func restoreDefaultState(ctx context.Context, _ *godog.Scenario) (context.Context, error) {
	err := deleteAllUsers()
	checkErr(err)
	return ctx, nil
}

func cleanUpData(ctx context.Context, _ *godog.Scenario, _ error) (context.Context, error) {
	err := deleteAllUsers()
	checkErr(err)
	return ctx, nil
}

func InitializeScenario(ctx *godog.ScenarioContext) {
	setupClients()

	ctx.Before(restoreDefaultState)
	ctx.After(cleanUpData)

	ctx.Step(`^I login with user$`, iLoginWithUser)
	ctx.Step(`^response must match json$`, responseMustMatchJSON)
	ctx.Step(`^response status code must be (\d+)$`, responseStatusCodeMustBe)
	ctx.Step(`^there are users with$`, thereAreUsersWith)
}

func iLoginWithUser(requests *godog.Table) error {
	for _, row := range requests.Rows {
		body := strings.NewReader(row.Cells[0].Value)
		if err := callRestEndpoint(http.MethodPost, authRestURL+"/login", body); err != nil {
			return err
		}
	}
	return nil
}

func responseMustMatchJSON(want *godog.DocString) error {
	return deepCompareJSON([]byte(want.Content), httpBody)
}

func responseStatusCodeMustBe(code int) error {
	if httpStatus != code {
		return fmt.Errorf("expected HTTP status code %d, but got %d", code, httpStatus)
	}
	return nil
}

func thereAreUsersWith(requests *godog.Table) error {
	for _, row := range requests.Rows {
		body := strings.NewReader(row.Cells[0].Value)
		if err := callRestEndpoint(http.MethodPost, userRestURL+"/register", body); err != nil {
			return err
		}
	}
	return nil
}

func deleteAllUsers() error {
	users, err := getAllUsers()
	if err != nil {
		return err
	}

	for _, user := range users {
		req := &apiv1.DeleteUserRequest{Id: user.ID}
		if _, err := userGrpcClient.DeleteUser(testCtx, req); err != nil {
			return err
		}
	}
	return nil
}

func getAllUsers() ([]*User, error) {
	if err := callRestEndpoint(http.MethodGet, userRestURL, nil); err != nil {
		return nil, err
	}

	var resp GetAllUsersResponse
	if err := json.Unmarshal(httpBody, &resp); err != nil {
		return nil, err
	}

	return resp.Data, nil
}

func callRestEndpoint(method, url string, body io.Reader) error {
	req, err := http.NewRequestWithContext(testCtx, method, url, body)
	if err != nil {
		return err
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	httpStatus = resp.StatusCode
	httpBody, err = io.ReadAll(resp.Body)
	return err
}

func deepCompareJSON(want, have []byte) error {
	var expected interface{}
	var actual interface{}

	err := json.Unmarshal(want, &expected)
	if err != nil {
		return err
	}
	err = json.Unmarshal(have, &actual)
	if err != nil {
		return err
	}

	if !reflect.DeepEqual(expected, actual) {
		return fmt.Errorf("expected JSON does not match actual, %v vs. %v", expected, actual)
	}
	return nil
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
