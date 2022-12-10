package service

import (
	"context"
	"crypto/rand"
	"math/big"
	"net/mail"
	"regexp"
	"strings"
	"time"

	"github.com/segmentio/ksuid"

	"github.com/indrasaputra/arjuna/service/user/entity"
)

const (
	allowedUsernameCharacters = "abcdefghijklmnopqrstuvwxyz0123456789"
	usernameLength            = 10
)

var (
	regexNameCompiler *regexp.Regexp
)

func init() {
	regexNameCompiler = regexp.MustCompile(`^[a-zA-Z\s]+$`)
}

// RegisterUser defines the interface to register a user.
type RegisterUser interface {
	// Register registers a user and store it in the storage.
	// It returns the ID of the newly created user.
	// It must check the uniqueness of the user.
	Register(ctx context.Context, user *entity.User) (string, error)
}

// RegisterUserWorkflow defines the interface for user registration workflow.
type RegisterUserWorkflow interface {
	// RegisterUser registers the user into the repository or 3rd party needed.
	// It also validates if the user's email is unique.
	// It returns the ID of the created user.
	RegisterUser(ctx context.Context, input *RegisterUserInput) (*RegisterUserOutput, error)
}

// RegisterUserInput holds input data for register user workflow
type RegisterUserInput struct {
	User *entity.User
}

// RegisterUserOutput holds output data for register user workflow
type RegisterUserOutput struct {
	UserID string
}

// UserRegistrar is responsible for registering a new user.
type UserRegistrar struct {
	workflow RegisterUserWorkflow
}

// NewUserRegistrar creates an instance of UserRegistrar.
func NewUserRegistrar(workflow RegisterUserWorkflow) *UserRegistrar {
	return &UserRegistrar{workflow: workflow}
}

// Register registers a user and store it in the storage.
// It returns the ID of the newly created user.
// It checks the email for duplication.
func (ur *UserRegistrar) Register(ctx context.Context, user *entity.User) (string, error) {
	if err := validateUser(user); err != nil {
		return "", err
	}
	sanitizeUser(user)

	// username is mandatory for Keycloak, but not for this current business.
	// hence, generating a random username is fine.
	user.Username = generateUsername(usernameLength)

	if err := setUserID(user); err != nil {
		return "", err
	}
	setUserAuditableProperties(user)

	input := &RegisterUserInput{User: user}
	output, err := ur.workflow.RegisterUser(ctx, input)
	if err != nil {
		return "", err
	}
	return output.UserID, nil
}

func validateUser(user *entity.User) error {
	if user == nil {
		return entity.ErrEmptyUser()
	}
	if !regexNameCompiler.MatchString(user.Name) {
		return entity.ErrInvalidName()
	}
	if _, err := mail.ParseAddress(user.Email); err != nil {
		return entity.ErrInvalidEmail()
	}
	return nil
}

func sanitizeUser(user *entity.User) {
	user.Name = strings.TrimSpace(user.Name)
	user.Email = strings.TrimSpace(user.Email)
}

func setUserID(user *entity.User) error {
	id, err := ksuid.NewRandom()
	if err != nil {
		return entity.ErrInternal("fail to create user's ID")
	}
	user.ID = id.String()
	return nil
}

func setUserAuditableProperties(user *entity.User) {
	user.CreatedAt = time.Now().UTC()
	user.UpdatedAt = time.Now().UTC()
	user.CreatedBy = user.ID
	user.UpdatedBy = user.ID
}

func generateUsername(n int) string {
	username := ""
	length := len(allowedUsernameCharacters)
	for i := 0; i < n; i++ {
		username += string(allowedUsernameCharacters[cryptoRandSecure(int64(length))])
	}
	return username
}

func cryptoRandSecure(max int64) int64 {
	nBig, _ := rand.Int(rand.Reader, big.NewInt(max))
	return nBig.Int64()
}
