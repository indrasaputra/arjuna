package keycloak

import (
	"context"
	"errors"
	"log"
	"net/http"
	"strings"

	kcsdk "github.com/indrasaputra/arjuna/pkg/sdk/keycloak"
	"github.com/indrasaputra/arjuna/service/user/entity"
	"github.com/indrasaputra/arjuna/service/user/internal/app"
)

const (
	credentialTypePassword = "password"
)

// Config defines Keycloak config.
type Config struct {
	// Client is required.
	Client kcsdk.Keycloak
	// Realm is required.
	Realm string
	// AdminUsername is required.
	AdminUsername string
	// AdminPassword is required.
	AdminPassword string
}

// Validate validates all fields in config.
// For example, all required fields must be set.
func (c *Config) Validate() error {
	if c.Client == nil {
		return errors.New("client must be set")
	}
	if strings.TrimSpace(c.Realm) == "" {
		return errors.New("realm must be set")
	}
	if strings.TrimSpace(c.AdminUsername) == "" {
		return errors.New("admin username must be set")
	}
	if strings.TrimSpace(c.AdminPassword) == "" {
		return errors.New("admin password must be set")
	}
	return nil
}

// User is responsible to connect user entity with Keycloak.
type User struct {
	config *Config
}

// NewUser creates an instance of User.
func NewUser(config *Config) (*User, error) {
	if err := config.Validate(); err != nil {
		return nil, err
	}
	return &User{config: config}, nil
}

// Create creates a new user to Keycloak.
func (u *User) Create(ctx context.Context, user *entity.User) (string, error) {
	jwt, err := u.config.Client.LoginAdmin(ctx, u.config.AdminUsername, u.config.AdminPassword)
	if err != nil {
		app.Logger.Errorf(ctx, "[KeycloakUser-Create] fail login admin: %v", err)
		return "", entity.ErrInternal(err.Error())
	}

	if err = u.createUser(ctx, user, jwt.AccessToken); err != nil {
		app.Logger.Errorf(ctx, "[KeycloakUser-Create] fail create user: %v", err)
		return "", err
	}
	// TODO: if get user somehow error, need to rollback user.
	res, err := u.config.Client.GetUserByEmail(ctx, jwt.AccessToken, u.config.Realm, user.Email)
	if err != nil {
		app.Logger.Errorf(ctx, "[KeycloakUser-Create] fail get user: %v", err)
		return "", decideError(err)
	}
	return res.ID, nil
}

// HardDelete hard-deletes user from Keycloak.
func (u *User) HardDelete(ctx context.Context, id string) error {
	jwt, err := u.config.Client.LoginAdmin(ctx, u.config.AdminUsername, u.config.AdminPassword)
	if err != nil {
		app.Logger.Errorf(ctx, "[KeycloakUser-HardDelete] fail login admin: %v", err)
		return entity.ErrInternal(err.Error())
	}

	err = u.config.Client.DeleteUser(ctx, jwt.AccessToken, u.config.Realm, id)
	if err != nil {
		app.Logger.Errorf(ctx, "[KeycloakUser-Create] fail delete user: %v", err)
		return decideError(err)
	}
	return nil
}

func (u *User) createUser(ctx context.Context, user *entity.User, accessToken string) error {
	userRep := createUserRepresentation(user)
	err := u.config.Client.CreateUser(ctx, accessToken, u.config.Realm, userRep)
	if err != nil {
		return decideError(err)
	}
	return nil
}

func getFirstAndLastName(name string) (string, string) {
	firstName := ""
	lastName := ""
	names := strings.Split(name, " ")

	firstName = names[0]
	if len(names) > 1 {
		lastName = strings.Join(names[1:], " ")
	}
	return firstName, lastName
}

func createUserRepresentation(user *entity.User) *kcsdk.UserRepresentation {
	firstName, lastName := getFirstAndLastName(user.Name)
	return &kcsdk.UserRepresentation{
		Username:  user.Username,
		FirstName: firstName,
		LastName:  lastName,
		Email:     user.Email,
		Enabled:   true,
		Credentials: []*kcsdk.CredentialRepresentation{
			{
				Type:      credentialTypePassword,
				Value:     user.Password,
				Temporary: false,
			},
		},
	}
}

func decideError(err error) error {
	kcerr, ok := err.(*kcsdk.Error)
	if !ok {
		log.Printf("[UserService-Repository] unknown error from keycloak: %v\n", err)
		return entity.ErrInternal(err.Error())
	}

	switch kcerr.Code {
	case http.StatusConflict:
		return entity.ErrAlreadyExists()
	case http.StatusNotFound:
		return entity.ErrNotFound()
	default:
		return entity.ErrInternal(err.Error())
	}
}
