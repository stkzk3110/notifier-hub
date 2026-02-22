#!/bin/bash
# Claude のレスポンス完了後に変更をまとめてコミットする

cd "$CLAUDE_PROJECT_DIR" || exit 0

# 変更・未追跡ファイルをすべてステージング
git add -A

# ステージングされた変更がなければスキップ
if git diff --cached --quiet; then
  exit 0
fi

# 変更ファイルの情報収集
CHANGED_FILES=$(git diff --cached --name-only)
ADDED_FILES=$(git diff --cached --name-only --diff-filter=A)
MODIFIED_FILES=$(git diff --cached --name-only --diff-filter=M)
DELETED_FILES=$(git diff --cached --name-only --diff-filter=D)
FILE_COUNT=$(echo "$CHANGED_FILES" | grep -c .)

# ---- type の判定 ----
TYPE="chore"

if [[ -n "$ADDED_FILES" ]]; then
  if echo "$ADDED_FILES" | grep -q "migrations/"; then
    TYPE="chore"
  elif echo "$ADDED_FILES" | grep -q "_test\.go$"; then
    TYPE="test"
  elif echo "$ADDED_FILES" | grep -q "\.go$"; then
    TYPE="feat"
  elif echo "$ADDED_FILES" | grep -q "\.github/"; then
    TYPE="ci"
  elif echo "$ADDED_FILES" | grep -q "\.md$"; then
    TYPE="docs"
  fi
elif [[ -n "$MODIFIED_FILES" ]]; then
  if echo "$MODIFIED_FILES" | grep -qv "\.go$" && echo "$MODIFIED_FILES" | grep -q "\.md$"; then
    TYPE="docs"
  elif echo "$MODIFIED_FILES" | grep -q "\.github/"; then
    TYPE="ci"
  elif echo "$MODIFIED_FILES" | grep -q "\.go$"; then
    TYPE="fix"
  fi
fi

# ---- scope の判定 ----
FIRST_FILE=$(echo "$CHANGED_FILES" | head -1)
DIR=$(dirname "$FIRST_FILE")
SCOPE=""

if [[ "$DIR" == *"worker/internal/"* ]]; then
  SCOPE=$(echo "$DIR" | sed 's|.*worker/internal/\([^/]*\).*|\1|')
elif [[ "$DIR" == *"worker/cmd"* ]]; then
  SCOPE="worker"
elif [[ "$DIR" == *".github"* ]]; then
  SCOPE="ci"
elif [[ "$DIR" == *"supabase"* ]]; then
  SCOPE="db"
elif [[ "$DIR" == *".claude"* ]]; then
  SCOPE="claude"
fi

# ---- 説明文の生成 ----
if [[ $FILE_COUNT -eq 1 ]]; then
  BASENAME=$(basename "$FIRST_FILE")
  case "$TYPE" in
    feat)    DESCRIPTION="${BASENAME}を追加" ;;
    fix)     DESCRIPTION="${BASENAME}を修正" ;;
    docs)    DESCRIPTION="${BASENAME}を更新" ;;
    test)    DESCRIPTION="${BASENAME}にテストを追加" ;;
    ci)      DESCRIPTION="${BASENAME}を設定" ;;
    chore)   DESCRIPTION="${BASENAME}を更新" ;;
    *)       DESCRIPTION="${BASENAME}を変更" ;;
  esac
else
  case "$TYPE" in
    feat)    DESCRIPTION="${FILE_COUNT}件のファイルを追加" ;;
    fix)     DESCRIPTION="${FILE_COUNT}件のファイルを修正" ;;
    docs)    DESCRIPTION="ドキュメントを更新" ;;
    test)    DESCRIPTION="テストを追加・更新" ;;
    ci)      DESCRIPTION="CI/CD設定を更新" ;;
    chore)   DESCRIPTION="${FILE_COUNT}件の設定ファイルを更新" ;;
    *)       DESCRIPTION="${FILE_COUNT}件のファイルを変更" ;;
  esac
fi

# ---- コミットメッセージ組み立て ----
if [[ -n "$SCOPE" ]]; then
  COMMIT_MSG="${TYPE}(${SCOPE}): ${DESCRIPTION}"
else
  COMMIT_MSG="${TYPE}: ${DESCRIPTION}"
fi

git commit -m "$COMMIT_MSG"
echo "コミット完了: $COMMIT_MSG" >&2
