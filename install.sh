#!/usr/bin/env bash
set -euo pipefail

REPO="https://github.com/colchuck-ai/intent.git"
BRANCH="main"
DEST="$HOME/.cursor/skills/intent"
TMP="$(mktemp -d)"

trap 'rm -rf "$TMP"' EXIT

git clone --depth 1 --branch "$BRANCH" "$REPO" "$TMP/intent"

rm -rf "$DEST"
mkdir -p "$DEST"
cp -R "$TMP/intent/.cursor/skills/intent/" "$DEST/"

echo "Installed intent skill to $DEST"
