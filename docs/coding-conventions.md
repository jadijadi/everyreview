# Coding Conventions

## General
- Follow the accepted ADRs in `adrs/`. If a change contradicts an `Accepted` ADR, write a superseding ADR first (see [ADR-0001](../adrs/0001-record-architecture-decisions.md)) — do not silently diverge in code.
- No premature abstraction: prefer duplication over a shared abstraction used in fewer than three places.
- Comments explain *why*, not *what*; well-named identifiers should make the "what" self-evident.

## Backend (Go)
- Formatted with `gofmt`/`goimports`; linted with `golangci-lint` (CI-enforced, see [ADR-0021](../adrs/0021-cicd-github-actions.md)).
- Package layout follows module boundaries from [ADR-0003](../adrs/0003-modular-monolith-backend.md): `backend/internal/<domain>/{handler,service,repository}.go` style, one package per domain module. A module's repository is the *only* code that issues SQL against that module's tables.
- Errors: wrap with `fmt.Errorf("...: %w", err)` to preserve chains; sentinel errors (`var ErrNotFound = errors.New(...)`) for conditions callers branch on.
- All request/response types are explicit Go structs matching `backend/api/openapi.yaml` — the spec is the source of truth; if they diverge, the contract test (ADR-0022) fails CI.
- Every exported function that does I/O takes a `context.Context` as its first argument.

## Android (Kotlin)
- Follows the official [Kotlin style guide](https://developer.android.com/kotlin/style-guide); enforced via `ktlint` in CI.
- Layering per [ADR-0023](../adrs/0023-mobile-architecture.md): `presentation/`, `domain/`, `data/` packages per feature. Domain layer has zero Android SDK imports — this is what keeps it unit-testable without an emulator.
- ViewModels expose a single `StateFlow<UiState>` per screen, not scattered mutable properties — keeps state changes traceable and testable.
- No business logic in Composables; a Composable renders state and forwards user intent to the ViewModel.

## Commit / PR conventions
- Conventional-commit-style prefixes (`feat:`, `fix:`, `refactor:`, `docs:`, `chore:`) recommended for changelog generation, not strictly enforced.
- A PR that changes `backend/api/openapi.yaml` in a breaking way must reference the corresponding ADR-0007 versioning decision and bump the API version if required.
- A PR introducing a new third-party dependency with architectural weight (a new datastore, a new external service, a new mobile library affecting more than one screen) needs an accompanying ADR.
