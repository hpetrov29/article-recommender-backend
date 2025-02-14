package post

import (
	"time"

	"github.com/google/uuid"
)

// Post struct contains information about an individual user.
// Meant to be used at the service/core layer
type Post struct {
	Id 			uuid.UUID
	Title       string
	Content 	string
	UserId      uuid.UUID
	DateCreated time.Time
	DateUpdated time.Time
}

// NewPost contains information required to create a new user.
// Meant to be used at the service/core layer
type NewPost struct {
	Title       string
	Content 	string
	UserId      uuid.UUID
}

// UpdatePost contains information required to update a user.
// Meant to be used at the service/core layer
type UpdatePost struct {
	Title       string
	Content 	string
}