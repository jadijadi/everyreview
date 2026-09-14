# EveryReview

Scan a product's barcode, see what other people think of it. Submit your own review, and other users see it immediately.

This is an **ADR-driven project**: every significant architectural decision is documented in [`adrs/`](adrs/README.md) before it's implemented. Read that first if you're contributing — see [`docs/contributing.md`](docs/contributing.md).

## Status
Architecture and documentation phase complete (backend + Android, per the [development process](#development-process) below). No application code yet — see `adrs/` for what's been decided and why.

## Repository layout
```
backend/           Go backend service (modular monolith) — REST API
android/            Android app (Kotlin, Jetpack Compose)
ios/                iOS app (Swift, SwiftUI) — Phase 2
infrastructure/     Terraform IaC (AWS)
docs/               System docs, diagrams, schema, dev setup, conventions
adrs/               Architecture Decision Records — read this first
scripts/            Cross-cutting dev/ops scripts
.github/            CI/CD workflows
```
Each top-level directory has its own `README.md` with component-specific instructions.

## Start here
- [System overview](docs/system-overview.md) — architecture, diagrams, primary flow
- [ADR index](adrs/README.md) — every architectural decision and why
- [Database schema](docs/database-schema.md)
- [API spec](backend/api/openapi.yaml)
- [Development setup](docs/development-setup.md)
- [Deployment](docs/deployment.md)

## Development process
1. ADRs (this phase — done for Phase 1 scope)
2. Database schema — done, see `docs/database-schema.md`
3. API design — done, see `backend/api/openapi.yaml`
4. Authentication design — done, see [ADR-0014](adrs/0014-authentication-strategy.md)
5. Infrastructure design — done, see [ADR-0018](adrs/0018-infrastructure-as-code-terraform.md), [ADR-0019](adrs/0019-compute-and-deployment.md)
6. Implement backend
7. Backend tests
8. Implement Android app
9. Android tests
10. Deploy dev environment
11. Implement iOS app (Phase 2)
12. Advanced features (see the long-term roadmap below)

## Long-term roadmap
Star ratings, review photos/video, product categories, product search, user profiles, reputation, likes, comments, product history, price tracking, OCR of labels, AI summaries, AI moderation, AI-generated descriptions, recommendations, offline caching, push notifications, admin dashboard, analytics, i18n, spam detection, abuse reporting, product merging/duplicate detection. The architecture in `adrs/` is deliberately shaped so these land as additive changes — see each relevant ADR's "Future Considerations" section for the intended extension point.

## License
TBD.
