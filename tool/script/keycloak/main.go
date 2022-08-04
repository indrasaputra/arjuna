package main

import (
	"context"
	"log"
)

const (
	KeycloakAddress       = "http://localhost:8080/"
	ProtocolOpenidConnect = "openid-connect"
	MasterAdminUsername   = "admin"
	MasterAdminPassword   = "admin"
	MasterRealm           = "master"
	ArjunaRealm           = "arjuna"
	ArjunaClientID        = "arjuna-client"
	ArjunaClientName      = "Arjuna Client"
	ArjunaClientSecret    = "arjuna-secret"
	ArjunaClientRootURL   = "https://www.keycloak.org/app/"
	ArjunaUserUsername    = "arjunauser"
	ArjunaUserPassword    = "arjunapassword"
	ArjunaUserFirstName   = "First"
	ArjunaUserLastName    = "User"
	ArjunaUserEmail       = "user.first@arjuna.com"
)

func main() {
	ctx := context.Background()
	client := NewKeycloakClient(KeycloakAddress)

	token, err := client.LoginAdmin(ctx, MasterAdminUsername, MasterAdminPassword, MasterRealm)
	if err != nil {
		log.Fatalf("fail login admin: %v\n", err)
	}

	if err := CreateArjunaRealm(ctx, client, token); err != nil {
		log.Fatalf("fail create realm: %v\n", err)
	}

	if err := CreateArjunaClient(ctx, client, token); err != nil {
		log.Fatalf("fail create client: %v\n", err)
	}

	if err := CreateArjunaUser(ctx, client, token); err != nil {
		log.Fatalf("fail create user: %v\n", err)
	}
}

func CreateArjunaRealm(ctx context.Context, client *KeycloakClient, token string) error {
	realm := &RealmRepresentation{
		ID:      ArjunaRealm,
		Realm:   ArjunaRealm,
		Enabled: true,
	}
	return client.CreateRealm(ctx, token, realm)
}

func CreateArjunaClient(ctx context.Context, client *KeycloakClient, token string) error {
	cl := &ClientRepresentation{
		ClientID:     ArjunaClientID,
		Name:         ArjunaClientName,
		Enabled:      true,
		RootURL:      ArjunaClientRootURL,
		Protocol:     ProtocolOpenidConnect,
		PublicClient: false,
		Secret:       ArjunaClientSecret,
	}
	return client.CreateClient(ctx, token, ArjunaRealm, cl)
}

func CreateArjunaUser(ctx context.Context, client *KeycloakClient, token string) error {
	user := &UserRepresentation{
		Username:  ArjunaUserUsername,
		FirstName: ArjunaUserFirstName,
		LastName:  ArjunaUserLastName,
		Email:     ArjunaUserEmail,
		Enabled:   true,
		Credentials: []*CredentialRepresentation{
			{
				Type:      "password",
				Value:     ArjunaUserPassword,
				Temporary: false,
			},
		},
	}
	return client.CreateUser(ctx, token, ArjunaRealm, user)
}
