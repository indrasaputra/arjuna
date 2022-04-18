package service

import (
	"context"
	"fmt"
	"net/mail"
	"regexp"
	"strings"

	"github.com/segmentio/ksuid"
	"golang.org/x/crypto/bcrypt"

	"github.com/indrasaputra/arjuna/service/user/entity"
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

	if err := setUserID(user); err != nil {
		return "", err
	}
	if err := hashUserPassword(user); err != nil {
		fmt.Println("fail hash password")
		return "", err
	}

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

func hashUserPassword(user *entity.User) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return entity.ErrInternal("fail to hash password")
	}
	user.Password = string(hash)
	return nil
}
