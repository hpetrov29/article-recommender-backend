package posts

import (
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/hpetrov29/resttemplate/business/core/post"
	"github.com/hpetrov29/resttemplate/internal/validate"
)

func parseFilter(r *http.Request) (post.QueryFilter, error) {
	const (
		filterByUserId   	= "user_id"
		filterByDateCreated = "date_created"
		filterByDateUpdated = "date_updated"
	)

	values := r.URL.Query()

	var filter post.QueryFilter

	if userId := values.Get(filterByUserId); userId != "" {
		uid, err := uuid.Parse(userId)
		fmt.Println(uid)
		if err != nil {
			return post.QueryFilter{}, validate.NewFieldsError(filterByUserId, err)
		}
		filter.WithUserId(uid)
	}

	if dateCreated := values.Get(filterByDateCreated); dateCreated != "" {
		dc, err := time.Parse("2006-01-02T15:04:05Z", dateCreated)
		if err != nil {
			return post.QueryFilter{}, validate.NewFieldsError(filterByDateCreated, err)
		}
		filter.WithDateCreated(dc)
	}

	if dateUpdated := values.Get(filterByDateUpdated); dateUpdated != "" {
		du, err := time.Parse("2006-01-02T15:04:05Z", dateUpdated)
		if err != nil {
			return post.QueryFilter{}, validate.NewFieldsError(filterByDateUpdated, err)
		}
		filter.WithDateUpdated(du)
	}
	if err := filter.Validate(); err != nil {
		return post.QueryFilter{}, err
	}

	return filter, nil
}