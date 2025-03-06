package comments

import (
	"time"

	"github.com/hpetrov29/resttemplate/business/core/comment"
)

type Comment struct {
	Id     	   	int64 		`json:"id"`
	UserId 		int64 		`json:"userId"`
	PostId		int64 		`json:"postId"`
	ParentId 	int64 		`json:"parentId"`
	Content   	string 		`json:"content"`
	CreatedAt 	time.Time 	`json:"createdAt"`
}

type NewComment struct {
	ParentId  	int64 `json:"parentId" validate:"omitempty"`
	Content   	string `json:"content" validate:"required"`
}

func toCoreNewComment(c NewComment, userId int64, postId int64) comment.NewComment {
	return comment.NewComment{
		UserId: userId,
		PostId: postId,
		ParentId: c.ParentId,
		Content: c.Content,
	}
} 

func toAppComment(c comment.Comment) Comment {
	return Comment{
		Id: c.Id,
		UserId: c.UserId,
		PostId: c.PostId,
		ParentId: c.ParentId,
		Content: c.Content,
		CreatedAt: c.CreatedAt,
	}
}

type DeleteComment struct {
	Id     	   	uint64 		`json:"id"`
	UserId 		uint64 		`json:"userId"`
}