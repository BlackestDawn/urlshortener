# urlshortener

[![CI](https://github.com/BlackestDawn/urlshortener/actions/workflows/pr-checks.yml/badge.svg)](https://github.com/BlackestDawn/urlshortener/actions/workflows/pr-checks.yml)
[![Coverage](https://raw.githubusercontent.com/BlackestDawn/urlshortener/badges/.badges/main/coverage.svg)](https://github.com/BlackestDawn/urlshortener/actions/workflows/coverage-baseline.yml)
![Go version](https://img.shields.io/badge/go-1.26-00ADD8?logo=go)
![License](https://img.shields.io/badge/license-CC0-lightgrey)

A production-shaped URL shortener written in Go — built less to be another "shorten a link" toy and more as a demonstration of how I structure, test, and ship a real backend service: clean architecture, an actual CI/CD pipeline, and infrastructure-as-code for a deployed environment.

**Live**: [shortener.alexstauch.app](https://shortener.alexstauch.app) ([API](#api) below — it's a JSON API, not a webpage)

## Why this project

Most URL-shortener portfolio projects stop at "it works locally." This one is built the way I'd approach a service at work:

- **Layered architecture** — domain logic has zero dependency on Gin or Postgres, so the core business rules are testable and swappable in isolation.
- **Deterministic short codes** — codes are derived from a BLAKE2b hash of the target URL rather than randomly generated, so re-shortening the same URL is idempotent by construction, with no collision-retry logic needed.
- **Real test pyramid** — unit tests against mocked interfaces, plus integration tests that spin up a real Postgres instance, plus a Docker smoke test in CI that boots the actual container image and hits it over HTTP.
- **CI/CD that deploys** — every merge to `main` builds, tests, and ships to a live Cloud Run staging environment automatically; production is a one-click, promote-the-exact-tested-image workflow, not a re-run of `go build` and hope.
- **Infrastructure as code** — the GCP project (Artifact Registry, IAM, Workload Identity Federation, Secret Manager) and the Postgres environments (via Neon, with a copy-on-write staging branch) are both defined in Terraform, not clicked together in a console.

## Architecture

```
HTTP request
     │
     ▼
┌─────────────────────────┐
│  cmd/server              │  Gin router + middleware (request ID, structured
│  (controllers, handlers) │  logging, panic recovery, body-size cap, rate
│                          │  limiting, timeout, centralized error mapping)
└────────────┬─────────────┘
             │ depends on
             ▼
┌─────────────────────────┐
│  internal/service        │  Business use cases (Shorten, Resolve, Stats,
│  (IShorten)               │  Delete) — depends only on a repository interface
└────────────┬─────────────┘
             │ depends on
             ▼
┌─────────────────────────┐
│  internal/domain         │  Framework-free core: ShortUrl type, URL
│  (IRepository interface)│  validation, code generation, sentinel errors
└────────────┬─────────────┘
             │ implemented by
             ▼
┌─────────────────────────┐
│  internal/repository      │  Postgres adapter — sqlc-generated typesafe
│  (PostgresRepository)    │  queries over jackc/pgx, explicit pool tuning
└─────────────────────────┘
```

Dependencies point inward: the domain layer knows nothing about HTTP or SQL, `main.go` wires the concrete Postgres repository into the service, and the service into the HTTP controller. Every layer boundary is an interface (`domain.IRepository`, `service.IShorten`), so each is independently unit-testable with generated mocks — no live database required except for the integration suite.

## API

| Method | Path | Description |
|---|---|---|
| `GET` | `/healthz` | Health check (own tighter rate limit) |
| `POST` | `/api/v1/links` | Shorten a URL |
| `GET` | `/api/v1/links/:code` | Look up the original URL for a code |
| `GET` | `/api/v1/links/:code/stats` | Click count and metadata for a code |
| `DELETE` | `/api/v1/links/:code` | Delete a short URL |
| `GET` | `/:code` | 308 redirect to the original URL, increments click count |

```bash
curl -X POST localhost:8080/api/v1/links \
  -H 'Content-Type: application/json' \
  -d '{"url": "https://example.com/some/very/long/path"}'
# → {"shortenedUrl":"https://short.example/9f8a1c3d2b7e04af"}
```

## Engineering details worth a closer look

- **Search-ready schema**: the `original_url` column is backed by a Postgres `pg_trgm` GIN index, added in a dedicated migration, for fast partial-text lookups as the dataset grows.
- **Tuned connection pool**: explicit `MaxOpenConns`/`MaxIdleConns`/`ConnMaxLifetime` on the `database/sql` pool rather than relying on driver defaults.
- **Graceful shutdown**: SIGINT/SIGTERM handling with a bounded shutdown context and resource cleanup, not a bare process kill.
- **Iterative performance work**: recent commits specifically tune the click-increment and redirect-resolve hot paths and connection limits — the kind of follow-up pass a project usually only gets post-launch.

## Testing

```
make test           # staticcheck + go vet + unit tests (mocked dependencies)
make integrationtest # real Postgres via testcontainers — repository + full HTTP API
make lint            # go vet, golangci-lint, sqlc vet
```

CI runs all of the above on every PR, plus a Docker smoke test that builds the production image and verifies `/healthz` over the network, plus a coverage-diff check against `main` so coverage can't silently regress.

## Deployment

```
merge to main ──▶ build & push image ──▶ migrate staging DB ──▶ deploy Cloud Run staging
                                                                        │
                                                            (manual approval)
                                                                        ▼
                                                    promote the *same* image ──▶ migrate prod DB ──▶ deploy Cloud Run prod
```

- **Compute**: Google Cloud Run (scale-to-zero, pay-per-request).
- **Database**: [Neon](https://neon.tech) serverless Postgres — prod and a copy-on-write staging branch.
- **Auth**: GitHub Actions authenticates to GCP via Workload Identity Federation — no long-lived service account keys.
- **Provisioning**: everything above is Terraform (`deploy/terraform/gcp`, `deploy/terraform/neon`), not manual console setup.

## Running locally

```bash
docker compose up   # postgres + migrations + app on :8080
```

or without Docker:

```bash
cp .env.example .env   # set DATABASE_URL
make migrate-up
make run
```

## Tech stack

Go 1.26 · Gin · sqlc · pgx · Postgres (Neon) · Docker · Terraform · GitHub Actions · Google Cloud Run
