package postsqldb

import (
	"bytes"
	"fmt"

	"github.com/hpetrov29/resttemplate/business/core/post"
	"github.com/hpetrov29/resttemplate/business/data/order"
)

var orderByFields = map[string]string{
	post.OrderByDateCreated:  "date_created",
	post.OrderByDateUpdated:  "date_updated",
}

func (s *Store) orderByClause(orderBy order.OrderBy, buf *bytes.Buffer) (error) {
	by, exists := orderByFields[orderBy.Field]
	if !exists {
		return fmt.Errorf("field %q does not exist", orderBy.Field)
	}

	buf.WriteString(" ORDER BY " + by + " " + orderBy.Direction)

	return nil
}
