# ADR-0006: REST over GraphQL for the Public API

## Status
Accepted

## Context
The API needs to serve two mobile clients (Android now, iOS later) with a fairly small, well-defined set of resources initially (products, reviews, users, auth), growing over time (search, likes, comments, reports).

## Problem
Should the public API be REST or GraphQL?

## Alternatives
- **GraphQL** — flexible client-driven queries, avoids over/under-fetching, single endpoint. Justified when clients have very different, evolving data needs or there are many nested/related resources queried together. Costs: added server complexity (resolvers, N+1 query management, custom caching since HTTP caching doesn't apply), a steeper learning curve, and harder-to-reason-about rate limiting/abuse prevention (arbitrary query cost). For this project's initial resource set (product by barcode, reviews list, submit review, auth), the fetch patterns are simple and don't vary much between clients.
- **REST** — resource-oriented, predictable URLs (`/products/{barcode}`, `/products/{id}/reviews`), leverages standard HTTP caching (ETags, CDN caching of GET requests — directly useful for product/review reads), simple to rate-limit per-route, simple to document with OpenAPI, and both Android and iOS have mature, low-boilerplate HTTP client tooling (Retrofit, URLSession) without needing a GraphQL client library.

## Decision
REST, specified with OpenAPI 3.x as the source of truth (`backend/api/openapi.yaml`), used for both human documentation and client codegen (Retrofit/Moshi models for Android, Swift codegen for iOS).

## Consequences
- Simple to cache GET endpoints (product lookups, review lists) at the CDN/HTTP layer — directly benefits the highest-traffic reads.
- Over-fetching is possible on complex future screens (e.g., a product detail screen needing product + reviews + rating summary in one call); addressed by designing composite "view" endpoints (e.g., `GET /products/{barcode}/summary`) rather than switching styles.
- No built-in schema-driven client codegen negotiation like GraphQL; mitigated by OpenAPI-based codegen in CI.

## Future Considerations
If a future feature (e.g., an admin dashboard or a highly dynamic recommendation feed) genuinely needs flexible, client-shaped queries across many relations, consider adding a scoped GraphQL endpoint alongside REST for that feature only, rather than migrating the whole API.
