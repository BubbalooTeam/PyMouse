#!/usr/bin/env bash
#
# deploy.sh — install a pre-built PyMouse binary and restart the service.
#
# The repo is public, so RELEASES can be downloaded with plain curl — no token,
# no gh login. Only PR-mode (downloading a workflow artifact) needs gh auth,
# because GitHub gates workflow artifacts behind authentication even on public
# repos.
#
# Modes (auto-detected from $1):
#
#   scripts/deploy.sh                 # latest release on PyMouse (production)
#   scripts/deploy.sh <tag>           # specific release tag (e.g. 2026-08-04-1cc4e4d)
#   scripts/deploy.sh pr              # latest PR workflow artifact (needs gh auth)
#   scripts/deploy.sh <run-id>        # specific PR workflow run id (needs gh auth)
#
# Requirements:
#   - curl (always) and gh (only for `pr` / `<run-id>` modes).
#   - systemd unit "pymouse" configured with ExecStart=$BIN_PATH.
#   - Repo checked out at $INSTALL_DIR (for .env, downloads dir, fallback).
#
set -euo pipefail

REPO="BubbalooTeam/PyMouse"
INSTALL_DIR="${PYMOUSE_INSTALL_DIR:-/root/PyMouse}"
BIN_PATH="${INSTALL_DIR}/bin/PyMouse"
ARTIFACT_NAME="pymouse-linux-amd64"
ASSET_NAME="PyMouse"   # name of the binary inside the release/artifact
TMP_DIR="$(mktemp -d)"
trap 'rm -rf "$TMP_DIR"' EXIT

cd "$INSTALL_DIR"

is_number() { [[ "$1" =~ ^[0-9]+$ ]]; }

arg="${1:-}"

# Resolve the requested binary into "$TMP_DIR/$ASSET_NAME".
if [[ -z "$arg" || "$arg" == "latest" ]]; then
    # --- latest release (public, no auth) ---
    echo "==> Downloading latest release from $REPO..."
    curl -fL --retry 3 -o "$TMP_DIR/$ASSET_NAME" \
        "https://github.com/$REPO/releases/latest/download/$ASSET_NAME"

elif [[ "$arg" == "pr" ]] || is_number "$arg"; then
    # --- workflow artifact (needs gh; GitHub gates artifacts even on public repos) ---
    if ! command -v gh >/dev/null 2>&1; then
        echo "ERROR: 'gh' is required for PR/run-id mode but was not found." >&2
        echo "       Install it (https://cli.github.com) and run 'gh auth login'." >&2
        exit 1
    fi
    if [[ "$arg" == "pr" ]]; then
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
    else
        RUN_ID="$arg"
    fi
    echo "==> Downloading artifact from run $RUN_ID..."
    gh run download "$RUN_ID" --repo "$REPO" --name "$ARTIFACT_NAME" --dir "$TMP_DIR"

else
    # --- specific release tag (public, no auth) ---
    echo "==> Downloading release '$arg' from $REPO..."
    curl -fL --retry 3 -o "$TMP_DIR/$ASSET_NAME" \
        "https://github.com/$REPO/releases/download/$arg/$ASSET_NAME"
fi

if [[ ! -f "$TMP_DIR/$ASSET_NAME" ]]; then
    echo "ERROR: binary was not downloaded to $TMP_DIR/$ASSET_NAME." >&2
    exit 1
fi

echo "==> Installing to $BIN_PATH"
mkdir -p "$(dirname "$BIN_PATH")"
install -m 0755 "$TMP_DIR/$ASSET_NAME" "$BIN_PATH"

echo "==> Restarting pymouse.service"
systemctl restart pymouse

echo "==> Done. Recent logs:"
journalctl -u pymouse -n 20 --no-pager
