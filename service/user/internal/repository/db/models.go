// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package db

import (
	"database/sql/driver"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/indrasaputra/arjuna/service/user/entity"
)

type UserOutboxStatus string

const (
	UserOutboxStatusREADY     UserOutboxStatus = "READY"
	UserOutboxStatusPROCESSED UserOutboxStatus = "PROCESSED"
	UserOutboxStatusDELIVERED UserOutboxStatus = "DELIVERED"
	UserOutboxStatusFAILED    UserOutboxStatus = "FAILED"
)

func (e *UserOutboxStatus) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = UserOutboxStatus(s)
	case string:
		*e = UserOutboxStatus(s)
	default:
		return fmt.Errorf("unsupported scan type for UserOutboxStatus: %T", src)
	}
	return nil
}

type NullUserOutboxStatus struct {
	UserOutboxStatus UserOutboxStatus
	Valid            bool // Valid is true if UserOutboxStatus is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullUserOutboxStatus) Scan(value interface{}) error {
	if value == nil {
		ns.UserOutboxStatus, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.UserOutboxStatus.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullUserOutboxStatus) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.UserOutboxStatus), nil
}

type User struct {
	ID        uuid.UUID
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
	CreatedBy uuid.UUID
	UpdatedBy uuid.UUID
	DeletedBy *uuid.UUID
}

type UsersOutbox struct {
	ID        uuid.UUID
	Payload   *entity.User
	Status    UserOutboxStatus
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
	CreatedBy uuid.UUID
	UpdatedBy uuid.UUID
	DeletedBy *uuid.UUID
}
