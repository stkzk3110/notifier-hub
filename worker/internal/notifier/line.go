package notifier

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type LineNotifier struct {
	token  string
	userID string
}

func NewLine(token, userID string) *LineNotifier {
	return &LineNotifier{token: token, userID: userID}
}

func (l *LineNotifier) Send(ctx context.Context, p Payload) error {
	message := fmt.Sprintf("【%s】%s\n%s", p.SourceName, p.Title, p.URL)

	reqBody, err := json.Marshal(map[string]any{
		"to": l.userID,
		"messages": []map[string]string{
			{"type": "text", "text": message},
		},
	})
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, "POST",
		"https://api.line.me/v2/bot/message/push",
		bytes.NewReader(reqBody),
	)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+l.token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req) //nolint:gosec // URL is hardcoded constant, not user input
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("LINE push failed: status %d", resp.StatusCode)
	}
	return nil
}
