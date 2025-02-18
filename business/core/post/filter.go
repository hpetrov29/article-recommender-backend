package post

import (
	"fmt"
	"time"

	"github.com/hpetrov29/resttemplate/internal/validate"
)

// QueryFilter holds the available fields a query can be filtered on.
type QueryFilter struct {
	UserId      *uint64
	CreatedAt *time.Time
	UpdatedAt *time.Time
}

// Validate checks the data in the model is considered clean.
func (qf *QueryFilter) Validate() error {
	if err := validate.Check(qf); err != nil {
		return fmt.Errorf("validate: %w", err)
	}
	return nil
}

// WithUserId is used to filter posts of a specific userId
func (qf *QueryFilter) WithUserId(userId uint64) {
	qf.UserId = &userId
}

// WithCreatedAt is used to filter posts on a specific creation time
func (qf *QueryFilter) WithCreatedAt(createdAt time.Time) {
	qf.CreatedAt = &createdAt
}

// WithUpdatedAt is used to filter posts on a specific update time
func (qf *QueryFilter) WithUpdatedAt(updatedAt time.Time) {
	qf.UpdatedAt = &updatedAt
}
