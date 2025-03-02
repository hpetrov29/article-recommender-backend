package like

import (
	"context"
	"time"

	"github.com/hpetrov29/resttemplate/internal/logger"
)

type Storer interface {
	Publish(ctx context.Context, like Like) error
}

// Core manages the set of APIs for posts api access
type Core struct {
	storer Storer
	log    *logger.Logger
}

// NewCore constructs and returns a new Core instance for post API access.
//
// Parameters:
//   - st: struct that implements the Storer interface for repository operations.
//   - log: pointer to the logger used for logging within the core.
func NewCore(s Storer, log *logger.Logger) *Core {
	return &Core{
		storer: s,
		log:    log,
	}
}

func (c *Core) Publish(ctx context.Context, newLike NewLike) (error) {
	like := Like{
		Value: newLike.Value,
		UserId: newLike.UserId,
		PostId: newLike.PostId,
		CreatedAt: time.Now(),
	}

	return c.storer.Publish(ctx, like)
}