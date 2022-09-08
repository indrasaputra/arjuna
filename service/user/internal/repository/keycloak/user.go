package keycloak

import (
	"context"
	"errors"
	"strings"

	kcsdk "github.com/indrasaputra/arjuna/pkg/sdk/keycloak"
	"github.com/indrasaputra/arjuna/service/user/entity"
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

// Insert inserts a new user to Keycloak.
func (u *User) Insert(ctx context.Context, user *entity.User) error {
	jwt, err := u.config.Client.LoginAdmin(ctx, u.config.AdminUsername, u.config.AdminPassword)
	if err != nil {
		return entity.ErrInternal(err.Error())
	}

	userRep := createUserRepresentation(user)
	if err = u.config.Client.CreateUser(ctx, jwt.AccessToken, u.config.Realm, userRep); err != nil {
		return entity.ErrInternal(err.Error())
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
		FirstName: firstName,
		LastName:  lastName,
		Email:     user.Email,
		Enabled:   true,
		Credentials: []*kcsdk.CredentialRepresentation{
			{
				Type:      "password",
				Value:     user.Password,
				Temporary: false,
			},
		},
	}
}
