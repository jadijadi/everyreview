# Architecture Decision Records

This directory records every significant architectural decision for EveryReview, per [ADR-0001](0001-record-architecture-decisions.md). New ADRs use [`template.md`](template.md), are numbered sequentially, and are never edited after acceptance — superseded decisions get a new ADR that supersedes the old one.

Implementation must follow `Accepted` ADRs. Proposals (including LLM-generated ones) land as `Proposed` ADRs before code is written against them.

| # | Title | Status |
|---|-------|--------|
| [0001](0001-record-architecture-decisions.md) | Record Architecture Decisions | Accepted |
| [0002](0002-repository-structure.md) | Monorepo Layout | Accepted |
| [0003](0003-modular-monolith-backend.md) | Modular Monolith Backend Architecture | Accepted |
| [0004](0004-backend-language-go.md) | Go as Backend Language | Accepted |
| [0005](0005-backend-framework.md) | `net/http` + `chi` Framework | Accepted |
| [0006](0006-api-style-rest.md) | REST over GraphQL | Accepted |
| [0007](0007-api-versioning.md) | URL Path API Versioning | Accepted |
| [0008](0008-database-postgresql.md) | PostgreSQL as Primary Datastore | Accepted |
| [0009](0009-schema-migrations.md) | `golang-migrate` + `sqlc` for Schema/Queries | Accepted |
| [0010](0010-caching-redis.md) | Redis for Caching/Rate Limiting | Accepted |
| [0011](0011-object-storage-and-image-handling.md) | S3 Object Storage & Image Uploads | Accepted |
| [0012](0012-cdn-cloudfront.md) | CloudFront CDN | Accepted |
| [0013](0013-search-strategy.md) | Postgres Full-Text Search (Initial) | Accepted |
| [0014](0014-authentication-strategy.md) | Anonymous-First Auth with JWT | Accepted |
| [0015](0015-rate-limiting.md) | Redis Token-Bucket Rate Limiting | Accepted |
| [0016](0016-logging-strategy.md) | Structured JSON Logging | Accepted |
| [0017](0017-monitoring-and-observability.md) | CloudWatch + OpenTelemetry Metrics | Accepted |
| [0018](0018-infrastructure-as-code-terraform.md) | Terraform on AWS | Accepted |
| [0019](0019-compute-and-deployment.md) | ECS Fargate Compute | Accepted |
| [0020](0020-secrets-management.md) | AWS Secrets Manager | Accepted |
| [0021](0021-cicd-github-actions.md) | GitHub Actions CI/CD | Accepted |
| [0022](0022-testing-strategy.md) | Layered Testing Strategy | Accepted |
| [0023](0023-mobile-architecture.md) | MVVM + Clean Architecture (Mobile) | Accepted |
| [0024](0024-android-tech-stack.md) | Android Tech Stack | Accepted |
| [0025](0025-barcode-scanning-library.md) | CameraX + ML Kit Barcode Scanning | Accepted |

## Adding a new ADR
1. Copy `template.md` to `NNNN-kebab-title.md` using the next sequential number.
2. Fill in every section — an ADR without alternatives or consequences is not useful.
3. Status starts as `Proposed`; a maintainer moves it to `Accepted` on merge.
4. Add a row to the table above.
5. Cross-link related ADRs using `[[NNNN-slug]]`-style references in prose (rendered as normal Markdown links here).
