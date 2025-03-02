package like

import "time"

type NewLike struct {
	Value   int8     // 1 = like, 0 = canceled like/dislike, -1 = dislike
	UserId 	uint64
	PostId 	uint64
}

type Like struct {
	Value     int8
	UserId    uint64
	PostId    uint64
	CreatedAt time.Time
}