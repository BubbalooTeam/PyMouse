# Deploy

The bot ships as a single pre-built binary. Assets (fonts, icons, locales) are
embedded into the executable via `go:embed`, so the server needs only the
binary, the `.env` file, and a writable directory for downloads — **not** the
`pymouse/assets/` or `locales/` trees.

The repo is **public**, so production binaries (GitHub Releases) download with
plain `curl` — no token, no `gh` login. Only pre-merge PR testing needs `gh`,
because GitHub gates workflow *artifacts* behind authentication even on public
repos.

The repo checkout is still kept on the server as a fallback (so you can `go
run` in an emergency) and because the bot reads `.env`, `youtubeCookies.txt`,
and writes to `pymouse/downloads/` relative to its working directory.

## How binaries are published

| Trigger | What gets published | Auth to download | Lifetime |
|---|---|---|---|
| **Pull request** (any branch) | Workflow artifact `pymouse-linux-amd64` | `gh` login needed | 30 days |
| **Push to `PyMouse`** (main) | **GitHub Release**, tagged `<YYYY-MM-DD>-<short-sha>` (e.g. `2026-08-04-1cc4e4d`), marked `latest` | **None** (public URL) | **Permanent** |

Releases do not expire, so every push to main is a distinct, rollback-able
version. PR artifacts are for pre-merge testing only.

## One-time server setup

### 1. systemd unit

Create `/etc/systemd/system/pymouse.service`:

```ini
[Unit]
Description=PyMouse Bot
After=network.target

[Service]
User=root
WorkingDirectory=/root/PyMouse
ExecStart=/root/PyMouse/bin/pymouse
Restart=always
EnvironmentFile=/root/PyMouse/.env

[Install]
WantedBy=multi-user.target
```

Then:

```bash
systemctl daemon-reload
systemctl enable pymouse
```

Key changes vs. the old unit:
- `ExecStart` runs the pre-built binary instead of `go run`.
- The `ExecStartPre=go clean -cache -modcache -i -r` line is **removed** —
  there is nothing to compile at startup, so restarts are instant.
- `EnvironmentFile` points at `.env` (the old unit pointed at `PyMouse.log`,
  which was incorrect).

### 2. GitHub CLI (only if you'll use `deploy.sh pr`)

`deploy.sh` uses plain `curl` for releases — no auth needed. The only mode that
needs `gh` is `deploy.sh pr` (downloading a PR workflow artifact). If you want
to test PRs on the server before merging, authenticate once:

```bash
gh auth login
```

Choose GitHub.com → HTTPS, and authenticate with access to
`BubbalooTeam/PyMouse`. If you never use `pr` mode, you can skip this entirely.

## Deploying

`scripts/deploy.sh` auto-detects what to install based on its argument:

| Command | Installs | Auth |
|---|---|---|
| `scripts/deploy.sh` | latest release on PyMouse (production default) | none |
| `scripts/deploy.sh <tag>` | specific release, e.g. `2026-08-04-1cc4e4d` | none |
| `scripts/deploy.sh pr` | artifact from the latest successful PR build | `gh` login |
| `scripts/deploy.sh <run-id>` | a specific PR workflow run (numeric) | `gh` login |

Typical flows:

```bash
# Production update after merging a PR (no auth needed):
cd /root/PyMouse && git pull && scripts/deploy.sh

# Roll back to a known-good version (no auth):
scripts/deploy.sh 2026-08-04-1cc4e4d

# Test a PR before merging (needs gh login):
scripts/deploy.sh pr
```

`scripts/deploy.sh` always restarts the `pymouse` service and tails the last
20 journal lines.

## Rollback / fallback

If a new binary misbehaves:

- **Previous version**: `scripts/deploy.sh <older-tag>` — releases are
  permanent, so any past main build is still available.
- **`go run` fallback** (the original method): the repo checkout is still on
  disk. Temporarily point the unit back at `go run`:

  ```bash
  # edit /etc/systemd/system/pymouse.service ExecStart to:
  #   ExecStart=/usr/local/go/bin/go run /root/PyMouse/main.go
  systemctl daemon-reload
  systemctl restart pymouse
  ```
