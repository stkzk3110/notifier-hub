package store

import (
	"context"
	"fmt"
	supa "github.com/supabase-community/supabase-go"
)

type Source struct {
	ID       string   `json:"id"`
	Name     string   `json:"name"`
	URL      string   `json:"url"`
	Type     string   `json:"type"` // "rss" | "scrape"
	Enabled  bool     `json:"enabled"`
	Keywords []string `json:"keywords"`
}

type SeenItem struct {
	SourceID string `json:"source_id"`
	URL      string `json:"url"`
	Title    string `json:"title"`
}

type NotificationChannel struct {
	ID      string                 `json:"id"`
	Type    string                 `json:"type"`
	Config  map[string]interface{} `json:"config"`
	Enabled bool                   `json:"enabled"`
}

type Store struct {
	client *supa.Client
}

func New(supabaseURL, supabaseKey string) (*Store, error) {
	client, err := supa.NewClient(supabaseURL, supabaseKey, &supa.ClientOptions{})
	if err != nil {
		return nil, fmt.Errorf("supabase client: %w", err)
	}
	return &Store{client: client}, nil
}

// 有効な監視ソース一覧を取得
func (s *Store) GetEnabledSources(ctx context.Context) ([]Source, error) {
	var sources []Source
	_, err := s.client.From("sources").
		Select("*", "exact", false).
		Eq("enabled", "true").
		ExecuteTo(&sources)
	if err != nil {
		return nil, fmt.Errorf("get sources: %w", err)
	}
	return sources, nil
}

// 既読チェック
func (s *Store) IsAlreadySeen(ctx context.Context, url string) (bool, error) {
	var items []SeenItem
	_, err := s.client.From("seen_items").
		Select("url", "exact", false).
		Eq("url", url).
		ExecuteTo(&items)
	if err != nil {
		return false, fmt.Errorf("check seen: %w", err)
	}
	return len(items) > 0, nil
}

// 既読登録
func (s *Store) MarkAsSeen(ctx context.Context, item SeenItem) error {
	_, _, err := s.client.From("seen_items").Insert(item, false, "", "", "").Execute()
	if err != nil {
		return fmt.Errorf("mark seen: %w", err)
	}
	return nil
}

// 有効な通知チャネル一覧を取得
func (s *Store) GetEnabledChannels(ctx context.Context) ([]NotificationChannel, error) {
	var channels []NotificationChannel
	_, err := s.client.From("notification_channels").
		Select("*", "exact", false).
		Eq("enabled", "true").
		ExecuteTo(&channels)
	if err != nil {
		return nil, fmt.Errorf("get channels: %w", err)
	}
	return channels, nil
}