package commentsqldb

import (
	"database/sql"
	"time"

	"github.com/hpetrov29/resttemplate/business/core/comment"
)

type Comment struct {
	Id     	   	int64 				`db:"id"`
	UserId 		int64 				`db:"user_id"`
	PostId 		int64				`db:"post_id"`
	ParentId 	sql.NullInt64 		`db:"parent_id"`
	Content   	string 				`db:"content"`
	CreatedAt 	time.Time 			`db:"created_at"`
}

func toDBComment(c comment.Comment) Comment {
	return Comment{
		Id: c.Id,
		UserId: c.UserId,
		PostId: c.PostId,
		ParentId: sql.NullInt64{Int64: c.ParentId, Valid: c.ParentId>0},
		Content: c.Content,
		CreatedAt: c.CreatedAt,
	}
}

func ToCoreComment(c Comment) comment.Comment {
	return comment.Comment{
		Id: c.Id,
		UserId: c.UserId,
		PostId: c.PostId,
		ParentId: c.ParentId.Int64,
		Content: c.Content,
		CreatedAt: c.CreatedAt,
	}
}

func ToCoreComments(comments []Comment) []comment.Comment {
	var slice []comment.Comment

	for _, c := range comments {
		cc := ToCoreComment(c)
		slice = append(slice, cc)
	}

	return slice
}