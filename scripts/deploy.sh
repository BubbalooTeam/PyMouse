#!/usr/bin/env bash
#
# deploy.sh — install a pre-built PyMouse binary and restart the service.
#
# Three modes (the script auto-detects which based on $1):
#
#   scripts/deploy.sh                 # latest release on PyMouse (production)
#   scripts/deploy.sh pr              # workflow artifact from the latest PR run
#   scripts/deploy.sh <tag>           # a specific release tag (e.g. 2026-08-04-1cc4e4d)
#   scripts/deploy.sh <run-id>        # a specific PR workflow run id (numeric)
#
# Requirements on the server:
#   - GitHub CLI (gh) installed and authenticated (gh auth login) with read
#     access to BubbalooTeam/PyMouse.
#   - The systemd unit "pymouse" configured to run ExecStart=$BIN_PATH.
#   - The repo checked out at $INSTALL_DIR (kept for .env, youtubeCookies.txt,
#     the downloads dir, and as a `go run` fallback).
#
set -euo pipefail

REPO="BubbalooTeam/PyMouse"
INSTALL_DIR="${PYMOUSE_INSTALL_DIR:-/root/PyMouse}"
BIN_PATH="${INSTALL_DIR}/bin/pymouse"
ARTIFACT_NAME="pymouse-linux-amd64"
TMP_DIR="$(mktemp -d)"
trap 'rm -rf "$TMP_DIR"' EXIT

cd "$INSTALL_DIR"

# Resolve the requested binary into "$TMP_DIR/pymouse".
arg="${1:-}"
downloaded=0

is_number() { [[ "$1" =~ ^[0-9]+$ ]]; }

if [[ -z "$arg" || "$arg" == "latest" ]]; then
    # --- latest release on PyMouse (the production default) ---
    echo "==> Downloading latest release from $REPO..."
    gh release download --repo "$REPO" --pattern pymouse --dir "$TMP_DIR" --clobber
    downloaded=1

elif [[ "$arg" == "pr" ]]; then
    # --- latest PR workflow artifact (for pre-merge testing) ---
    echo "==> Locating latest successful PR build run..."
    RUN_ID="$(gh run list \
        --repo "$REPO" \
        --workflow build.yml \
        --event pull_request \
        --status success \
        --limit 1 \
        --json databaseId \
        --jq '.[0].databaseId')"
    if [[ -z "$RUN_ID" ]]; then
        echo "ERROR: no successful PR build found. Push to the PR branch first." >&2
        exit 1
    fi
    echo "==> Using PR run $RUN_ID"
    gh run download "$RUN_ID" --repo "$REPO" --name "$ARTIFACT_NAME" --dir "$TMP_DIR"
    downloaded=1

elif is_number "$arg"; then
    # --- numeric: a specific workflow run id (artifact) ---
    echo "==> Downloading artifact from run $arg..."
    gh run download "$arg" --repo "$REPO" --name "$ARTIFACT_NAME" --dir "$TMP_DIR"
    downloaded=1

else
    # --- otherwise: treat $arg as a release tag ---
    echo "==> Downloading release '$arg' from $REPO..."
    gh release download "$arg" --repo "$REPO" --pattern pymouse --dir "$TMP_DIR" --clobber
    downloaded=1
fi

if [[ "$downloaded" -ne 1 ]] || [[ ! -f "$TMP_DIR/pymouse" ]]; then
    echo "ERROR: binary was not downloaded to $TMP_DIR/pymouse." >&2
    exit 1
fi

echo "==> Installing to $BIN_PATH"
mkdir -p "$(dirname "$BIN_PATH")"
install -m 0755 "$TMP_DIR/pymouse" "$BIN_PATH"

echo "==> Restarting pymouse.service"
systemctl restart pymouse

echo "==> Done. Recent logs:"
journalctl -u pymouse -n 20 --no-pager
