package postsqldb

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/hpetrov29/resttemplate/business/core/post"
)

func (s *Store) applyFilter(filter post.QueryFilter, data map[string]interface{}, buf *bytes.Buffer) {
	var wc []string

	if filter.Id != nil {
		data["id"] = *filter.Id
		wc = append(wc, "id = :id")
	}

	if filter.UserId != nil {
		data["user_id"] = filter.UserId.String()
		fmt.Println(data["user_id"])
		wc = append(wc, "user_id = :user_id")
	}

	if filter.DateCreated != nil {
		data["date_created"] = *filter.DateCreated
		wc = append(wc, "date_created = :date_created")
	}

	if filter.DateUpdated != nil {
		data["date_updated"] = *filter.DateUpdated
		wc = append(wc, "date_updated = :date_updated")
	}

	if len(wc) > 0 {
		buf.WriteString(" WHERE ")
		buf.WriteString(strings.Join(wc, " AND "))
	}
}