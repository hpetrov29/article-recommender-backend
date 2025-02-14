package post

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/hpetrov29/resttemplate/internal/logger"
)

// Set of error variables for CRUD operations.
var (
	ErrNotFound  = errors.New("post not found")
)

type Storer interface {
	Create(ctx context.Context, post Post) (sql.Result, error)
	Delete(ctx context.Context, post Post) error
	QueryById(ctx context.Context, id string) (Post, error)
	GetPosts(ctx context.Context) ([]Post, error)
}

// Core manages the set of APIs for posts api access
type Core struct {
	storer Storer
	log *logger.Logger
}

// NewCore constructs and returns a new Core instance for post API access.
//
// Parameters:
//   - st: struct that implements the Storer interface for repository operations.
//   - log: pointer to the logger used for logging within the core.
func NewCore(s Storer, log *logger.Logger) *Core {
	return &Core{
		storer: s, 
		log: log,
	}
}


// Create adds a new post in the repository.
//
// Parameters:
//   - ctx: the context for the request, used for managing timeouts and cancellations.
//   - newPost: the contents of the new post to be created.
func (c *Core) Create(ctx context.Context, newPost NewPost) (Post, error) {
	now := time.Now()
	post := Post{
		Title: newPost.Title,
		Content: newPost.Content,
		UserId: newPost.UserId,
		DateCreated: now,
		DateUpdated: now,
	}
	
	if _, err := c.storer.Create(ctx, post); err != nil {
		return Post{}, fmt.Errorf("create: %w", err)
	}
	return post, nil
}

// Delete removes a specified post from the repository.
//
// Parameters:
//   - ctx: the context for the request, used for managing timeouts and cancellations.
//   - post: the post to be deleted.
func (c *Core) Delete(ctx context.Context, post Post) error {
	if err := c.storer.Delete(ctx, post); err != nil {
		return fmt.Errorf("delete: %w", err)
	}

	return nil
}

func (c *Core) GetPostById(ctx context.Context, id string) (Post, error) {
	post, err := c.storer.QueryById(ctx, id)
	if err != nil {
		return Post{}, err
	}
	return post, nil
}

func (c *Core) GetPosts(ctx context.Context) ([]Post, error) {
	c.storer.GetPosts(ctx)
	return []Post{}, nil
}