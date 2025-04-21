package postsqldb

import (
	"database/sql"
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

// =============================================================================

// dbComemnt represents the structure used to transfer comment data
// associated to a post between the core and repository layers
type dbComment struct {
	Id     	   	int64 				`db:"id"`
	UserId 		int64 				`db:"user_id"`
	ParentId 	sql.NullInt64 		`db:"parent_id"`
	Content   	string 				`db:"content"`
	CreatedAt 	time.Time 			`db:"created_at"`
	Level 		int 				`db:"lvl"`
	Root 		int 				`db:"root"`
}

// Converts dbComment (repository layer) to post.Comment (core layer)
func toCoreComment(c dbComment) post.Comment {
	return post.Comment{
		Id:        c.Id,
		UserId:    c.UserId,
		ParentId:  c.ParentId.Int64,
		Content:   c.Content,
		CreatedAt: c.CreatedAt,
	}
}

// Converts a slice of dbComment (repository layer) to a slice of post.Comment (core layer)
func toCoreComments(comments []dbComment) []post.Comment {
	slice := make([]post.Comment, len(comments))

	for i, c := range comments {
		slice[i] = toCoreComment(c)
	}

	return slice
}