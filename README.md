# Marketing Backend — Greenpark Qualified Demand Control Tower API

A small, dependency-free Go HTTP API that serves the data for the **Dashboard
Marketing Greenpark** — the CEO "Qualified Demand Control Tower" that tracks the
full demand funnel from Impression to Cash-In, lead quality, channel & project
performance, digital assets, content, the CEO command panel and the alert system.

It is built with a clean, layered architecture so the file-backed data source can
later be swapped for real sources (ad platforms, CRM, analytics) without touching
the transport or business logic. Master data is **editable** (login + admin CRUD)
and persisted to a JSON file, so input flows straight back into the dashboard reads.

## Architecture

```
cmd/server            composition root — wires everything and runs the server
internal/
  config              env-based runtime configuration (with defaults)
  domain              core entities (FunnelStage, KPI, Channel, Project, …) — no deps
  passwd              salted SHA-256 password hashing helpers
  auth                token login + in-memory bearer-token sessions
  repository          storage boundary (interface) + file/Postgres seeded store (CRUD + users)
  service             business logic — composes data, derives the summary, write use-cases
  transport/http      router, handlers (auth + reads + writes), middleware, JSON helpers
```

Dependency direction points inward: `transport → service → repository → domain`.
Each layer depends only on the interfaces of the one beneath it.

## Auth & roles

- `POST /api/auth/login` → `{ token, user }`. Send `Authorization: Bearer <token>` on every other call.
- **admin** (`admin / admin123`) — full master-data write access.
- **viewer** (`viewer / viewer123`) — read-only; writes return `403`.

Master-data edits (`POST`/`PUT`/`DELETE`) persist to the store and immediately
change the derived `summary` returned by the dashboard.

## Run

```bash
cd backend/marketing
go run ./cmd/server
# marketing API listening on http://localhost:8086
```

Configuration via environment variables:

| Variable                   | Default                     | Description                                   |
| -------------------------- | --------------------------- | --------------------------------------------- |
| `MARKETING_PORT`           | `8086`                      | HTTP port                                     |
| `MARKETING_ALLOW_ORIGIN`   | `*`                         | CORS allowed origin                           |
| `MARKETING_DATA_PATH`      | `data/marketing-data.json`  | JSON file the master data persists to         |
| `MARKETING_DATABASE_URL`   | _(empty)_                   | PostgreSQL DSN; when set, used instead of file |

## API

All responses are JSON. Read-only `GET` endpoints under `/api` (require auth):

| Method · Path                 | Description                                  |
| ----------------------------- | -------------------------------------------- |
| `GET /api/health`             | Liveness probe (public)                      |
| `GET /api/dashboard`          | Full payload (all sections + derived summary)|
| `GET /api/summary`            | Executive summary (derived)                  |
| `GET /api/context`            | Header context (period, goal, booking YTD)   |
| `GET /api/funnel`             | Full demand funnel (Impression → Cash-In)    |
| `GET /api/kpis`               | North Star KPI ribbon                        |
| `GET /api/lead-quality`       | Lead quality & MQL scoring                   |
| `GET /api/handover`           | MQL → SAL handover metrics                   |
| `GET /api/channels`           | Channel performance matrix                   |
| `GET /api/projects`           | Project demand & readiness                   |
| `GET /api/projects/{name}`    | Single project (404 if unknown)              |
| `GET /api/assets`             | Digital asset registry                       |
| `GET /api/ig-accounts`        | Project Instagram accounts                   |
| `GET /api/content`            | Content & winning campaigns                  |
| `GET /api/commands`           | CEO command panel                            |
| `GET /api/alerts`             | Red / Yellow / Green alert system            |
| `GET /api/reason-codes`       | Funnel leakage reason codes                  |

Admin writes: `PUT` for singletons (`context`, `spend`, `lead-quality`, `content`,
`alerts`) and `POST`/`DELETE /{id}` for the collections (`funnel`, `kpis`,
`handover`, `channels`, `projects`, `assets`, `ig-accounts`, `commands`,
`reason-codes`). `POST` with empty `_id` creates; with an existing `_id` updates.

The `summary` is **derived** in the service layer from the context, funnel,
channels, commands and alerts (progress to goal, total leads/MQL/spend, booking,
red-alert count, open commands).

## Test

```bash
go build ./...
go vet ./...
```
