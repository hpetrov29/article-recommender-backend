package comment

import "time"

type Comment struct {
	Id     	   	uint64
	UserId 		uint64
	PostId		uint64
	ParentId 	uint64
	Content   	string
	CreatedAt 	time.Time
}

type NewComment struct {
	UserId 		uint64		
	PostId		uint64
	ParentId  	uint64 
	Content   	string
}