package likemessaging

import (
	"encoding/json"

	"github.com/hpetrov29/resttemplate/business/core/like"
)

type Like struct {
	Value   	int8
	UserId 		uint64
	PostId 		uint64
}

func toDbLike(l like.Like) Like {
	return Like{
		Value: 		l.Value,
		UserId: 	l.UserId,
		PostId: 	l.PostId,
	}
}

func toBytes(l Like) ([]byte, error) {
	return json.Marshal(l)
}