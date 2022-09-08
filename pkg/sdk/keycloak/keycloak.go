package keycloak

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// Doer is an interface to be used as HTTP call.
type Doer interface {
	// Do does HTTP request.
	Do(*http.Request) (*http.Response, error)
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
	ID      string `json:"id"`
	Realm   string `json:"realm"`
	Enabled bool   `json:"enabled"`
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
	// CreateRealm creates a realm. It needs admin's token.
	CreateRealm(ctx context.Context, token string, realm *RealmRepresentation) error
	// CreateClient creates a client. It needs admin's token.
	CreateClient(ctx context.Context, token string, realm string, client *ClientRepresentation) error
	// CreateUser creates a user. It needs admin's token.
	CreateUser(ctx context.Context, token string, realm string, user *UserRepresentation) error
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

	req, _ := http.NewRequestWithContext(ctx, http.MethodPost, url, payload)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	res, err := c.doer.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = res.Body.Close() }()

	var jwt *JWT
	err = json.NewDecoder(res.Body).Decode(&jwt)
	if err != nil {
		return nil, err
	}
	return jwt, nil
}

// CreateRealm creates a new realm in Keycloak.
func (c *Client) CreateRealm(ctx context.Context, token string, realm *RealmRepresentation) error {
	url := fmt.Sprintf("%s/admin/realms", c.baseURL)
	payload, _ := json.Marshal(realm)
	return c.doRequestWithJSON(ctx, token, http.MethodPost, url, payload)
}

// CreateClient creates a new client in Keycloak.
func (c *Client) CreateClient(ctx context.Context, token string, realm string, client *ClientRepresentation) error {
	url := fmt.Sprintf("%s/admin/realms/%s/clients", c.baseURL, realm)
	payload, _ := json.Marshal(client)
	return c.doRequestWithJSON(ctx, token, http.MethodPost, url, payload)
}

// CreateUser creates a new user in Keycloak.
func (c *Client) CreateUser(ctx context.Context, token string, realm string, user *UserRepresentation) error {
	url := fmt.Sprintf("%s/admin/realms/%s/users", c.baseURL, realm)
	payload, _ := json.Marshal(user)
	return c.doRequestWithJSON(ctx, token, http.MethodPost, url, payload)
}

func (c *Client) doRequestWithJSON(ctx context.Context, token, method, url string, payload []byte) error {
	req, _ := http.NewRequestWithContext(ctx, method, url, bytes.NewBuffer(payload))
	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Add("Content-Type", "application/json")

	res, err := c.doer.Do(req)
	if err != nil {
		return err
	}
	defer func() { _ = res.Body.Close() }()

	if res.StatusCode != http.StatusCreated {
		return fmt.Errorf("got HTTP status code: %d", res.StatusCode)
	}
	return nil
}
