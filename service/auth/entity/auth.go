package entity

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Token represents token.
type Token struct {
	AccessToken           string
	RefreshToken          string
	AccessTokenExpiresIn  uint32
	RefreshTokenExpiresIn uint32
}

// Account represents account.
type Account struct {
	ID       string
	UserID   string
	Email    string
	Password string
	Auditable
}

// Claims represents token claims.
type Claims struct {
	AccountID string `json:"account_id"`
	UserID    string `json:"user_id"`
	Email     string `json:"email"`
	jwt.RegisteredClaims
}

// Auditable defines logical data related to audit.
type Auditable struct {
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
	CreatedBy string
	UpdatedBy string
	DeletedBy string
}
