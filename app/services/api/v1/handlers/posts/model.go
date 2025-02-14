package posts

import (
	"time"

	"github.com/google/uuid"
	"github.com/hpetrov29/resttemplate/business/core/post"
	"github.com/hpetrov29/resttemplate/internal/validate"
)

// AppPost represents the contents of a post in the app layer.
type AppPost struct {
	Id           	string   	`json:"id"`
	Title        	string   	`json:"title"`
	Content      	string   	`json:"content"`
	UserId       	string	 	`json:"userId"`
	DateCreated  	string   	`json:"dateCreated"`
	DateUpdated  	string   	`json:"dateUpdated"`
}

func toAppPost(post post.Post) AppPost {
	return AppPost{
	Id:  			post.Id.String(),
	Title: 			post.Title,
	Content: 		post.Content,
	UserId: 		post.UserId.String(),
	DateCreated: 	post.DateCreated.Format(time.RFC3339),
	DateUpdated:  	post.DateUpdated.Format(time.RFC3339),
	}
}

func toAppPosts(posts []post.Post) []AppPost {
	items := make([]AppPost, len(posts))
	for i, post := range posts {
		items[i] = toAppPost(post)
	}

	return items
}

// =============================================================================

// AppNewUser contains information needed to create a new user.
type AppNewPost struct {
	Title            string   	  `json:"title" validate:"required"`
	Content          string   	  `json:"content" validate:"required"`
	UserId        	 string 	  
}

func toCoreNewPost(app AppNewPost) (post.NewPost, error) {
	userId, err := uuid.Parse(app.UserId)
	if err != nil {
		return post.NewPost{}, err
	}
	post := post.NewPost{
		Title: app.Title,
		Content: app.Content,
		UserId: userId,
	}

	return post, nil
}

// Validate checks the data in the model is considered clean.
func (app AppNewPost) Validate() error {
	if err := validate.Check(app); err != nil {
		return err
	}

	return nil
}

// =============================================================================

// AppUpdateUser contains information needed to update a user.
type UpdatePost struct {
	Title        	string   	`json:"title" validate:"required"`
	Content      	string   	`json:"content" validate:"required"`
}

// Validate checks the data in the model is considered clean.
func (up UpdatePost) Validate() error {
	if err := validate.Check(up); err != nil {
		return err
	}

	return nil
}
