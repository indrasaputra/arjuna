package keycloak

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// Doer is an interface to be used as HTTP call.
type Doer interface {
	// Do does HTTP request.
	Do(*http.Request) (*http.Response, error)
}

// Error defines keycloak error.
type Error struct {
	Message string
	Code    int
}

// Error returns error message.
func (e *Error) Error() string {
	return e.Message
}

// NewError creates an instance of Error.
func NewError(code int, message string) *Error {
	return &Error{
		Code:    code,
		Message: message,
	}
}

type keycloakErrorResponse struct {
	Error       string `json:"error"`
	Description string `json:"error_description"`
}

// JWT represents JWT.
type JWT struct {
	AccessToken      string `json:"access_token"`
	IDToken          string `json:"id_token"`
	ExpiresIn        int    `json:"expires_in"`
	RefreshExpiresIn int    `json:"refresh_expires_in"`
	RefreshToken     string `json:"refresh_token"`
	TokenType        string `json:"token_type"`
	NotBeforePolicy  int    `json:"not-before-policy"`
	SessionState     string `json:"session_state"`
	Scope            string `json:"scope"`
}

// RealmRepresentation represents Keycloak realm data structure.
type RealmRepresentation struct {
	ID                  string `json:"id"`
	Realm               string `json:"realm"`
	Enabled             bool   `json:"enabled"`
	AccessTokenLifespan int    `json:"accessTokenLifespan"`
}

// ClientRepresentation represents Keycloak client data structure.
type ClientRepresentation struct {
	ClientID     string `json:"clientId"`
	Name         string `json:"name"`
	Enabled      bool   `json:"enabled"`
	RootURL      string `json:"rootUrl"`
	Protocol     string `json:"protocol"`
	PublicClient bool   `json:"publicClient"`
	Secret       string `json:"secret"`
}

// UserRepresentation represents Keycloak user data structure.
type UserRepresentation struct {
	ID          string                      `json:"id"`
	Username    string                      `json:"username"`
	FirstName   string                      `json:"firstName"`
	LastName    string                      `json:"lastName"`
	Email       string                      `json:"email"`
	Enabled     bool                        `json:"enabled"`
	Credentials []*CredentialRepresentation `json:"credentials"`
}

// CredentialRepresentation represents Keycloak credential data structure.
type CredentialRepresentation struct {
	Type      string `json:"type"`
	Value     string `json:"value"`
	Temporary bool   `json:"temporary"`
}

// Keycloak defines all use cases of keycloak.
type Keycloak interface {
	// LoginAdmin logs in as admin in Master realm.
	LoginAdmin(ctx context.Context, username, password string) (*JWT, error)
	// LoginUser logs in as user in preferred realm.
	LoginUser(ctx context.Context, realm, clientID, email, password string) (*JWT, error)
	// CreateRealm creates a realm. It needs admin's token.
	CreateRealm(ctx context.Context, token string, realm *RealmRepresentation) error
	// CreateClient creates a client. It needs admin's token.
	CreateClient(ctx context.Context, token string, realm string, client *ClientRepresentation) error
	// CreateUser creates a user. It needs admin's token.
	CreateUser(ctx context.Context, token string, realm string, user *UserRepresentation) error
	// GetUserByEmail gets a user by email. It needs admin's token.
	GetUserByEmail(ctx context.Context, token string, realm string, email string) (*UserRepresentation, error)
	// DeleteUser deletes a user. It needs admin's token.
	DeleteUser(ctx context.Context, token string, realm string, id string) error
	// GetAllUsers gets all users. It needs admin's token.
	GetAllUsers(ctx context.Context, token string, realm string) ([]*UserRepresentation, error)
}

// Client is keycloak client and responsible to communicate with Keycloak server.
// It implements Keycloak interface.
type Client struct {
	doer    Doer
	baseURL string
}

// NewClient creates an instance of Client.
func NewClient(doer Doer, baseURL string) *Client {
	return &Client{
		doer:    doer,
		baseURL: strings.TrimRight(baseURL, "/"),
	}
}

// LoginAdmin logs in as admin.
func (c *Client) LoginAdmin(ctx context.Context, username, password string) (*JWT, error) {
	url := fmt.Sprintf("%s/realms/%s/protocol/openid-connect/token", c.baseURL, "master")
	payload := strings.NewReader(fmt.Sprintf("client_id=admin-cli&username=%s&password=%s&grant_type=password", username, password))
	return c.doLogin(ctx, url, payload)
}

// LoginUser logs in user.
func (c *Client) LoginUser(ctx context.Context, realm, clientID, username, password string) (*JWT, error) {
	url := fmt.Sprintf("%s/realms/%s/protocol/openid-connect/token", c.baseURL, realm)
	payload := strings.NewReader(fmt.Sprintf("client_id=%s&username=%s&password=%s&grant_type=password", clientID, username, password))
	return c.doLogin(ctx, url, payload)
}

// CreateRealm creates a new realm in Keycloak.
func (c *Client) CreateRealm(ctx context.Context, token string, realm *RealmRepresentation) error {
	url := fmt.Sprintf("%s/admin/realms", c.baseURL)
	payload, _ := json.Marshal(realm)
	return c.doRequestWithJSON(ctx, token, http.MethodPost, url, payload, http.StatusCreated)
}

// CreateClient creates a new client in Keycloak.
func (c *Client) CreateClient(ctx context.Context, token string, realm string, client *ClientRepresentation) error {
	url := fmt.Sprintf("%s/admin/realms/%s/clients", c.baseURL, realm)
	payload, _ := json.Marshal(client)
	return c.doRequestWithJSON(ctx, token, http.MethodPost, url, payload, http.StatusCreated)
}

// CreateUser creates a new user in Keycloak.
func (c *Client) CreateUser(ctx context.Context, token string, realm string, user *UserRepresentation) error {
	url := fmt.Sprintf("%s/admin/realms/%s/users", c.baseURL, realm)
	payload, _ := json.Marshal(user)
	return c.doRequestWithJSON(ctx, token, http.MethodPost, url, payload, http.StatusCreated)
}

// GetUserByEmail gets a user by email.
func (c *Client) GetUserByEmail(ctx context.Context, token string, realm string, email string) (*UserRepresentation, error) {
	url := fmt.Sprintf("%s/admin/realms/%s/users?email=%s", c.baseURL, realm, email)
	users, err := c.doGetUsers(ctx, token, http.MethodGet, url)
	if err != nil {
		return nil, err
	}
	if len(users) == 0 {
		return nil, NewError(http.StatusNotFound, "user not found")
	}
	return users[0], nil
}

// DeleteUser deletes a user in Keycloak.
func (c *Client) DeleteUser(ctx context.Context, token string, realm string, id string) error {
	url := fmt.Sprintf("%s/admin/realms/%s/users/%s", c.baseURL, realm, id)
	return c.doRequestWithJSON(ctx, token, http.MethodDelete, url, nil, http.StatusNoContent)
}

// GetAllUsers gets all users in Keycloak.
func (c *Client) GetAllUsers(ctx context.Context, token string, realm string) ([]*UserRepresentation, error) {
	url := fmt.Sprintf("%s/admin/realms/%s/users", c.baseURL, realm)
	return c.doGetUsers(ctx, token, http.MethodGet, url)
}

func (c *Client) doLogin(ctx context.Context, url string, payload io.Reader) (*JWT, error) {
	req, _ := http.NewRequestWithContext(ctx, http.MethodPost, url, payload)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	res, err := c.doer.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = res.Body.Close() }()

	if res.StatusCode != http.StatusOK {
		return nil, decideError(res)
	}

	var jwt *JWT
	err = json.NewDecoder(res.Body).Decode(&jwt)
	if err != nil {
		return nil, err
	}
	return jwt, nil
}

func (c *Client) doRequestWithJSON(ctx context.Context, token, method, url string, payload []byte, expectedCode int) error {
	req, _ := http.NewRequestWithContext(ctx, method, url, bytes.NewBuffer(payload))
	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Add("Content-Type", "application/json")

	res, err := c.doer.Do(req)
	if err != nil {
		return err
	}
	defer func() { _ = res.Body.Close() }()

	if res.StatusCode != expectedCode {
		return decideError(res)
	}
	return nil
}

func (c *Client) doGetUsers(ctx context.Context, token, method, url string) ([]*UserRepresentation, error) {
	req, _ := http.NewRequestWithContext(ctx, method, url, nil)
	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Add("Content-Type", "application/json")

	res, err := c.doer.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = res.Body.Close() }()
	var users []*UserRepresentation
	_ = json.NewDecoder(res.Body).Decode(&users)
	return users, nil
}

func decideError(res *http.Response) error {
	var body keycloakErrorResponse
	if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
		return NewError(http.StatusInternalServerError, "problem with decoding json")
	}

	switch res.StatusCode {
	case http.StatusConflict:
		return NewError(http.StatusConflict, body.Description)
	case http.StatusUnauthorized:
		return NewError(http.StatusUnauthorized, body.Description)
	case http.StatusBadRequest:
		return NewError(http.StatusBadRequest, body.Description)
	default:
		return NewError(http.StatusInternalServerError, "internal server error")
	}
}
