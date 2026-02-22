package scraper

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

const (
	userAgent      = "Mozilla/5.0 (compatible; notifier-hub/1.0)"
	minTitleLength = 10 // ナビリンク等の短すぎるテキストを除外
)

type ScrapeScaper struct {
	client *http.Client
}

func NewScrape() *ScrapeScaper {
	return &ScrapeScaper{client: &http.Client{}}
}

func (s *ScrapeScaper) Fetch(ctx context.Context, rawURL string, keywords []string) ([]Item, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", rawURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", userAgent)

	resp, err := s.client.Do(req) //nolint:gosec // URL は DB の管理値であり任意入力ではない
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("http error: %s", resp.Status)
	}

	base, err := url.Parse(rawURL)
	if err != nil {
		return nil, err
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	seen := map[string]bool{}
	var items []Item

	doc.Find("a[href]").Each(func(_ int, sel *goquery.Selection) {
		href, _ := sel.Attr("href")

		// ページ内リンク・JS・メールは除外
		if strings.HasPrefix(href, "#") ||
			strings.HasPrefix(href, "javascript:") ||
			strings.HasPrefix(href, "mailto:") {
			return
		}

		// リンクテキストを正規化（改行・連続スペースを単一スペースに）
		title := strings.Join(strings.Fields(sel.Text()), " ")
		if len([]rune(title)) < minTitleLength {
			return
		}

		// 相対URLを絶対URLに解決
		ref, err := url.Parse(href)
		if err != nil {
			return
		}
		absURL := base.ResolveReference(ref).String()

		// 同一URLの重複除外
		if seen[absURL] {
			return
		}
		seen[absURL] = true

		if !matchesKeywords(title, keywords) {
			return
		}

		items = append(items, Item{
			Title: title,
			URL:   absURL,
		})
	})

	return items, nil
}
