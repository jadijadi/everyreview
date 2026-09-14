# ADR-0009: SQL-First Migrations with `golang-migrate`

## Status
Accepted

## Context
The database schema will evolve continuously (ADR-0001's long-term feature list implies dozens of future schema changes: categories, photos, likes, reports, etc.). The team uses Go and PostgreSQL (see [[0004-backend-language-go]], [[0008-database-postgresql]]).

## Problem
How are database schema changes authored, reviewed, and applied across environments?

## Alternatives
- **ORM-managed migrations (e.g., GORM auto-migrate)** — convenient but often "magic," produces hard-to-review generated DDL, and encourages an ORM-first modeling style that fights explicit SQL performance tuning.
- **Application-managed ad hoc scripts** — no structure, high risk of environment drift.
- **Dedicated SQL migration tool (`golang-migrate`)** — plain, numbered, up/down `.sql` files, applied by a CLI/library, fully explicit and reviewable in PRs as raw SQL, no dependency on an ORM's code-generation model. Works identically in CI, local dev, and production via the same migration files.

## Decision
Use `golang-migrate` with versioned, paired up/down `.sql` files in `backend/migrations/`. Migrations run automatically as a pre-deploy CI/CD step (see [[0021-cicd]]) before the new backend version is rolled out, never from application startup code (to avoid multiple replicas racing to migrate concurrently). Data access in application code uses plain SQL via `database/sql` + `sqlc` for compile-time-checked query code generation from those same `.sql` queries — not a full ORM — keeping query performance visible and explicit.

## Consequences
- Every schema change is an explicit, reviewable SQL diff.
- No hidden ORM-generated queries; query performance is easy to reason about and `sqlc` catches type mismatches at build time.
- Slightly more manual than an ORM for simple CRUD (writing SQL by hand), an accepted trade-off for clarity and performance control given the growth-to-millions-of-rows requirement.

## Future Considerations
If schema changes need zero-downtime patterns for large tables (e.g., adding a NOT NULL column to a 100M-row `products` table), adopt expand/contract migration patterns (add nullable → backfill → add constraint in separate migrations) — document this as a convention in `docs/coding-conventions.md` once it's needed.
