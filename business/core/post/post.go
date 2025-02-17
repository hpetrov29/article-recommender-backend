package post

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/hpetrov29/resttemplate/business/data/order"
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
	Query(ctx context.Context, filter QueryFilter, orderBy order.OrderBy, pageNumber int, rowsPerPage int) ([]Post, error)
}

type IdGenerator interface {
	GenerateId() (uint64, error)
}

// Core manages the set of APIs for posts api access
type Core struct {
	storer Storer
	log *logger.Logger
	idGenerator IdGenerator
}

// NewCore constructs and returns a new Core instance for post API access.
//
// Parameters:
//   - st: struct that implements the Storer interface for repository operations.
//   - log: pointer to the logger used for logging within the core.
func NewCore(s Storer, log *logger.Logger, idGen IdGenerator) *Core {
	return &Core{
		storer: s, 
		log: log,
		idGenerator: idGen,
	}
}


// Create adds a new post in the repository.
//
// Parameters:
//   - ctx: the context for the request, used for managing timeouts and cancellations.
//   - newPost: the contents of the new post to be created.
func (c *Core) Create(ctx context.Context, newPost NewPost) (Post, error) {
	now := time.Now()
	
	id, err := c.idGenerator.GenerateId()
	if err != nil {
		return Post{}, err
	}

	contentId, err := c.idGenerator.GenerateId()
	if err != nil {
		return Post{}, err
	}

	post := Post{
		Id: id,
		UserId: newPost.UserId,
		Title: newPost.Title,
		Description: newPost.Description,
		ContentId: contentId,
		CreatedAt: now,
		UpdatedAt: now,
	}
	
	if _, err := c.storer.Create(ctx, post); err != nil {
		return Post{}, fmt.Errorf("error creating a post: %w", err)
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

func (c *Core) Query(ctx context.Context, filter QueryFilter, orderBy order.OrderBy, pageNumber int, rowsPerPage int) ([]Post, error) {
	posts, err := c.storer.Query(ctx, filter, orderBy, pageNumber, rowsPerPage)
	if err != nil {
		return nil, err
	}

	return posts, nil
}