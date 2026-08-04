# Deploy

The bot ships as a single pre-built binary. Assets (fonts, icons, locales) are
embedded into the executable via `go:embed`, so the server does **not** need the
`pymouse/assets/` or `locales/` trees on disk — only the binary, the `.env`
file, and a writable directory for downloads.

The repo checkout is still kept on the server as a fallback (so you can `go run`
in an emergency) and because the bot reads `.env`, `youtubeCookies.txt`, and
writes to `pymouse/downloads/` relative to its working directory.

## Pipeline

1. Push to `PyMouse` (or merge a PR) → `.github/workflows/build.yml` builds a
   `linux/amd64` binary with `CGO_ENABLED=0` and publishes it as the artifact
   `pymouse-linux-amd64`.
2. On the server, `scripts/deploy.sh` downloads the latest successful artifact,
   installs it to `/root/PyMouse/bin/pymouse`, and restarts the service.

Deploy is manual on purpose — merging a PR does **not** automatically update
production. Run `scripts/deploy.sh` when you are ready.

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
  there is nothing to compile at startup anymore, so restarts are instant.
- `EnvironmentFile` points at `.env` (the old unit pointed at `PyMouse.log`,
  which was incorrect).

### 2. GitHub CLI

`scripts/deploy.sh` uses `gh` to download artifacts. Authenticate once:

```bash
gh auth login
```

Choose GitHub.com → HTTPS, and paste a personal access token with `read`
access to `BubbalooTeam/PyMouse` (workflow `contents: read` + `actions: read`).

## Deploying

```bash
cd /root/PyMouse
git pull                              # keep the checkout in sync (for .env, fallback)
scripts/deploy.sh                     # pull the latest build + restart
```

To deploy a specific workflow run: `scripts/deploy.sh <run-id>`.

## Rollback / fallback

If the new binary misbehaves, you have two options:

- **Previous binary**: re-run `scripts/deploy.sh <older-run-id>`.
- **`go run` fallback** (the old way): temporarily point the unit back at
  `go run` — the repo checkout is still on disk.

  ```bash
  # edit /etc/systemd/system/pymouse.service ExecStart to:
  #   ExecStart=/usr/local/go/bin/go run /root/PyMouse/main.go
  systemctl daemon-reload
  systemctl restart pymouse
  ```
