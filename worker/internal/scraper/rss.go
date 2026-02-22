package scraper

import (
	"context"
	"strings"

	"github.com/mmcdole/gofeed"
)

type RSSScaper struct {
	parser *gofeed.Parser
}

func NewRSS() *RSSScaper {
	return &RSSScaper{parser: gofeed.NewParser()}
}

func (r *RSSScaper) Fetch(ctx context.Context, url string, keywords []string) ([]Item, error) {
	feed, err := r.parser.ParseURLWithContext(url, ctx)
	if err != nil {
		return nil, err
	}

	var items []Item
	for _, entry := range feed.Items {
		if !matchesKeywords(entry.Title+" "+entry.Description, keywords) {
			continue
		}
		items = append(items, Item{
			Title:       entry.Title,
			URL:         entry.Link,
			Description: entry.Description,
			PublishedAt: entry.Published,
		})
	}
	return items, nil
}

// keywordsが空の場合は全件対象。1つでも含まれていればマッチ。
func matchesKeywords(text string, keywords []string) bool {
	if len(keywords) == 0 {
		return true
	}
	for _, kw := range keywords {
		if strings.Contains(text, kw) {
			return true
		}
	}
	return false
}
