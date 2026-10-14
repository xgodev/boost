#!/bin/bash
# PostToolUse (Edit|Write): lembrete da skill boost-claude-sync ao tocar
# arquivo de componente (Iron Law #6). Um aviso por sessão.

input=$(cat)
file_path=$(printf '%s' "$input" | jq -r '.tool_input.file_path // empty')
session_id=$(printf '%s' "$input" | jq -r '.session_id // "nosession"')
[ -z "$file_path" ] && exit 0

repo_root="${CLAUDE_PROJECT_DIR:-$(cd "$(dirname "$0")/../.." && pwd)}"

case "$file_path" in
  "$repo_root"/factory/* | \
  "$repo_root"/wrapper/* | \
  "$repo_root"/bootstrap/* | \
  "$repo_root"/extra/* | \
  "$repo_root"/fx/* | \
  "$repo_root"/model/* | \
  "$repo_root"/config.go | \
  "$repo_root"/start.go) ;;
  *) exit 0 ;;
esac

marker="${TMPDIR:-/tmp}/boost-claude-sync-reminder-${session_id}"
[ -f "$marker" ] && exit 0
: > "$marker"

cat <<'EOF'
{"hookSpecificOutput":{"hookEventName":"PostToolUse","additionalContext":"Componente do boost tocado nesta sessão. Antes de concluir/push: se a mudança for observável (API pública, config key, default, comportamento, README), siga a skill de projeto boost-claude-sync (.claude/skills/boost-claude-sync/SKILL.md) para sincronizar a reference em xgodev/boost-claude — Iron Law #6."}}
EOF
