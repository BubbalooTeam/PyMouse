#!/usr/bin/env bash
#
# deploy.sh — pull the latest pre-built PyMouse binary from GitHub Actions and
# restart the systemd service on this host.
#
# Requirements on the server:
#   - GitHub CLI (gh) installed and authenticated (gh auth login) with read
#     access to BubbalooTeam/PyMouse.
#   - The systemd unit "pymouse" configured to run ExecStart=/root/PyMouse/bin/pymouse
#     (see deploy/README.md or the PR description for the unit file).
#   - /root/PyMouse checked out (kept as a fallback to `go run` and to host
#     .env / youtubeCookies.txt / the downloads dir).
#
# Usage:
#   scripts/deploy.sh            # deploy latest run on PyMouse (default)
#   scripts/deploy.sh <run-id>   # deploy a specific workflow run
#
set -euo pipefail

REPO="BubbalooTeam/PyMouse"
BRANCH="${PYMOUSE_DEPLOY_BRANCH:-PyMouse}"
INSTALL_DIR="${PYMOUSE_INSTALL_DIR:-/root/PyMouse}"
BIN_PATH="${INSTALL_DIR}/bin/pymouse"
ARTIFACT_NAME="pymouse-linux-amd64"
TMP_DIR="$(mktemp -d)"
trap 'rm -rf "$TMP_DIR"' EXIT

cd "$INSTALL_DIR"

echo "==> Fetching latest successful build on '$BRANCH'..."
if [[ -n "${1:-}" ]]; then
  RUN_ID="$1"
else
  RUN_ID="$(gh run list \
    --repo "$REPO" \
    --branch "$BRANCH" \
    --workflow build.yml \
    --status success \
    --limit 1 \
    --json databaseId \
    --jq '.[0].databaseId')"
fi

if [[ -z "$RUN_ID" ]]; then
  echo "ERROR: no successful build found on '$BRANCH'. Trigger a push first." >&2
  exit 1
fi

echo "==> Using run $RUN_ID"
echo "==> Downloading artifact '$ARTIFACT_NAME'..."
gh run download "$RUN_ID" \
  --repo "$REPO" \
  --name "$ARTIFACT_NAME" \
  --dir "$TMP_DIR"

if [[ ! -f "$TMP_DIR/pymouse" ]]; then
  echo "ERROR: artifact did not contain 'pymouse' binary." >&2
  exit 1
fi

echo "==> Installing to $BIN_PATH"
mkdir -p "$(dirname "$BIN_PATH")"
install -m 0755 "$TMP_DIR/pymouse" "$BIN_PATH"

echo "==> Restarting pymouse.service"
systemctl restart pymouse

echo "==> Done. Tailing logs (Ctrl-C to exit)..."
journalctl -u pymouse -n 20 --no-pager
