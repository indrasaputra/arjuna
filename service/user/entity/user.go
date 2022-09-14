package entity

import "time"

// User defines the logical data of a user.
type User struct {
	ID         string
	KeycloakID string
	Name       string
	Username   string
	Email      string
	Password   string
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
