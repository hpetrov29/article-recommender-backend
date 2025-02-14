package user

import (
	"net/mail"
	"time"

	"github.com/google/uuid"
)

// User struct contains information about an individual user.
// Meant to be used at the service/core layer
type User struct {
	ID uuid.UUID
	Name string
	Email mail.Address
	Roles []Role
	PasswordHash []byte
	Department string
	Enabled bool
	DateCreated time.Time
	DateUpdated time.Time
}

// NewUser contains information required to create a new user.
// Meant to be used at the service/core layer
type NewUser struct {
	Name string
	Email mail.Address
	Roles []Role
	Department string
	Password string
	PasswordConfirm string
}

// UpdateUser contains information required to update a user.
// Meant to be used at the service/core layer
type UpdateUser struct {
	Name *string
	Email *mail.Address
	Roles []Role
	Department *string
	Password *string
	PasswordConfirm *string
	Enabled *bool
}