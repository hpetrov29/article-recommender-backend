package postsqldb

import (
	"time"

	"github.com/hpetrov29/resttemplate/business/core/post"
)

// dbPost represents the structure used to transfer post data
// between the application and the database.
type dbPost struct {
	Id          uint64    `db:"id"`
	UserId      uint64    `db:"user_id"`
	Title       string    `db:"title"`
	Description string    `db:"description"`
	FrontImage  string    `db:"front_image"`
	ContentId   uint64    `db:"content_id"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}

// toDBPost converts a post.Post instance (found in the service layer) to a dbPost struct suited for database operations.
//
// Parameters:
//   - post: the Post instance to be converted.
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

// toCorePost converts a dbPost instance (found in the repository layer) to a post.Post struct.
//
// Parameters:
//   - dbPost: the dbPost instance to be converted.
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


func toCorePostSlice(dbProducts []dbPost) []post.Post {
	prds := make([]post.Post, len(dbProducts))
	for i, dbPost := range dbProducts {
		prds[i] = toCorePost(dbPost)
	}
	return prds
}