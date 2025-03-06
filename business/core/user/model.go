package user

import (
	"net/mail"
	"time"
)

// User struct contains information about an individual user.
// Meant to be used at the service/core layer
type User struct {
	Id int64
	Username string
	Email mail.Address
	Roles []Role
	PasswordHash []byte
	CreatedAt time.Time
}

// NewUser contains information required to create a new user.
// Meant to be used at the service/core layer
type NewUser struct {
	Username string
	Email mail.Address
	Roles []Role
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