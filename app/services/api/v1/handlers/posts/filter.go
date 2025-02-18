package posts

import (
	"net/http"
	"strconv"
	"time"

	"github.com/hpetrov29/resttemplate/business/core/post"
	"github.com/hpetrov29/resttemplate/internal/validate"
)

func parseFilter(r *http.Request) (post.QueryFilter, error) {
	const (
		filterByUserId   	= "user_id"
		filterByCreatedAt = "created_at"
		filterByUpdatedAt = "updated_at"
	)

	values := r.URL.Query()

	var filter post.QueryFilter

	if userId := values.Get(filterByUserId); userId != "" {
		uid, err := strconv.ParseUint(userId, 10, 64)
		if err != nil {
			return post.QueryFilter{}, validate.NewFieldsError(filterByUserId, err)
		}
		filter.WithUserId(uid)
	}

	if createdAt := values.Get(filterByCreatedAt); createdAt != "" {
		dc, err := time.Parse("2006-01-02T15:04:05Z", createdAt)
		if err != nil {
			return post.QueryFilter{}, validate.NewFieldsError(filterByCreatedAt, err)
		}
		filter.WithCreatedAt(dc)
	}

	if updatedAt := values.Get(filterByUpdatedAt); updatedAt != "" {
		du, err := time.Parse("2006-01-02T15:04:05Z", updatedAt)
		if err != nil {
			return post.QueryFilter{}, validate.NewFieldsError(filterByUpdatedAt, err)
		}
		filter.WithUpdatedAt(du)
	}
	if err := filter.Validate(); err != nil {
		return post.QueryFilter{}, err
	}

	return filter, nil
}