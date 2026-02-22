# CLAUDE.md

このファイルはClaude Code (claude.ai/code) がこのリポジトリで作業する際のガイドです。
**ユーザーへの返答は必ず日本語で行うこと。**

## コミットメッセージ規約

[Conventional Commits](https://www.conventionalcommits.org/ja/v1.0.0/) に準拠し、**日本語**で記述すること。

### フォーマット

```
<type>(<scope>): <日本語の説明>
```

### type 一覧

| type | 用途 |
|---|---|
| `feat` | 新機能の追加 |
| `fix` | バグ修正 |
| `docs` | ドキュメントのみの変更 |
| `chore` | ビルド・設定・ツールなどの変更 |
| `refactor` | リファクタリング（バグ修正・機能追加を除く） |
| `test` | テストの追加・修正 |
| `ci` | CI/CD の変更 |

### 例

```
feat(notifier): LINEノティファイアを実装
fix(scraper): RSSフィードのパースエラーを修正
chore(ci): golangci-lint ワークフローを追加
docs: CLAUDE.md にコミット規約を追記
```

---

## コマンド

Goのコマンドはすべて `worker/` ディレクトリ内で実行する：

```bash
# ローカル実行（環境変数が必要）
SUPABASE_URL=... SUPABASE_SERVICE_ROLE_KEY=... go run ./cmd/main.go

# ビルド
go build ./cmd/main.go

# テスト全件
go test ./...

# 特定テストのみ実行
go test ./internal/scraper/ -run TestFoo

# 依存関係の整理
go mod tidy
```

Supabaseのローカル開発（リポジトリルートで実行）：
```bash
supabase start        # ローカルSupabaseを起動
supabase db reset     # マイグレーションをゼロから適用
supabase stop
```

## アーキテクチャ

食品チェーン（マクドナルド・スターバックス・KFC等）の新商品・期間限定情報をスクレイピングし、LINE や Google Calendar へ通知するGoバッチワーカー。GitHub Actionsで毎朝8時JST（UTC 23時）に自動実行（`workflow_dispatch`で手動実行も可）。

### データフロー

```
Supabase DB (sources) → Scraper.Fetch() → IsAlreadySeen() → Notifier.Send() → MarkAsSeen()
```

`cmd/main.go` がパイプライン全体を制御する。DBから有効なソースを取得し、`type` フィールド（`rss` | `scrape`）で振り分けてスクレイパーに渡し、`seen_items` で重複除外、有効な通知チャネル全件に送信する。

### 主要インターフェース

- `scraper.Scraper` — `Fetch(ctx, url, keywords) ([]Item, error)`。実装: `rss.go`（未実装）、`scrape.go`（未実装）
- `notifier.Notifier` — `Send(ctx, Payload) error`。実装: `line.go`（実装済み）、`calendar.go`（未実装）
- `store.Store` — `supabase-go` クライアントのラッパー。メソッド: `GetEnabledSources`, `GetEnabledChannels`, `IsAlreadySeen`, `MarkAsSeen`

### スクレイパーの追加手順

1. `worker/internal/scraper/` に `scraper.Scraper` を実装する
2. `main.go` のソースタイプ `switch` に `case` を追加する

### ノティファイアの追加手順

1. `worker/internal/notifier/` に `notifier.Notifier` を実装する
2. `main.go` の `buildNotifiers()` に `case` を追加する
3. 新しいマイグレーションで `notification_channels.type` のCHECK制約にチャネル種別を追加する

### Supabase

- マイグレーションは `supabase/migrations/` に格納。既存ファイルを編集せず、必ず新規ファイルを作成する。
- `notification_channels.config` は `jsonb`。LINE ノティファイアは `config["token"]` を参照する。
- `seen_items.url` には `unique` 制約があり、DB レベルでも重複排除される。

### 環境変数

| 変数名 | 説明 |
|---|---|
| `SUPABASE_URL` | SupabaseプロジェクトのURL |
| `SUPABASE_SERVICE_ROLE_KEY` | サービスロールキー（RLSをバイパスする） |

LINEノティファイアのトークンは環境変数ではなく、`notification_channels.config` の各行に格納する。
