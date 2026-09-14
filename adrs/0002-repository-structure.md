# ADR-0002: Monorepo Layout for Backend, Mobile Clients, and Infrastructure

## Status
Accepted

## Context
The project has three client surfaces (backend, Android, later iOS) plus infrastructure-as-code and cross-cutting docs/ADRs. They evolve together in the early phases (API contract changes ripple into both clients immediately) but will eventually be built, tested, and released independently by potentially different teams.

## Problem
Should backend, Android, iOS, and infrastructure live in one repository or separate repositories?

## Alternatives
- **Polyrepo** — one repo per component. Clean CI isolation and independent access control, but coordinating an API change across 2-3 repos with separate PRs/versioning is slow, especially for a small team, and shared OpenAPI contracts drift.
- **Monorepo, single top-level** — everything in one repo with clear top-level boundaries. Atomic cross-cutting changes (e.g., API contract + Android client update in one PR), single source of truth for ADRs, shared CI config, simpler for a small team.
- **Monorepo with git submodules per client** — same as above but with submodule complexity for no real benefit at this scale.

## Decision
Single monorepo with top-level directories per concern:

```
backend/           Go backend service(s)
android/           Android app (Kotlin)
ios/                iOS app (Swift) — added in Phase 2
infrastructure/     Terraform IaC, environment configs
docs/               System docs, diagrams, API docs, conventions
adrs/               Architecture Decision Records
scripts/            Cross-cutting dev/ops scripts (setup, codegen, release)
.github/            CI/CD workflows, issue/PR templates
```

Each top-level project directory has its own `README.md` describing local setup, build, and test instructions specific to that component. The API contract (OpenAPI spec) lives in `backend/api/openapi.yaml` and is the single source of truth consumed by both mobile clients for codegen.

## Consequences
- One PR can atomically change backend API + OpenAPI spec + Android client, keeping them in sync.
- CI must scope jobs per changed path (e.g., only run Android CI when `android/**` changes) to keep build times reasonable.
- All contributors have read access to all components; if backend and mobile teams later need separate access control, this decision must be revisited.

## Future Considerations
If the org grows to the point where backend and mobile are owned by fully separate teams with different release cadences and access-control needs, split into polyrepo, using the OpenAPI spec (versioned and published, e.g., as a package or a dedicated `api-contract` repo) as the sync point.
