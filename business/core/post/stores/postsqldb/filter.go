package postsqldb

import (
	"bytes"
	"strings"

	"github.com/hpetrov29/resttemplate/business/core/post"
)

func (s *Store) applyFilter(filter post.QueryFilter, data map[string]interface{}, buf *bytes.Buffer) {
	var wc []string
	if filter.UserId != nil {
		data["user_id"] = filter.UserId
		wc = append(wc, "user_id = :user_id")
	}

	var timeConditions []string
	if filter.CreatedAt != nil {
		data["created_at"] = *filter.CreatedAt
		timeConditions = append(timeConditions, "created_at > :created_at")
	}

	if filter.UpdatedAt != nil {
		data["updated_at"] = *filter.UpdatedAt
		timeConditions = append(timeConditions, "updated_at > :updated_at")
	}

	if len(timeConditions) > 0 {
		wc = append(wc, "("+strings.Join(timeConditions, " OR ")+")")
	}

	if len(wc) > 0 {
		buf.WriteString(" WHERE ")
		buf.WriteString(strings.Join(wc, " AND "))
	}
}