package comment

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/hpetrov29/resttemplate/internal/logger"
)

var (
	ErrNotFound  = errors.New("comment not found")
)

type Storer interface {
	Create(ctx context.Context, comment Comment) (sql.Result, error)
	Delete(ctx context.Context, id uint64) error
}

type IdGenerator interface {
	GenerateId() (uint64, error)
}

// Core manages the set of APIs for comment api access
type Core struct {
	storer      Storer
	log         *logger.Logger
	idGenerator IdGenerator
}

// NewCore constructs and returns a new Core instance for comment API access.
//
// Parameters:
//   - st: struct that implements the Storer interface for repository operations.
//   - log: pointer to the logger used for logging within the core.
func NewCore(s Storer, log *logger.Logger, idGen IdGenerator) *Core {
	return &Core{
		storer:      s,
		log:         log,
		idGenerator: idGen,
	}
}

func (c *Core) Create(ctx context.Context, newComment NewComment) (Comment, error) {
	commentId, err := c.idGenerator.GenerateId()
	if err != nil {
		return Comment{}, err
	}

	comment := Comment{
		Id: int64(commentId),
		UserId: newComment.UserId,
		PostId: newComment.PostId,
		ParentId: newComment.ParentId,
		Content: newComment.Content,
		CreatedAt: time.Now(),
	}

	_, err = c.storer.Create(ctx, comment)
	if err != nil {
		return Comment{}, fmt.Errorf("error creating a comment: %w", err)
	}

	return comment, nil
}

func (c *Core) Delete(ctx context.Context, id uint64) error {
	err := c.storer.Delete(ctx, id); if err != nil {
		return fmt.Errorf("error deleting comment with id: %v: %w", id, err)
	}

	return nil
}

