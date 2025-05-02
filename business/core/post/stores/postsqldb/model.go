package postsqldb

import (
	"time"

	"github.com/hpetrov29/resttemplate/business/core/post"
)

// dbPost represents the structure used to transfer post data
// between the core and repository layers
type dbPost struct {
	Id          int64    	`db:"id"`
	UserId      int64    	`db:"user_id"`
	Title       string    	`db:"title"`
	Description string    	`db:"description"`
	FrontImage  string    	`db:"front_image"`
	ContentId   int64    	`db:"content_id"`
	CreatedAt   time.Time 	`db:"created_at"`
	UpdatedAt   time.Time 	`db:"updated_at"`
}

// Converts post.Post (core layer) to dbPost (repository layer)
func toDBPost(post post.Post) dbPost {
	return dbPost{
		Id:           	post.Id,
		UserId: 		post.UserId,	
		Title:          post.Title,
		Description:    post.Description,
		FrontImage: 	post.FrontImage,
		ContentId: 		post.ContentId,
		CreatedAt: 		post.CreatedAt.UTC(),
		UpdatedAt: 		post.UpdatedAt.UTC(),
	}
}

// Converts dbPost (repository layer) to post.Post (core layer)
func toCorePost(dbPost dbPost) post.Post {
	post := post.Post{
		Id:           	dbPost.Id,
		UserId: 		dbPost.UserId,	
		Title:          dbPost.Title,
		Description:    dbPost.Description,
		FrontImage: 	dbPost.FrontImage,
		ContentId: 		dbPost.ContentId,
		CreatedAt: 		dbPost.CreatedAt.In(time.Local),
		UpdatedAt: 		dbPost.UpdatedAt.In(time.Local),
	}

	return post
}

// Converts a slice of dbPost (repository layer) to a slice of post.Post (core layer)
func toCorePostSlice(dbPosts []dbPost) []post.Post {
	prds := make([]post.Post, len(dbPosts))
	for i, dbPost := range dbPosts {
		prds[i] = toCorePost(dbPost)
	}
	return prds
}