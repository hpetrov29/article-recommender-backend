package commentsqldb

import (
	"database/sql"
	"time"

	"github.com/hpetrov29/resttemplate/business/core/comment"
)

type Comment struct {
	Id     	   	uint64 				`db:"id"`
	UserId 		uint64 				`db:"user_id"`
	PostId 		uint64				`db:"post_id"`
	ParentId 	sql.Null[uint64] 	`db:"parent_id"`
	Content   	string 				`db:"content"`
	CreatedAt 	time.Time 			`db:"created_at"`
}

func toDBComment(c comment.Comment) Comment {
	return Comment{
		Id: c.Id,
		UserId: c.UserId,
		PostId: c.PostId,
		ParentId: sql.Null[uint64]{V: c.ParentId, Valid: c.ParentId>0},
		Content: c.Content,
		CreatedAt: c.CreatedAt,
	}
}