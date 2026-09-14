# ADR-0013: Postgres Full-Text Search Initially, with a Clear Upgrade Path

## Status
Accepted

## Context
Product search (by name/brand) is on the long-term feature list but not in the initial release (initial flow is barcode-first). Catalog size starts small and is expected to grow toward "every product with a barcode" over years, not overnight.

## Problem
What powers product/review text search, now and as the catalog scales?

## Alternatives
- **Dedicated search engine from day one (OpenSearch/Elasticsearch, Algolia, Meilisearch)** — best-in-class relevance, faceting, typo tolerance, but a second datastore to keep in sync with Postgres (via CDC or dual writes), extra operational cost and complexity, unjustified before there's a search feature at all.
- **Postgres full-text search (`tsvector`/`tsquery`, `pg_trgm` for fuzzy/typo-tolerant matching)** — no new infrastructure, transactionally consistent with the data it indexes (no sync lag), good enough relevance and performance for hundreds of thousands to low millions of rows with proper GIN indexes.
- **Client-side/naive `LIKE '%term%'` queries** — no indexing benefit, doesn't scale even at moderate size; rejected outright.

## Decision
Implement search using PostgreSQL full-text search (`tsvector` generated columns + GIN indexes) plus `pg_trgm` for fuzzy brand/name matching, scoped to the `product` module (see [[0003-modular-monolith-backend]]) behind a `SearchProducts` interface. Because that interface is the only way other code accesses search, swapping the implementation later doesn't touch calling code.

## Consequences
- Zero additional infrastructure or data-sync complexity for the initial search feature.
- Relevance ranking and faceted search are more limited than a dedicated engine, acceptable while catalog size and query sophistication are both modest.
- Search load shares the primary database's resources with transactional traffic; mitigated by read replicas if it becomes contention.

## Future Considerations
Migrate to OpenSearch (AWS-native, avoids a new vendor) when: catalog size makes GIN index performance/relevance inadequate, or product requirements need faceting/typo-tolerance/ranking beyond what Postgres offers. Because search is already isolated behind an interface, migration means implementing that interface against OpenSearch and standing up a CDC pipeline (e.g., via a Postgres logical replication slot or outbox pattern) to keep it in sync — not a rewrite of calling code.
