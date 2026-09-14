# ADR-0010: Redis for Caching, Rate Limiting, and Ephemeral State

## Status
Accepted

## Context
Barcode lookups are read-heavy and highly repetitive (many users scan the same popular products). The system also needs rate limiting (see [[0015-rate-limiting]]) and, later, features like session/token state and real-time-ish counters (likes, view counts).

## Problem
What, if anything, sits between the API and PostgreSQL to reduce read load and support ephemeral/shared state across backend replicas?

## Alternatives
- **No cache, rely on Postgres + connection pooling** — simplest, but every barcode scan hits the primary database directly; popular products (common consumer goods) would create hot rows under load, and there's no shared place for cross-replica rate-limit counters once the backend runs more than one instance.
- **In-process (in-memory) cache per backend instance** — fast, zero extra infra, but not shared across replicas (inconsistent behavior, cold cache on every deploy/restart) and useless for cross-replica rate limiting.
- **Redis (managed via AWS ElastiCache)** — shared, sub-millisecond cache across all backend replicas, supports TTL-based caching (product lookups, rating summaries), atomic counters and sliding-window structures well suited to rate limiting, and pub/sub if needed later for real-time features (e.g., live review updates).

## Decision
Redis, deployed as AWS ElastiCache for Redis. Used for: (1) caching product-by-barcode lookups and rating summaries with short TTLs (invalidated on write), (2) rate-limiting counters (see [[0015-rate-limiting]]), (3) later, session/short-lived token state if needed. Redis is a cache/ephemeral store, not a system of record — nothing is ever written to Redis without also being durably persisted in PostgreSQL or S3; a total Redis flush must never cause data loss, only a temporary performance dip.

## Consequences
- Reduces read load on PostgreSQL for hot paths (barcode lookup is the single most frequent operation in the product).
- Adds an operational component (ElastiCache) and a cache-invalidation responsibility on every write path that touches cached data.
- Enables correct multi-replica rate limiting once the backend scales beyond one instance.

## Future Considerations
If real-time features (live-updating review feeds, presence) grow significantly, evaluate Redis pub/sub or a dedicated real-time layer (e.g., WebSocket gateway backed by Redis Streams) — the caching decision here doesn't block that.
