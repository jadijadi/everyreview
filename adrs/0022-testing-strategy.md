# ADR-0022: Layered Testing Strategy — Unit, Integration, Contract, E2E

## Status
Accepted

## Context
The system spans a Go backend, Android (and later iOS) clients, and infrastructure, all evolving concurrently. Bugs are cheapest to catch at the lowest layer that can catch them; a single testing philosophy needs to apply consistently as the codebase (and number of contributors, human and LLM) grows.

## Problem
What kinds of tests exist, at what layer, and what's required for a PR to merge?

## Alternatives
- **Unit tests only** — fast and cheap, but modular-monolith module boundaries (see [[0003-modular-monolith-backend]]) and real Postgres/Redis behavior (constraints, query correctness) aren't verified by mocks alone.
- **Heavy end-to-end tests only** — high confidence but slow, flaky, and expensive to run on every PR; poor at pinpointing failure cause.
- **Layered pyramid: unit → integration → contract → a thin layer of E2E**, matched to what each layer is actually good at.

## Decision
- **Unit tests** (Go: standard `testing` + `testify`; Kotlin: JUnit + MockK) — pure business logic within a module, dependencies mocked via interfaces. Required for all new logic; run on every PR (see [[0021-cicd]]).
- **Integration tests** (Go) — real PostgreSQL and Redis via `testcontainers-go` in CI, testing each module's actual data access layer (query correctness, constraints) without mocking the database. Required for any code touching persistence.
- **Contract tests** — the OpenAPI spec (see [[0006-api-style-rest]]) is validated against actual handler responses in CI (e.g., via a schema-validation middleware in test mode), preventing the API from silently drifting from its documented contract that both mobile clients depend on.
- **Android**: unit tests for ViewModels/UseCases (see [[0023-mobile-architecture]]), instrumented UI tests (Compose testing APIs) for critical flows (scan → view reviews → submit review) run against an emulator in CI.
- **End-to-end tests** — a small, deliberately thin suite exercising the full stack (real backend in a test environment + API client) covering only the golden path (scan barcode → see product → submit review → see it appear) and the "product not found" path. Kept small because E2E suites are the most expensive to maintain; broader coverage comes from the layers below.
- **Merge requirement**: PRs must pass unit + integration + contract tests for any changed component; E2E runs on merge to `main` and nightly, not blocking every PR (to keep PR feedback fast).

## Consequences
- Most bugs are caught cheaply and specifically (unit/integration) rather than via slow, ambiguous E2E failures.
- Contract tests keep backend and mobile clients honest about the API shape they actually share, reducing integration surprises between the two codebases evolving somewhat independently.
- Requires test infrastructure (testcontainers, emulators) in CI, adding some pipeline complexity and runtime — mitigated by path-filtered, parallelized CI jobs (see [[0021-cicd]]).

## Future Considerations
As AI-driven features (moderation, summaries, recommendations) are added, extend this strategy with evaluation-style tests (golden-output comparisons, not strict equality) rather than trying to force deterministic unit tests onto inherently probabilistic components.
