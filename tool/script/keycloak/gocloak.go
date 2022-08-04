package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

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

type RealmRepresentation struct {
	ID      string `json:"id"`
	Realm   string `json:"realm"`
	Enabled bool   `json:"enabled"`
}

type ClientRepresentation struct {
	ClientID     string `json:"clientId"`
	Name         string `json:"name"`
	Enabled      bool   `json:"enabled"`
	RootURL      string `json:"rootUrl"`
	Protocol     string `json:"protocol"`
	PublicClient bool   `json:"publicClient"`
	Secret       string `json:"secret"`
}

type UserRepresentation struct {
	Username    string                      `json:"username"`
	FirstName   string                      `json:"firstName"`
	LastName    string                      `json:"lastName"`
	Email       string                      `json:"email"`
	Enabled     bool                        `json:"enabled"`
	Credentials []*CredentialRepresentation `json:"credentials"`
}

type CredentialRepresentation struct {
	Type      string `json:"type"`
	Value     string `json:"value"`
	Temporary bool   `json:"temporary"`
}

type KeycloakClient struct {
	BaseURL string
	client  *http.Client
}

func NewKeycloakClient(baseURL string) *KeycloakClient {
	client := &http.Client{
		Timeout: 1 * time.Minute,
	}

	return &KeycloakClient{
		BaseURL: strings.TrimRight(baseURL, "/"),
		client:  client,
	}
}

func (kc *KeycloakClient) LoginAdmin(ctx context.Context, username, password, realm string) (string, error) {
	url := fmt.Sprintf("%s/realms/%s/protocol/openid-connect/token", kc.BaseURL, realm)
	payload := strings.NewReader(fmt.Sprintf("client_id=admin-cli&username=%s&password=%s&grant_type=password", username, password))

	req, _ := http.NewRequestWithContext(ctx, http.MethodPost, url, payload)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	res, err := kc.client.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	var jwt *JWT
	err = json.NewDecoder(res.Body).Decode(&jwt)
	if err != nil {
		return "", err
	}
	return jwt.AccessToken, nil
}

func (kc *KeycloakClient) CreateRealm(ctx context.Context, token string, realm *RealmRepresentation) error {
	url := fmt.Sprintf("%s/admin/realms", kc.BaseURL)
	payload, _ := json.Marshal(realm)

	req, _ := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(payload))
	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Add("Content-Type", "application/json")

	res, err := kc.client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	return nil
}

func (kc *KeycloakClient) CreateClient(ctx context.Context, token string, realm string, client *ClientRepresentation) error {
	url := fmt.Sprintf("%s/admin/realms/%s/clients", kc.BaseURL, realm)
	payload, _ := json.Marshal(client)

	req, _ := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(payload))
	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Add("Content-Type", "application/json")

	res, err := kc.client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	return nil
}

func (kc *KeycloakClient) CreateUser(ctx context.Context, token string, realm string, user *UserRepresentation) error {
	url := fmt.Sprintf("%s/admin/realms/%s/users", kc.BaseURL, realm)
	payload, _ := json.Marshal(user)

	req, _ := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(payload))
	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Add("Content-Type", "application/json")

	res, err := kc.client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	return nil
}
