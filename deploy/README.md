# Deploy

The bot ships as a single self-contained binary (assets embedded via
`go:embed`, so the server does not need `pymouse/assets/` or `locales/` on
disk). There are two deploy scripts:

| Script | Default for | Source of the binary |
|---|---|---|
| **`scripts/deploy.sh`** | day-to-day updates | **Compiled locally** from the checked-out source |
| `scripts/deploy-release.sh` | rollbacks, no-toolchain hosts | **Downloaded** from a GitHub release or PR artifact |

## Why local compile is the default

The server already has the Go toolchain and the repo checkout, and the build
cache is now **warm** (the old unit wiped it on every start with
`go clean -cache -modcache`; that line is gone). An incremental local build
takes seconds and has real advantages over pulling a remote binary:

- Works fully offline (after `git pull`).
- Works even if CI is broken (no dependency on a successful workflow run).
- The running binary is guaranteed to match the checked-out commit.
- One fewer moving part (no network, no asset naming, no CI timing).

The YouTube extractor also requires `yt-dlp`, Node.js, and npm on the server.
The deploy script checks all three before replacing the running binary.

The CI (`.github/workflows/build.yml`) still publishes a **permanent GitHub
Release** on every push to `PyMouse`, so `deploy-release.sh` is available for
rollbacks and future multi-server setups — it's just not the main path.

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
ExecStart=/root/PyMouse/bin/PyMouse
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
  compilation is now a deliberate `deploy.sh` step, not something that happens
  on every restart. Service restarts are instant.
- `EnvironmentFile` points at `.env` (the old unit pointed at `PyMouse.log`,
  which was incorrect).

### 2. GitHub CLI (only for `deploy-release.sh pr`)

`deploy-release.sh` uses plain `curl` for releases — no auth needed. The only
mode that needs `gh` is `deploy-release.sh pr` (downloading a PR workflow
artifact). If you want to test PRs on the server before merging:

```bash
gh auth login
```

Skip this entirely if you never use `pr` mode.

## Deploying

### Day-to-day (local compile)

```bash
cd /root/PyMouse
git pull                              # get the new source
scripts/deploy.sh                     # build + install + restart
```

`deploy.sh` drops the Go build cache (keeping the module cache, so no
re-download), builds `bin/PyMouse` with `CGO_ENABLED=0 -trimpath -ldflags=-s`,
and restarts the service.

### From a release (rollback / no-toolchain)

```bash
scripts/deploy-release.sh                 # latest release on PyMouse
scripts/deploy-release.sh <tag>           # specific release, e.g. 2026-08-04-1cc4e4d
scripts/deploy-release.sh pr              # latest PR workflow artifact (needs gh)
scripts/deploy-release.sh <run-id>        # specific PR workflow run (needs gh)
```

## How releases are published (for `deploy-release.sh`)

| Trigger | What gets published | Auth to download | Lifetime |
|---|---|---|---|
| **Pull request** (any branch) | Workflow artifact `pymouse-linux-amd64` | `gh` login needed | 30 days |
| **Push to `PyMouse`** (main) | **GitHub Release**, tagged `<YYYY-MM-DD>-<short-sha>`, marked `latest` | **None** (public URL) | **Permanent** |

## Rollback / fallback

If a new build misbehaves:

- **Previous version (release):** `scripts/deploy-release.sh <older-tag>` —
  releases are permanent, so any past main build is still available.
- **Previous version (source):** `git checkout <older-sha> && scripts/deploy.sh`.
- **`go run` fallback** (the original method, pre-binary): the repo checkout is
  still on disk. Temporarily point the unit back at `go run`:

  ```bash
  # edit /etc/systemd/system/pymouse.service ExecStart to:
  #   ExecStart=/usr/local/go/bin/go run /root/PyMouse/main.go
  systemctl daemon-reload
  systemctl restart pymouse
  ```
