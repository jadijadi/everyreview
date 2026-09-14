# Development Setup

## Prerequisites
- Go 1.22+
- Docker + Docker Compose (Postgres, Redis locally)
- Android Studio (latest stable) + JDK 17 for the Android app
- `golang-migrate` CLI, `sqlc` CLI
- AWS CLI configured (only needed for infrastructure work, not day-to-day backend/app dev)
- Terraform 1.x (only needed for infrastructure work)

## Backend
```bash
cd backend
cp .env.example .env.local        # dev-only dummy secrets, see ADR-0020
docker compose up -d              # starts local Postgres + Redis
migrate -path migrations -database "$DATABASE_URL" up
go run ./cmd/api
```
The API listens on `localhost:8080` by default. `GET /v1/products/0000000000000` against a freshly migrated DB should return a newly created placeholder product.

Run tests:
```bash
go test ./...                     # unit tests
go test -tags=integration ./...   # integration tests, spins up testcontainers (needs Docker)
```

## Android
1. Open `android/` in Android Studio.
2. Set `API_BASE_URL` in `local.properties` to point at your local backend (`http://10.0.2.2:8080/v1` from the emulator) or the shared `dev` environment.
3. Run on an emulator with Google Play services (required for ML Kit barcode scanning, see [ADR-0025](../adrs/0025-barcode-scanning-library.md)) or a physical device.

## Regenerating API client code
Both the Android network layer and the backend's contract tests are generated/validated from `backend/api/openapi.yaml`. After changing the spec:
```bash
scripts/generate-api-clients.sh
```

## Infrastructure (only if you're changing `infrastructure/`)
```bash
cd infrastructure/environments/dev
terraform init
terraform plan
```
Never run `terraform apply` locally against `staging` or `prod` — those apply only via CI ([ADR-0021](../adrs/0021-cicd-github-actions.md)).
