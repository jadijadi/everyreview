# ADR-0015: Token-Bucket Rate Limiting at the API Layer, Backed by Redis

## Status
Accepted

## Context
The API is reachable by anonymous accounts (see [[0014-authentication-strategy]]) and mobile clients, both of which can be abused (scripted review spam, barcode-lookup scraping, placeholder-product creation flooding). The backend runs as multiple replicas (see [[0019-compute-and-deployment]]), so limiting must be consistent across instances.

## Problem
How do we prevent abusive traffic from degrading the service or polluting data, without over-throttling legitimate mobile users?

## Alternatives
- **No rate limiting, rely on infra-level protection only (AWS WAF/Shield)** — protects against volumetric DDoS but does nothing against application-level abuse (e.g., one authenticated-but-malicious user submitting thousands of spam reviews per minute).
- **In-process rate limiting per backend instance** — simple, but each replica tracks its own counters, so a client can get `N × replica_count` effective requests by hitting different instances; inconsistent.
- **Centralized token-bucket/sliding-window limiting backed by Redis** (see [[0010-caching-redis]]), keyed by user ID (authenticated) or IP+device-fingerprint (unauthenticated edge cases), enforced as `chi` middleware (see [[0005-backend-framework]]) at the API layer, with different limits per route class (cheap reads like barcode lookup: generous; expensive/abuse-prone writes like review submission: strict).

## Decision
Redis-backed token-bucket rate limiting as middleware, applied per-route-class with distinct limits (e.g., barcode lookups: ~60/min/user; review submission: ~5/min/user; anonymous account creation: strict per-IP limit to prevent farming). Limit configuration lives in code (not a separate service) but as a single declarative table, easy to tune without touching handler logic. Rate-limited responses return `429` with a `Retry-After` header.

## Consequences
- Consistent enforcement across all backend replicas.
- Adds a Redis round-trip to the request path (already present for caching, so marginal cost).
- Requires ongoing tuning of per-route limits as real usage patterns emerge; start conservative and loosen based on data rather than guessing generously upfront.

## Future Considerations
If abuse patterns get more sophisticated (distributed low-and-slow spam, content farms), pair rate limiting with the future AI moderation/spam-detection features from the long-term roadmap rather than trying to solve it purely with request-rate thresholds.
