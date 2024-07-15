package entity

import "time"

// User defines the logical data of a user.
type User struct {
	ID       string
	Name     string
	Email    string
	Password string
	Auditable
}

// UserOutboxStatus enumerates user outbox status.
type UserOutboxStatus string

var (
	// UserOutboxStatusReady means ready to be picked up.
	UserOutboxStatusReady UserOutboxStatus = "READY"
	// UserOutboxStatusProcessed means being processed.
	UserOutboxStatusProcessed UserOutboxStatus = "PROCESSED"
	// UserOutboxStatusDelivered means successfully sent to server.
	UserOutboxStatusDelivered UserOutboxStatus = "DELIVERED"
	// UserOutboxStatusFailed means failure.
	UserOutboxStatusFailed UserOutboxStatus = "FAILED"
)

// UserOutbox defines logical data of user outbox.
type UserOutbox struct {
	ID      string
	Status  UserOutboxStatus
	Payload *User
	Auditable
}

// RegisterUserInput holds input data for register user workflow.
type RegisterUserInput struct {
	User *User
}

// RegisterUserOutput holds output data for register user workflow.
type RegisterUserOutput struct {
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
