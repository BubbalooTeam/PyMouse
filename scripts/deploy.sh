#!/usr/bin/env bash
#
# deploy.sh — build the bot from the local checkout and restart the service.
#
# This is the DEFAULT deploy path: compile from source with the Go toolchain
# already on the host. Works fully offline (after `git pull`), works even if
# CI is broken, and guarantees the running binary matches the checked-out
# commit.
#
# To deploy a pre-built binary from GitHub instead (latest release, a specific
# tag, or a PR artifact), see scripts/deploy-release.sh.
#
# Usage:
#   scripts/deploy.sh             # build + install + restart
#
# What it does:
#   1. go clean -cache            # drop the build cache (avoids stale-object
#                                  bugs; the Go toolchain will rebuild). We do
#                                  NOT pass -modcache, so downloaded modules
#                                  are reused and there is no re-download.
#   2. go build (CGO disabled, -trimpath, -ldflags="-s -w") -> bin/PyMouse
#   3. systemctl restart pymouse
#
set -euo pipefail

INSTALL_DIR="${PYMOUSE_INSTALL_DIR:-/root/PyMouse}"
BIN_PATH="${INSTALL_DIR}/bin/PyMouse"
GO_BIN="${GO_BIN:-go}"

cd "$INSTALL_DIR"

if ! command -v "$GO_BIN" >/dev/null 2>&1; then
    echo "ERROR: 'go' not found on PATH (set GO_BIN=/usr/local/go/bin/go if needed)." >&2
    exit 1
fi

for command_name in node npm; do
    if ! command -v "$command_name" >/dev/null 2>&1; then
        echo "ERROR: '$command_name' is required by the YouTube extractor." >&2
        exit 1
    fi
done

echo "==> node: $(node --version); npm: $(npm --version)"

echo "==> HEAD: $(git -C "$INSTALL_DIR" rev-parse --short HEAD) $(git -C "$INSTALL_DIR" log -1 --format=%s)"

echo "==> Cleaning build cache (keeping module cache)..."
"$GO_BIN" clean -cache

echo "==> Building PyMouse (linux/amd64, CGO disabled)..."
mkdir -p "$(dirname "$BIN_PATH")"
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    "$GO_BIN" build -trimpath -ldflags="-s -w" -o "$BIN_PATH" .

echo "==> Verifying binary..."
file "$BIN_PATH"
test -x "$BIN_PATH"

echo "==> Restarting pymouse.service"
systemctl restart pymouse

echo "==> Done. Recent logs:"
journalctl -u pymouse -n 20 --no-pager
