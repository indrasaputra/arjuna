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
	//
	// It must check the uniqueness of the user.
	Register(ctx context.Context, user *entity.User) (string, error)
}

// RegisterUserRepository defines the interface to save a user into the repository.
type RegisterUserRepository interface {
	// Insert inserts the user into the repository.
	// It also validates if the user's email is unique.
	// It returns the ID of the created user.
	Insert(ctx context.Context, user *entity.User) error
}

// UserRegistrator is responsible for registering a new user.
type UserRegistrator struct {
	repo RegisterUserRepository
}

// NewUserRegistrator creates an instance of UserRegistrator.
func NewUserRegistrator(repo RegisterUserRepository) *UserRegistrator {
	return &UserRegistrator{repo: repo}
}

// Register registers a user and store it in the storage.
// It returns the ID of the newly created user.
//
// It checks the email for duplication.
func (ur *UserRegistrator) Register(ctx context.Context, user *entity.User) (string, error) {
	if err := validateUser(user); err != nil {
		return "", err
	}
	sanitizeUser(user)

	// username is mandatory for Keycloak, but not for business.
	// hence, generating a random username is fine.
	user.Username = generateUsername(usernameLength)

	if err := setUserID(user); err != nil {
		return "", err
	}
	setUserAuditableProperties(user)

	if err := ur.repo.Insert(ctx, user); err != nil {
		return "", err
	}
	return user.ID, nil
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
