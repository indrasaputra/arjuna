package entity

import "time"

// Token represents token.
type Token struct {
	AccessToken           string
	TokenType             string
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

// Auditable defines logical data related to audit.
type Auditable struct {
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
	CreatedBy string
	UpdatedBy string
	DeletedBy string
}
