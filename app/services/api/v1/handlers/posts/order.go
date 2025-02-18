package posts

import (
	"errors"
	"net/http"

	"github.com/hpetrov29/resttemplate/business/core/post"
	"github.com/hpetrov29/resttemplate/business/data/order"
	"github.com/hpetrov29/resttemplate/internal/validate"
)

const (
	orderByCreatedAt = "created_at"
	orderByUpdatedAt = "updated_at"
)

var orderByFields = map[string]string{
	orderByCreatedAt:   post.OrderByCreatedAt,
	orderByUpdatedAt:   post.OrderByUpdatedAt,
}

func parseOrder(r *http.Request) (order.OrderBy, error) {
	orderBy, err := order.Parse(r, order.NewBy(orderByCreatedAt, order.ASC))
	if err != nil {
		return order.OrderBy{}, err
	}

	if _, exists := orderByFields[orderBy.Field]; !exists {
		return order.OrderBy{}, validate.NewFieldsError(orderBy.Field, errors.New("order field does not exist"))
	}

	orderBy.Field = orderByFields[orderBy.Field]

	return orderBy, nil
}