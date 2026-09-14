# ADR-0003: Modular Monolith as Backend Service Architecture

## Status
Accepted

## Context
The long-term feature list is large (moderation, search, AI summaries, recommendations, price tracking, OCR, notifications, admin dashboard, analytics). A microservices architecture is often proposed to isolate these, but the initial product is a single small team shipping a barcode-scan-to-reviews MVP with unknown traffic and unknown which features will actually be used.

## Problem
Should the backend be built as microservices from day one, or as a single deployable service?

## Alternatives
- **Microservices from day one** — each capability (product, review, auth, moderation, search) as a separate deployable service. Enables independent scaling and deployment, but multiplies operational overhead (service discovery, distributed tracing, network failure modes, data consistency across services) before there is any traffic or team size to justify it. High risk of over-engineering an unproven product.
- **Single monolith, no internal structure** — fastest to start, but risks becoming an unmaintainable ball of mud as features accrete, making a later split into services expensive.
- **Modular monolith** — one deployable Go service, internally organized into packages with explicit boundaries per domain (`product`, `review`, `auth`, `user`, `moderation`, `media`), each owning its own data access and exposing a narrow internal API to other modules. No cross-module SQL queries; modules communicate through Go interfaces (in-process) rather than HTTP.

## Decision
Build a modular monolith. One deployable backend binary/container, internally split into domain modules under `backend/internal/<domain>/`, each with its own package boundary, tests, and (where reasonable) its own tables. Cross-module interaction happens through explicit interfaces defined at module boundaries, not shared database access — this is what makes a later extraction into a separate service (e.g., pulling out `moderation` or `search` when it needs independent scaling) a refactor rather than a rewrite.

## Consequences
- Single deployment unit: simpler CI/CD, simpler local dev (one process to run), no distributed transactions, no network calls between modules.
- Enforced module boundaries require discipline (and lint/CI checks) to avoid modules reaching into each other's tables directly.
- Scaling is initially vertical/horizontal-by-replica for the whole service, not per-module — acceptable until a specific module (e.g., image processing, AI summarization) has meaningfully different resource or scaling needs.

## Future Considerations
Extract a module into its own service when it has a distinct scaling profile, a distinct release cadence, or needs a different runtime (e.g., an AI/ML inference service). Because module boundaries are already enforced as interfaces, extraction means standing up a new service behind that same interface (now over the network, e.g., gRPC or REST) rather than restructuring the whole codebase. See also [[0025-android-tech-stack]] for how clients remain unaffected by this internal change (they only ever see the public REST API).
