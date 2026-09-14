# Backend

Go modular-monolith REST API. See [ADR-0003](../adrs/0003-modular-monolith-backend.md), [ADR-0004](../adrs/0004-backend-language-go.md), [ADR-0005](../adrs/0005-backend-framework.md), [ADR-0006](../adrs/0006-api-style-rest.md).

## Layout (planned)
```
backend/
  cmd/api/                entrypoint (main.go)
  internal/
    auth/                 auth module: handler, service, repository
    product/               product module
    review/                review module
    media/                 media/upload module
    platform/              shared infra: db, cache, logging, middleware
  migrations/              golang-migrate SQL files (ADR-0009)
  api/openapi.yaml          API contract — source of truth for all clients
```

## Local development
See [`docs/development-setup.md`](../docs/development-setup.md).

## Testing
See [ADR-0022](../adrs/0022-testing-strategy.md). `go test ./...` for unit tests, `go test -tags=integration ./...` for integration tests against real Postgres/Redis via testcontainers.

## Conventions
See [`docs/coding-conventions.md`](../docs/coding-conventions.md).

## MVP scope (current state)

This is a first, intentionally small implementation for testing the core scan → view → submit loop end to end. It deliberately diverges from the documents above in a few ways — flagging them here rather than silently diverging, per `docs/coding-conventions.md`:

- **No auth.** `/auth/*` from `backend/api/openapi.yaml` isn't implemented. `POST /v1/products/{productId}/reviews` has no `bearerAuth` requirement.
- **No `users`/`identities` tables.** `reviews` has a plain `author_name text` column instead of a `user_id` FK — there's no account system yet to attach reviews to. Revisit when [ADR-0014](../adrs/0014-authentication-strategy.md) is actually implemented (will need a follow-up migration).
- **No `rating_summaries` table.** `GET /v1/products/{barcode}` computes `COUNT`/`AVG` over `reviews` directly per request — fine at MVP scale; the schema doc itself notes this denormalization is "not yet an ADR-level decision."
- **No `sqlc` codegen.** Repositories use plain parameterized `database/sql` + the `pgx` stdlib driver.
- **No `openapi-generator` client codegen.** The Android app hand-writes its Retrofit interface/DTOs against the actual endpoints below rather than generating from the (fuller) `openapi.yaml`.
- **No Redis.** `docker-compose.yml` here only runs Postgres.

### Endpoints implemented
- `GET /v1/products/{barcode}` — looks up a product by barcode; creates a placeholder row if it doesn't exist yet. Returns the product plus an inline-computed rating summary.
- `GET /v1/products/{productId}/reviews` — newest-first, up to 50 reviews (no cursor pagination yet).
- `POST /v1/products/{productId}/reviews` — body `{author_name?, body, rating?}`, no auth.

### Quick start
```bash
cd backend
cp .env.example .env.local   # then adjust PORT/DATABASE_URL if needed
docker compose up -d          # starts Postgres on localhost:55432 (5432 was taken locally)
docker compose exec -T postgres psql -U everyreview -d everyreview < migrations/0001_init.up.sql
DATABASE_URL="postgres://everyreview:everyreview@localhost:55432/everyreview?sslmode=disable" go run ./cmd/api
```
The API listens on `localhost:8080` by default (override with `PORT`). `go test ./...` runs the unit tests (fake in-memory repositories — no Postgres integration tests yet, see ADR-0022 for the fuller plan to add later).
