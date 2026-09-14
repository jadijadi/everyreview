# ADR-0008: PostgreSQL as the Primary Datastore

## Status
Accepted

## Context
Core data is relational by nature: products, reviews, users, ratings, and their relationships (a review belongs to exactly one product and one user; ratings aggregate from reviews). The system must eventually support "every product with a barcode" — potentially tens of millions of product rows — plus full-text search on product names/brands and review text, and JSON-shaped, evolving product attributes (nutrition, specs) in later phases.

## Problem
Which primary database should the backend use?

## Alternatives
- **DynamoDB / other managed NoSQL** — scales effortlessly and has a lower ops burden on AWS, but its access-pattern-first modeling fights against the relational, ad-hoc-query needs here (rating aggregation, future admin/analytics queries, joins between products/reviews/users), and strong consistency + relational integrity (e.g., foreign keys preventing orphaned reviews) is harder to get right.
- **MongoDB** — flexible schema is appealing for the "tolerate incomplete product information" requirement, but this project needs strong relational integrity and transactional aggregation (rating counts/averages) more than flexible documents, and Postgres's `jsonb` columns already give schema flexibility where needed without giving up relational guarantees elsewhere.
- **MySQL** — a reasonable alternative, roughly equivalent for this use case; Postgres is chosen for its more capable `jsonb` support (future product attributes/specs), native full-text search (initial search strategy, see [[0013-search-strategy]]), and stronger extension ecosystem (e.g., `pg_trgm` for fuzzy matching, PostGIS if location-based features are ever added).
- **PostgreSQL** — mature, ACID-compliant, excellent relational modeling, `jsonb` for flexible/evolving attributes, native full-text search, well supported by every major cloud (including AWS RDS/Aurora), large talent pool.

## Decision
PostgreSQL, run as Amazon RDS for PostgreSQL initially (see [[0019-compute-and-deployment]]), with `jsonb` columns used for genuinely variable/sparse product attributes (nutrition facts, specs) rather than for core normalized entities (products, reviews, users, ratings), which stay in normalized tables. See `docs/database-schema.md` for the concrete schema.

## Consequences
- Strong relational integrity for the core domain (foreign keys, constraints) at the database layer, not just application layer.
- `jsonb` gives schema flexibility for sparse/evolving product attributes without a schema migration for every new attribute.
- Vertical scaling has limits; horizontal read scaling via read replicas is available on RDS when read traffic (product/review lookups) grows past a single instance's capacity.
- Full-text search via Postgres (`tsvector`) is adequate at moderate scale but not a substitute for a dedicated search engine at very large product catalogs (see [[0013-search-strategy]]).

## Future Considerations
If the product catalog reaches hundreds of millions of rows or write throughput on a single primary becomes a bottleneck, consider Amazon Aurora PostgreSQL (drop-in compatible, better read-replica scaling and storage auto-scaling) or partitioning/sharding hot tables (e.g., `reviews` by product or time). Migration path from RDS Postgres to Aurora Postgres is low-risk since the engine is wire-compatible.
