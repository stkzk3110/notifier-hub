//go:build integration

package scraper

import (
	"context"
	"testing"
	"time"
)

func TestScrapeScaper_McDonalds(t *testing.T) {
	s := NewScrape()
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	items, err := s.Fetch(ctx, "https://www.mcdonalds.co.jp/company/news/", []string{})
	if err != nil {
		t.Fatalf("fetch error: %v", err)
	}
	t.Logf("取得件数: %d 件", len(items))
	for i, item := range items {
		t.Logf("[%d] %s\n    %s", i+1, item.Title, item.URL)
	}
	if len(items) == 0 {
		t.Error("アイテムが0件です。セレクタまたはURLを確認してください")
	}
}

func TestScrapeScaper_Starbucks(t *testing.T) {
	s := NewScrape()
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	items, err := s.Fetch(ctx, "https://www.starbucks.co.jp/press_release/", []string{})
	if err != nil {
		t.Fatalf("fetch error: %v", err)
	}
	t.Logf("取得件数: %d 件", len(items))
	for i, item := range items {
		t.Logf("[%d] %s\n    %s", i+1, item.Title, item.URL)
	}
	if len(items) == 0 {
		t.Error("アイテムが0件です。セレクタまたはURLを確認してください")
	}
}

func TestScrapeScaper_Keywords(t *testing.T) {
	s := NewScrape()
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	keywords := []string{"期間限定", "新発売", "季節限定"}
	items, err := s.Fetch(ctx, "https://www.mcdonalds.co.jp/company/news/", keywords)
	if err != nil {
		t.Fatalf("fetch error: %v", err)
	}
	t.Logf("キーワードフィルタ後: %d 件", len(items))
	for i, item := range items {
		t.Logf("[%d] %s", i+1, item.Title)
	}
}
