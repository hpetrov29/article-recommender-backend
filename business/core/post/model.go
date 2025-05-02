package post

import (
	"time"
)

// =============================================================================
// Post metadata related models

// Post struct contains all information about a single post
type Post struct {
	Id          int64
	UserId      int64
	Title       string
	Description string
	FrontImage  string
	ContentId   int64
	Content 	Content
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// NewPost contains information required to create a new post
// Meant to be used at the service/core layer
type NewPost struct {
	UserId      int64
	Title       string
	Description string
	Content 	Content
}

// UpdatePost contains information required to update a post
// Meant to be used at the service/core layer
type UpdatePost struct {
	
}

// =============================================================================
// Post content related models

// Content contains the entire content of a post
type Content struct {
	Blocks []Block
}

// Block contains information about each content block such as type, styling, url, etc.
type Block struct {
	Type    string
	Content string
	Styles  []Style
	URL     string
	Caption string
}

// Style contains text styling information (e.g., bold, italic), offset and length
type Style struct {
	Offset int
	Length int
	Style  string
}