package entity

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
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
	Email    string `json:"email"`
	Password string `json:"password"`
	Auditable
	ID     uuid.UUID `json:"id"`
	UserID uuid.UUID `json:"user_id"`
}

// Claims represents token claims.
type Claims struct {
	jwt.RegisteredClaims
	Email     string    `json:"email"`
	AccountID uuid.UUID `json:"account_id"`
	UserID    uuid.UUID `json:"user_id"`
}

// Auditable defines logical data related to audit.
type Auditable struct {
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
	DeletedBy *uuid.UUID
	CreatedBy uuid.UUID
	UpdatedBy uuid.UUID
}
