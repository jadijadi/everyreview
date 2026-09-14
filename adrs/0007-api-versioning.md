# ADR-0007: API Versioning via URL Path Prefix

## Status
Accepted

## Context
Two mobile clients (Android now, iOS later) will be on the field on different release cadences than the backend, and users on old app versions cannot be force-upgraded instantly. The API contract will change over time as features are added (see [[0006-api-style-rest]]).

## Problem
How do we evolve the API without breaking already-installed client versions?

## Alternatives
- **No versioning, additive-only changes forever** — simplest, but eventually some change (removing a field, changing semantics) becomes unavoidable and breaks old clients silently.
- **Header-based versioning** (`Accept: application/vnd.everyreview.v2+json`) — flexible, but invisible in logs/URLs, harder to route/cache differently per version at the CDN/load-balancer layer, and less discoverable.
- **URL path versioning** (`/v1/products/...`, `/v2/products/...`) — explicit, visible in logs and CDN cache keys, trivial to route at the load balancer or reverse proxy to different backend versions if ever needed, simple for client SDKs to target.

## Decision
Version the API in the URL path: `/v1/...`. Within a version, changes must be backward-compatible (additive fields, new optional query params, new endpoints). Breaking changes require a new version prefix (`/v2/...`), and the backend supports at most two major versions concurrently, with a documented deprecation window (minimum 6 months) communicated via the `Deprecation` and `Sunset` HTTP headers on the old version's responses.

## Consequences
- Old app installs keep working against `/v1` while new features ship under `/v1` (additive) or, when unavoidable, `/v2`.
- Running two versions concurrently means some duplicated handler code during a transition window; kept minimal by only bumping the major version for true breaking changes, not every feature.
- Simple to reason about and to test.

## Future Considerations
If breaking changes become frequent, consider narrower per-resource versioning instead of a single global version — but only if the global-version approach demonstrably causes unnecessary churn.
