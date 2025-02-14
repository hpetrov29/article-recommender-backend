package postsqldb

import (
	"time"

	"github.com/google/uuid"
	"github.com/hpetrov29/resttemplate/business/core/post"
)

// dbPost represents the structure used to transfer post data
// between the application and the database.
type dbPost struct {
	Id           uuid.UUID      `db:"id"`
	Title        string         `db:"title"`
	Content      string         `db:"content"`
	UserId		 uuid.UUID		`db:"user_id"`	
	DateCreated  time.Time      `db:"date_created"`
	DateUpdated  time.Time      `db:"date_updated"`
}

// toDBPost converts a post.Post instance (found in the service layer) to a dbPost struct suited for database operations.
//
// Parameters:
//   - post: the Post instance to be converted.
func toDBPost(post post.Post) dbPost {
	return dbPost{
		Id:           	post.Id,
		Title:          post.Title,
		Content:        post.Content,
		UserId: 		post.UserId,	
		DateCreated: 	post.DateCreated.UTC(),
		DateUpdated: 	post.DateUpdated.UTC(),
	}
}

// toCorePost converts a dbPost instance (found in the repository layer) to a post.Post struct.
//
// Parameters:
//   - dbPost: the dbPost instance to be converted.
func toCorePost(dbPost dbPost) post.Post {
	post := post.Post{
		Id:           dbPost.Id,
		Title:        dbPost.Title,
		Content: 	  dbPost.Content,
		UserId: 	  dbPost.UserId,
		DateCreated:  dbPost.DateCreated.In(time.Local),
		DateUpdated:  dbPost.DateUpdated.In(time.Local),
	}

	return post
}