package post

import (
	"time"
)

// Post struct contains information about an individual user.
// Meant to be used at the service/core layer
type Post struct {
	Id          uint64
	UserId      uint64
	Title       string
	Description string
	FrontImage  string
	ContentId   uint64
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// NewPost contains information required to create a new user.
// Meant to be used at the service/core layer
type NewPost struct {
	UserId      uint64
	Title       string
	Description string
	ContentId 	uint64
}

// UpdatePost contains information required to update a user.
// Meant to be used at the service/core layer
type UpdatePost struct {
	Title       string
	Content 	string
}