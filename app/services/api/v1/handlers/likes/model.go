package likes

import (
	"github.com/hpetrov29/resttemplate/business/core/like"
)

type AppLike struct {
	Value     int8 		`json:"value"`
	UserId    uint64 	`json:"userId"`
	PostId    uint64 	`json:"postId"`
}

func toAppLike(l like.NewLike) AppLike {
	return AppLike{
		Value: l.Value,
		UserId: l.UserId,
		PostId: l.PostId,
	}
}