// Server main program.
package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/indrasaputra/arjuna/pkg/sdk/keycloak"
)

const (
	keycloakAddress       = "http://localhost:8080/"
	protocolOpenidConnect = "openid-connect"
	masterAdminUsername   = "admin"
	masterAdminPassword   = "admin"
	arjunaRealm           = "arjuna"
	arjunaClientID        = "arjuna-client"
	arjunaClientName      = "Arjuna Client"
	arjunaUserUsername    = "arjunauser"
	arjunaUserPassword    = "arjunapassword"
	arjunaUserFirstName   = "First"
	arjunaUserLastName    = "User"
	arjunaUserEmail       = "user.first@arjuna.com"

	thirtyMinutes = 1800
)

func main() {
	ctx := context.Background()
	httpClient := &http.Client{Timeout: time.Minute}
	keycloakClient := keycloak.NewClient(httpClient, keycloakAddress)

	token, err := keycloakClient.LoginAdmin(ctx, masterAdminUsername, masterAdminPassword)
	if err != nil {
		log.Fatalf("fail login admin: %v\n", err)
	}

	if err := CreateArjunaRealm(ctx, keycloakClient, token.AccessToken); err != nil {
		log.Fatalf("fail create realm: %v\n", err)
	}

	if err := CreateArjunaClient(ctx, keycloakClient, token.AccessToken); err != nil {
		log.Fatalf("fail create client: %v\n", err)
	}

	if err := CreateArjunaUser(ctx, keycloakClient, token.AccessToken); err != nil {
		log.Fatalf("fail create user: %v\n", err)
	}
}

// CreateArjunaRealm creates realm arjuna.
func CreateArjunaRealm(ctx context.Context, client *keycloak.Client, token string) error {
	realm := &keycloak.RealmRepresentation{
		ID:                  arjunaRealm,
		Realm:               arjunaRealm,
		Enabled:             true,
		AccessTokenLifespan: thirtyMinutes,
	}
	return client.CreateRealm(ctx, token, realm)
}

// CreateArjunaClient creates client arjuna.
func CreateArjunaClient(ctx context.Context, client *keycloak.Client, token string) error {
	cl := &keycloak.ClientRepresentation{
		ClientID:     arjunaClientID,
		Name:         arjunaClientName,
		Enabled:      true,
		Protocol:     protocolOpenidConnect,
		PublicClient: true,
	}
	return client.CreateClient(ctx, token, arjunaRealm, cl)
}

// CreateArjunaUser creates user arjuna.
func CreateArjunaUser(ctx context.Context, client *keycloak.Client, token string) error {
	user := &keycloak.UserRepresentation{
		Username:  arjunaUserUsername,
		FirstName: arjunaUserFirstName,
		LastName:  arjunaUserLastName,
		Email:     arjunaUserEmail,
		Enabled:   true,
		Credentials: []*keycloak.CredentialRepresentation{
			{
				Type:      "password",
				Value:     arjunaUserPassword,
				Temporary: false,
			},
		},
	}
	return client.CreateUser(ctx, token, arjunaRealm, user)
}
