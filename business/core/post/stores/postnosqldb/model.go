package postnosqldb

import (
	"github.com/hpetrov29/resttemplate/business/core/post"
)

type Content struct {
	Id 		int64 	`bson:"_id"`
	Blocks 	[]Block `bson:"blocks"`
}

type Block struct {
	Type    string  `bson:"type"`
	Content string  `bson:"content,omitempty"`
	Styles  []Style `bson:"styles,omitempty"`
	URL     string  `bson:"url,omitempty"`
	Caption string  `bson:"caption,omitempty"`
}

type Style struct {
	Offset int    `bson:"offset"`
	Length int    `bson:"length"`
	Style  string `bson:"style"`
}

// toDbContent returns the content stored inside a post.Post (core layer) instance
// to a Content struct suited for nosql db operations.
//
// Parameters:
//   - post: the Post instance from which the content is retrieved.
func toDbContent(content post.Content, contentId int64) Content {
	return Content{
		Id: 	contentId,
		Blocks: toDbBlocks(content.Blocks),
	}
}

// Converts a slice of post.Block (core layer) to Block (repository layer)
func toDbBlocks(blocks []post.Block) []Block {
	converted := make([]Block, len(blocks))
	for i, b := range blocks {
		converted[i] = Block{
			Type:    b.Type,
			Content: b.Content,
			Styles:  toDbStyles(b.Styles),
			URL:     b.URL,
			Caption: b.Caption,
		}
	}
	return converted
}

// Converts a slice of post.Style (core layer) to Style (repository layer)
func toDbStyles(styles []post.Style) []Style {
	converted := make([]Style, len(styles))
	for i, s := range styles {
		converted[i] = Style{
			Offset: s.Offset,
			Length: s.Length,
			Style:  s.Style,
		}
	}
	return converted
}

// Converts Content (repository layer) to post.Content (core layer)
func toCoreContent(c Content) post.Content {
	return post.Content{
		Blocks: toCoreBlocks(c.Blocks),
	}
}

// Converts a slice of Block (repository layer) to post.Block (core layer)
func toCoreBlocks(blocks []Block) []post.Block {
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

// Converts a slice of Style (repository layer) to post.Style (core layer)
func toCoreStyles(styles []Style) []post.Style {
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