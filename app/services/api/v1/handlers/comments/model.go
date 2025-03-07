package comments

import (
	"time"

	"github.com/hpetrov29/resttemplate/business/core/comment"
)

// AppComment represents the payload of a newly created comment in the app layer
type AppComment struct {
	Id     	   	int64 		`json:"id"`
	UserId 		int64 		`json:"userId"`
	PostId		int64 		`json:"postId"`
	ParentId 	int64 		`json:"parentId"`
	Content   	string 		`json:"content"`
	CreatedAt 	time.Time 	`json:"createdAt"`
}

// NewAppComment represents the contents that need to be decoded from the request 
// body in the app layer for comment creation
type NewAppComment struct {
	ParentId  	int64 `json:"parentId" validate:"omitempty"`
	Content   	string `json:"content" validate:"required"`
}

// Converts NewAppComment (app layer) to comment.NewComment (core layer)
func toCoreNewComment(c NewAppComment, userId int64, postId int64) comment.NewComment {
	return comment.NewComment{
		UserId: userId,
		PostId: postId,
		ParentId: c.ParentId,
		Content: c.Content,
	}
} 

// Converts comment.Comment (core layer) to comment.AppComment (app layer)
func toAppComment(c comment.Comment) AppComment {
	return AppComment{
		Id: c.Id,
		UserId: c.UserId,
		PostId: c.PostId,
		ParentId: c.ParentId,
		Content: c.Content,
		CreatedAt: c.CreatedAt,
	}
}

// DeleteComment represents the contents that need to be decoded from the request
// body in the app layer for comment deletion.
type DeleteComment struct {
	Id     	   	uint64 		`json:"id"`
	UserId 		uint64 		`json:"userId"`
}