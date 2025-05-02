package posts

import (
	"time"

	"github.com/hpetrov29/resttemplate/business/core/post"
	"github.com/hpetrov29/resttemplate/internal/validate"
)

// AppPost represents the contents of a post in the app layer.
type AppPost struct {
	Id           	int64   		`json:"id"`
	UserId       	int64	 		`json:"userId"`
	Title        	string   		`json:"title"`
	Description 	string 			`json:"description"`
	FrontImage  	string 			`json:"frontImage"`
	ContentId   	int64   		`json:"contentId"`
	Content   		*AppContent   	`json:"content,omitempty"`
	CreatedAt   	string   		`json:"createdAt"`
	UpdatedAt  		string   		`json:"updatedAt"`
}

func toAppPost(post post.Post) AppPost {
	return AppPost{
	Id:  			post.Id,
	UserId: 		post.UserId,
	Title: 			post.Title,
	Description:	post.Description,
	FrontImage:		post.FrontImage,
	ContentId: 		post.ContentId,
	Content: 		toAppContent(post.Content),
	CreatedAt: 		post.CreatedAt.Format(time.RFC3339),
	UpdatedAt:  	post.UpdatedAt.Format(time.RFC3339),
	}
}

// Converts a slice of post.Post (core layer) to a slice of AppPost (app layer)
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
	Title            string   	  	`json:"title" validate:"required"`
	Description      string   	  	`json:"description" validate:"required"`
	Content 		 AppContent 	`json:"content" validate:"required"`
}

func toCoreNewPost(app AppNewPost, userId int64) post.NewPost {
	post := post.NewPost{
		UserId: userId,
		Title: app.Title,
		Description: app.Description,
		Content:     toCoreContent(app.Content),
	}

	return post
}

// Validate checks the data in the model is considered clean.
func (app AppNewPost) Validate() error {
	if err := validate.Check(app); err != nil {
		return err
	}

	return nil
}

// =============================================================================

// AppUpdatePost contains information needed to update a post.
type AppUpdatePost struct {
	Title        	string   	`json:"title" validate:"required"`
	Content      	string   	`json:"content" validate:"required"`
}

// Validate checks the data in the model is considered clean.
func (up AppUpdatePost) Validate() error {
	if err := validate.Check(up); err != nil {
		return err
	}

	return nil
}

// =============================================================================
// Content related models and functions

// Content contains the entire content of a post.
type AppContent struct {
	Blocks []AppBlock `json:"blocks,omitempty" validate:"required,min=1,dive"`
}

// Block contains information about each content block such as type, styling, url, etc.
type AppBlock struct {
	Type    string    	`json:"type" validate:"required"`
	Content string    	`json:"content,omitempty"`
	Styles  []AppStyle  `json:"styles,omitempty"`
	URL     string    	`json:"url,omitempty"`
	Caption string    	`json:"caption,omitempty"`
}

// Style contains text styling information (e.g., bold, italic), offset and length.
type AppStyle struct {
	Offset int    `json:"offset" validate:"required"`
	Length int    `json:"length" validate:"required"`
	Style  string `json:"style" validate:"required"`
}

// Converts AppContent (app layer) to post.Content (core layer)
func toCoreContent(c AppContent) post.Content {
	return post.Content{
		Blocks: toCoreBlocks(c.Blocks),
	}
}

// Converts a slice of AppBlock (app layer) to post.Block (core layer)
func toCoreBlocks(blocks []AppBlock) []post.Block {
	converted := make([]post.Block, len(blocks))
	for i, b := range blocks {
		converted[i] = post.Block{
			Type:    b.Type,
			Content: b.Content,
			Styles:  toCoreStyles(b.Styles),
			URL:     b.URL,
			Caption: b.Caption,
		}
	}
	return converted
}

// Converts a slice of AppStyle (app layer) to post.Style (core layer)
func toCoreStyles(styles []AppStyle) []post.Style {
	converted := make([]post.Style, len(styles))
	for i, s := range styles {
		converted[i] = post.Style{
			Offset: s.Offset,
			Length: s.Length,
			Style:  s.Style,
		}
	}
	return converted
}

// Converts post.Content (core layer) to AppContent (app layer)
func toAppContent(c post.Content) *AppContent {
	if len(c.Blocks) == 0 {
		return nil
	}
	return &AppContent{
		Blocks: toAppBlocks(c.Blocks),
	}
}

// Converts a slice of post.Block (core layer) to AppBlock (app layer)
func toAppBlocks(blocks []post.Block) []AppBlock {
	converted := make([]AppBlock, len(blocks))
	for i, b := range blocks {
		converted[i] = AppBlock{
			Type:    b.Type,
			Content: b.Content,
			Styles:  toAppStyles(b.Styles),
			URL:     b.URL,
			Caption: b.Caption,
		}
	}
	return converted
}

// Converts a slice of post.Style (core layer) to AppStyle (app layer)
func toAppStyles(styles []post.Style) []AppStyle {
	converted := make([]AppStyle, len(styles))
	for i, s := range styles {
		converted[i] = AppStyle{
			Offset: s.Offset,
			Length: s.Length,
			Style:  s.Style,
		}
	}
	return converted
}