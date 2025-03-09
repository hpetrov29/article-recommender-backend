package postcache

import (
	"time"

	"github.com/hpetrov29/resttemplate/business/core/post"
)

type dbPost struct {
	Id           	int64   		`json:"id"`
	UserId       	int64	 		`json:"userId"`
	Title        	string   		`json:"title"`
	Description 	string 			`json:"description"`
	FrontImage  	string 			`json:"frontImage"`
	ContentId   	int64   		`json:"contentId"`
	Content   		dbContent   	`json:"content"`
	CreatedAt   	time.Time   	`json:"createdAt"`
	UpdatedAt  		time.Time   	`json:"updatedAt"`
}

// Content contains the entire content of a post.
type dbContent struct {
	Blocks []dbBlock `json:"blocks"`
}

// Block contains information about each content block such as type, styling, url, etc.
type dbBlock struct {
	Type    string    	`json:"type"`
	Content string    	`json:"content"`
	Styles  []dbStyle  	`json:"styles"`
	URL     string    	`json:"url"`
	Caption string    	`json:"caption"`
}

// Style contains text styling information (e.g., bold, italic), offset and length.
type dbStyle struct {
	Offset int    `json:"offset"`
	Length int    `json:"length"`
	Style  string `json:"style"`
}

func toDBPost(post post.Post) dbPost {
	return dbPost{
		Id: post.Id,
		UserId: post.UserId,
		Title: post.Title,
		Description: post.Description,
		FrontImage: post.FrontImage,
		ContentId: post.ContentId,
		Content: toDbContent(post.Content),
		CreatedAt: post.CreatedAt,
		UpdatedAt: post.CreatedAt,
	}
}

// Converts AppContent (app layer) to post.Content (core layer)
func toDbContent(c post.Content) dbContent {
	return dbContent{
		Blocks: toDbBlocks(c.Blocks),
	}
}

// Converts a slice of AppBlock (app layer) to post.Block (core layer)
func toDbBlocks(blocks []post.Block) []dbBlock {
	converted := make([]dbBlock, len(blocks))
	for i, b := range blocks {
		converted[i] = dbBlock{
			Type:    b.Type,
			Content: b.Content,
			Styles:  toDbStyles(b.Styles),
			URL:     b.URL,
			Caption: b.Caption,
		}
	}
	return converted
}

// Converts a slice of AppStyle (app layer) to post.Style (core layer)
func toDbStyles(styles []post.Style) []dbStyle {
	converted := make([]dbStyle, len(styles))
	for i, s := range styles {
		converted[i] = dbStyle{
			Offset: s.Offset,
			Length: s.Length,
			Style:  s.Style,
		}
	}
	return converted
}

func toCorePost(p dbPost) post.Post {
	return post.Post{
		Id: p.Id,
		UserId: p.UserId,
		Title: p.Title,
		Description: p.Description,
		FrontImage: p.FrontImage,
		ContentId: p.ContentId,
		Content: toCoreContent(p.Content),
		CreatedAt: p.CreatedAt,
		UpdatedAt: p.CreatedAt,
	}
}

// Converts AppContent (app layer) to post.Content (core layer)
func toCoreContent(c dbContent) post.Content {
	return post.Content{
		Blocks: toCoreBlocks(c.Blocks),
	}
}

// Converts a slice of AppBlock (app layer) to post.Block (core layer)
func toCoreBlocks(blocks []dbBlock) []post.Block {
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
func toCoreStyles(styles  []dbStyle) []post.Style {
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