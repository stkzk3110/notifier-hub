package scraper

import "context"

// 取得したアイテムの共通構造
type Item struct {
	Title       string
	URL         string
	Description string
	PublishedAt string // 発売日の文字列（抽出できれば）
}

// Scraperインターフェース（RSSもスクレイピングも同じ口に統一）
type Scraper interface {
	Fetch(ctx context.Context, url string, keywords []string) ([]Item, error)
}