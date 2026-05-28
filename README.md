# e8markets

Small full-stack app integrating with the [TradeLocker Public API](https://public-api.tradelocker.com/docs/getting-started). React frontend, Go REST API, Postgres, Redis. The frontend never calls TradeLocker directly; there is no premade TradeLocker SDK.

Further detail: [docs/docs.md](docs/docs.md).

## Setup

Requires Docker, Make, and curl.

```bash
make dev    # start
make stop   # stop
make reset  # stop + wipe DB volumes
```

- **App:** http://localhost:5173
- **API:** http://localhost:8080/health
- **DB (Adminer):** http://localhost:9901/?pgsql=postgres-db&username=postgres&db=e8markets&ns=public

First run copies `.env.example` → `.env.local` in `client/` and `server/`. TradeLocker credentials are entered at login (not stored in env files).

## Docker

Infra (Postgres, Redis, Adminer):

```bash
cd server/docker && docker compose up -d --wait
```

The full stack (infra + API + client) is started with `make dev` via [`scripts/dev.sh`](scripts/dev.sh), which builds and runs `e8m-server` and `e8m-client` containers on the same network.

```bash
make dev
```

## Architecture

```
Browser (React)  →  Go API  →  TradeLocker API
                      ├── Redis (sessions)
                      └── Postgres (logs, position snapshots)
```

- **Frontend:** React, TanStack Router & Query, Vite (`client/`)
- **Backend:** Go, custom TradeLocker HTTP client (`server/internal/tradelocker`)
- **Sessions:** HttpOnly cookie + Redis; access token refresh before expiry on protected routes
- **Persistence:** `api_logs`, `sync_runs`, `position_snapshots` in Postgres

## Known limitations

- **Closed trades** — not synced; stored history is open-position snapshots at sync time.
- **Charts** — no historical price chart (optional assignment feature).
- **Docker** — `make dev` orchestrates services; not a single root `docker compose up` for everything.
- **Auto-sync** — client polls every 3 minutes when the tab is visible; no backend cron job.
- **UI** — functional, not polished.
- **Snapshot fields** — PnL / current price only when TradeLocker returns them in the positions payload.
- **No shared BE/FE types** — request/response shapes are defined separately in Go and TypeScript; an OpenAPI spec with type generation (e.g. for the client) was not added.
- **Account choosing** — not fully tested; E8 / TradeLocker test restrictions allow only one account, so multi-account selection could not be verified end-to-end.
