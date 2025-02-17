package posts

import (
	"time"

	"github.com/hpetrov29/resttemplate/business/core/post"
	"github.com/hpetrov29/resttemplate/internal/validate"
)

// AppPost represents the contents of a post in the app layer.
type AppPost struct {
	Id           	uint64   	`json:"id"`
	UserId       	uint64	 	`json:"userId"`
	Title        	string   	`json:"title"`
	Description 	string 		`json:"description"`
	FrontImage  	string 		`json:"frontImage"`
	ContentId   	uint64   	`json:"contentId"`
	CreatedAt   	string   	`json:"createdAt"`
	UpdatedAt  		string   	`json:"updatedAt"`
}

func toAppPost(post post.Post) AppPost {
	return AppPost{
	Id:  			post.Id,
	UserId: 		post.UserId,
	Title: 			post.Title,
	Description:	post.Description,
	FrontImage:		post.FrontImage,
	ContentId: 		post.ContentId,
	CreatedAt: 		post.CreatedAt.Format(time.RFC3339),
	UpdatedAt:  	post.UpdatedAt.Format(time.RFC3339),
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
	Description      string   	  `json:"description" validate:"required"`
}

func toCoreNewPost(app AppNewPost) (post.NewPost, error) {
	post := post.NewPost{
		Title: app.Title,
		Description: app.Description,
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
