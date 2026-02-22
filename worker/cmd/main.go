package main

import (
	"context"
	"log"
	"os"

	"github.com/yourname/notify-hub/worker/internal/notifier"
	"github.com/yourname/notify-hub/worker/internal/scraper"
	"github.com/yourname/notify-hub/worker/internal/store"
)

func main() {
	ctx := context.Background()

	// 環境変数から設定を読み込む
	supabaseURL := mustEnv("SUPABASE_URL")
	supabaseKey := mustEnv("SUPABASE_SERVICE_ROLE_KEY")

	st, err := store.New(supabaseURL, supabaseKey)
	if err != nil {
		log.Fatalf("store init: %v", err)
	}

	// 監視ソース取得
	sources, err := st.GetEnabledSources(ctx)
	if err != nil {
		log.Fatalf("get sources: %v", err)
	}

	// 通知チャネル構築
	channels, err := st.GetEnabledChannels(ctx)
	if err != nil {
		log.Fatalf("get channels: %v", err)
	}
	notifiers := buildNotifiers(channels)

	// スクレイパー
	rssScaper := scraper.NewRSS()

	for _, src := range sources {
		var items []scraper.Item
		var fetchErr error

		switch src.Type {
		case "rss":
			items, fetchErr = rssScaper.Fetch(ctx, src.URL, src.Keywords)
		default:
			log.Printf("unsupported type %s for %s, skipping", src.Type, src.Name)
			continue
		}

		if fetchErr != nil {
			log.Printf("fetch error [%s]: %v", src.Name, fetchErr)
			continue
		}

		for _, item := range items {
			seen, err := st.IsAlreadySeen(ctx, item.URL)
			if err != nil {
				log.Printf("seen check error: %v", err)
				continue
			}
			if seen {
				continue
			}

			// 各通知チャネルに送信
			payload := notifier.Payload{
				SourceName:  src.Name,
				Title:       item.Title,
				URL:         item.URL,
				Description: item.Description,
			}
			for _, n := range notifiers {
				if err := n.Send(ctx, payload); err != nil {
					log.Printf("notify error: %v", err)
				}
			}

			// 既読登録
			_ = st.MarkAsSeen(ctx, store.SeenItem{
				SourceID: src.ID,
				URL:      item.URL,
				Title:    item.Title,
			})
		}
	}

	log.Println("done")
}

func buildNotifiers(channels []store.NotificationChannel) []notifier.Notifier {
	var notifiers []notifier.Notifier
	for _, ch := range channels {
		switch ch.Type {
		case "line":
			token, ok1 := ch.Config["token"].(string)
			userID, ok2 := ch.Config["user_id"].(string)
			if !ok1 || !ok2 {
				log.Printf("line channel missing config (token/user_id), skipping")
				continue
			}
			notifiers = append(notifiers, notifier.NewLine(token, userID))
		// case "google_calendar": 後で追加
		}
	}
	return notifiers
}

func mustEnv(key string) string {
	v := os.Getenv(key)
	if v == "" {
		log.Fatalf("env %s is required", key)
	}
	return v
}