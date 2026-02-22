package notifier

import "context"

type Payload struct {
	Title       string
	URL         string
	Description string
	SourceName  string
}

// 通知チャネルのインターフェース（LINE・Calendarを統一）
type Notifier interface {
	Send(ctx context.Context, payload Payload) error
}