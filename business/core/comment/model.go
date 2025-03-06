package comment

import "time"

type Comment struct {
	Id     	   	int64
	UserId 		int64
	PostId		int64
	ParentId 	int64
	Content   	string
	CreatedAt 	time.Time
}

type NewComment struct {
	UserId 		int64		
	PostId		int64
	ParentId  	int64 
	Content   	string
}