package toggle_test

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/cucumber/godog"
	"github.com/joho/godotenv"
)

var (
	ctx         = context.Background()
	client      = http.DefaultClient
	userURLHTTP = "http://localhost:8000/v1/users"

	httpStatus int
	httpBody   []byte
)

func TestMain(_ *testing.M) {
	status := godog.TestSuite{
		Name:                "user v1",
		ScenarioInitializer: InitializeScenario,
	}.Run()

	os.Exit(status)
}

func restoreDefaultState(ctx context.Context, sc *godog.Scenario) (context.Context, error) {
	return ctx, nil
}

func cleanUpData(ctx context.Context, sc *godog.Scenario, _ error) (context.Context, error) {
	return ctx, nil
}

func InitializeScenario(ctx *godog.ScenarioContext) {
	_ = godotenv.Load()
	url := os.Getenv("SERVER_URL")
	if url != "" {
		userURLHTTP = url
	}

	ctx.Before(restoreDefaultState)
	ctx.After(cleanUpData)

	ctx.Step(`^the user is empty$`, theUserIsEmpty)
	ctx.Step(`^I register user with body$`, iRegisterUserWithBody)
	ctx.Step(`^response status code must be (\d+)$`, responseStatusCodeMustBe)
	ctx.Step(`^response must match json$`, responseMustMatchJSON)
}

func theUserIsEmpty() error {
	return nil
}

func iRegisterUserWithBody(requests *godog.Table) error {
	for _, row := range requests.Rows {
		body := strings.NewReader(row.Cells[0].Value)
		if err := callEndpoint(http.MethodPost, userURLHTTP+"/register", body); err != nil {
			return err
		}
	}
	return nil
}

func responseStatusCodeMustBe(code int) error {
	if httpStatus != code {
		return fmt.Errorf("expected HTTP status code %d, but got %d", code, httpStatus)
	}
	return nil
}

func responseMustMatchJSON(want *godog.DocString) error {
	return deepCompareJSON([]byte(want.Content), httpBody)
}

// func deleteAllUsers() error {
// 	toggles, err := getAllUsers()
// 	if err != nil {
// 		return err
// 	}

// 	for _, toggle := range toggles {
// 		if err = callEndpoint(http.MethodPut, fmt.Sprintf("%s/%s/disable", toggleURL, toggle.Key), nil); err != nil {
// 			return err
// 		}
// 		if err = callEndpoint(http.MethodDelete, fmt.Sprintf("%s/%s", toggleURL, toggle.Key), nil); err != nil {
// 			return err
// 		}
// 	}
// 	return nil
// }

// func getAllUsers() {
// 	if err := callEndpoint(http.MethodGet, toggleURL, nil); err != nil {
// 		return nil, err
// 	}

// 	var resp GetAllResponse
// 	if err := json.Unmarshal(httpBody, &resp); err != nil {
// 		return nil, err
// 	}

// 	return resp.Toggles, nil
// }

func callEndpoint(method, url string, body io.Reader) error {
	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return err
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	httpStatus = resp.StatusCode
	httpBody, err = ioutil.ReadAll(resp.Body)
	log.Println(string(httpBody))
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
