#!/bin/bash

INPUT=$(cat)
FILE_PATH=$(echo "$INPUT" | jq -r '.tool_input.file_path // empty')

# .go ファイル以外はスキップ
if [[ ! "$FILE_PATH" =~ \.go$ ]]; then
  exit 0
fi

# golangci-lint がインストールされているか確認
if ! command -v golangci-lint &>/dev/null; then
  echo "golangci-lint not found. Install: brew install golangci-lint" >&2
  exit 0
fi

# worker/ ディレクトリ内のファイルのみ対象
WORKER_DIR="$(cd "$(dirname "$0")/../.." && pwd)/worker"
if [[ ! "$FILE_PATH" =~ worker/ ]]; then
  exit 0
fi

echo "Running golangci-lint..." >&2
cd "$WORKER_DIR" || exit 0
golangci-lint run ./... 2>&1
