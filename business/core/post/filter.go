package post

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/hpetrov29/resttemplate/internal/validate"
)

// QueryFilter holds the available fields a query can be filtered on.
type QueryFilter struct {
	Id 			*uuid.UUID
	UserId      *uuid.UUID
	DateCreated *time.Time
	DateUpdated *time.Time
}

// Validate checks the data in the model is considered clean.
func (qf *QueryFilter) Validate() error {
	if err := validate.Check(qf); err != nil {
		return fmt.Errorf("validate: %w", err)
	}
	return nil
}

// WithUserId is used to filter posts of a specific userId
func (qf *QueryFilter) WithUserId(userId uuid.UUID) {
	qf.UserId = &userId
}

// WithDateCreated is used to filter posts on a specific creation time
func (qf *QueryFilter) WithDateCreated(dateCreated time.Time) {
	qf.DateCreated = &dateCreated
}

// WithDateUpdated is used to filter posts on a specific update time
func (qf *QueryFilter) WithDateUpdated(dateUpdated time.Time) {
	qf.DateUpdated = &dateUpdated
}
