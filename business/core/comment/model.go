package comment

import "time"

// Comment struct contains all information about a comment.
type Comment struct {
	Id     	   	int64
	UserId 		int64
	PostId		int64
	ParentId 	int64
	Content   	string
	CreatedAt 	time.Time
}

// NewComment struct contains all information required by the core layer
// from the app layer to insert a comment into the repository
type NewComment struct {
	UserId 		int64		
	PostId		int64
	ParentId  	int64 
	Content   	string
}