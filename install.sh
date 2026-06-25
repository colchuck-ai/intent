#!/usr/bin/env bash
set -euo pipefail

REPO="https://github.com/colchuck-ai/intent.git"
BRANCH="main"
AGENT="cursor"
DEST=""

usage() {
  cat <<'EOF'
Install the intent Agent Skill.

Usage: install.sh [--agent <name>] [DIR]

Options:
  --agent <name>   Install into a known agent's skills directory:
                     cursor   ~/.cursor/skills/intent   (default)
                     claude   ~/.claude/skills/intent
                     agents   ~/.agents/skills/intent
  -h, --help       Show this help.

Arguments:
  DIR              Install directly into DIR, overriding --agent.

Examples:
  install.sh                                   # ~/.cursor/skills/intent
  install.sh --agent claude                    # ~/.claude/skills/intent
  install.sh ~/.config/my-agent/skills/intent  # custom directory

When piping from curl, pass options after `-s --`:
  curl -fsSL <url>/install.sh | bash -s -- --agent claude
EOF
}

while [ $# -gt 0 ]; do
  case "$1" in
    --agent) AGENT="${2:-}"; shift 2 ;;
    --agent=*) AGENT="${1#*=}"; shift ;;
    -h|--help) usage; exit 0 ;;
    --) shift; break ;;
    -*) echo "Unknown option: $1" >&2; usage >&2; exit 1 ;;
    *) DEST="$1"; shift ;;
  esac
done

# A directory may also follow `--`.
if [ -z "$DEST" ] && [ $# -gt 0 ]; then
  DEST="$1"
fi

# Resolve the destination from the agent when no explicit directory is given.
if [ -z "$DEST" ]; then
  case "$AGENT" in
    cursor) DEST="$HOME/.cursor/skills/intent" ;;
    claude) DEST="$HOME/.claude/skills/intent" ;;
    agents) DEST="$HOME/.agents/skills/intent" ;;
    *) echo "Unknown agent: $AGENT" >&2; usage >&2; exit 1 ;;
  esac
fi

TMP="$(mktemp -d)"
trap 'rm -rf "$TMP"' EXIT

git clone --depth 1 --branch "$BRANCH" "$REPO" "$TMP/intent"

rm -rf "$DEST"
mkdir -p "$DEST"
cp -R "$TMP/intent/intent/" "$DEST/"

echo "Installed intent skill to $DEST"
