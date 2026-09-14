# Database Schema (v1)

Datastore: PostgreSQL (see [ADR-0008](../adrs/0008-database-postgresql.md)). Migrations: `golang-migrate`, files in `backend/migrations/` (see [ADR-0009](../adrs/0009-schema-migrations.md)).

This schema covers the Phase 1 scope (barcode scan → product → reviews → rating summary, anonymous + email auth) plus columns/tables clearly anticipated by the long-term roadmap (ratings, photos, likes) so the initial migrations don't need immediate breaking changes. Fields explicitly out of scope for v1 (categories, nutrition, comments, price history) are called out but not created yet — they get their own migrations when built, per [ADR-0001](../adrs/0001-record-architecture-decisions.md).

## Entity overview

```mermaid
erDiagram
    USERS ||--o{ IDENTITIES : has
    USERS ||--o{ REVIEWS : writes
    USERS ||--o{ REVIEW_REPORTS : files
    PRODUCTS ||--o{ REVIEWS : receives
    PRODUCTS ||--|| RATING_SUMMARIES : has
    REVIEWS ||--o{ REVIEW_MEDIA : has
    REVIEWS ||--o{ REVIEW_LIKES : receives
    REVIEWS ||--o{ REVIEW_REPORTS : receives

    USERS {
        uuid id PK
        text display_name
        text status
        timestamptz created_at
    }
    IDENTITIES {
        uuid id PK
        uuid user_id FK
        text provider
        text provider_subject
        text password_hash
        timestamptz created_at
    }
    PRODUCTS {
        uuid id PK
        text barcode UK
        text name
        text brand
        text image_object_key
        text source
        jsonb attributes
        timestamptz created_at
        timestamptz updated_at
    }
    REVIEWS {
        uuid id PK
        uuid product_id FK
        uuid user_id FK
        text body
        smallint rating
        text status
        timestamptz created_at
        timestamptz updated_at
        timestamptz deleted_at
    }
    REVIEW_MEDIA {
        uuid id PK
        uuid review_id FK
        text object_key
        text media_type
        smallint position
    }
    RATING_SUMMARIES {
        uuid product_id PK_FK
        integer review_count
        numeric average_rating
        timestamptz updated_at
    }
    REVIEW_LIKES {
        uuid review_id FK
        uuid user_id FK
        timestamptz created_at
    }
    REVIEW_REPORTS {
        uuid id PK
        uuid review_id FK
        uuid reporter_user_id FK
        text reason
        text status
        timestamptz created_at
    }
```

## Table notes

### `users`
Every client gets a row here immediately (anonymous account creation, see [ADR-0014](../adrs/0014-authentication-strategy.md)). `status` is `active | suspended | deleted` (soft states for moderation). No email/password columns here — credentials live in `identities`.

### `identities`
One row per linked auth method for a user. `provider` is `anonymous | email | google | apple | ...`. `provider_subject` is the provider's stable identifier (for `email`, the email address itself, unique per provider). `password_hash` (Argon2id) is only populated for `provider = 'email'`. Unique constraint on `(provider, provider_subject)`. This table is what makes adding an OAuth provider a data change, not a schema change (per ADR-0014).

### `products`
`barcode` is unique and indexed — the single hottest lookup path in the system (`GET /v1/products/{barcode}`). All descriptive fields except `barcode` are nullable ("tolerate incomplete product information" per the product requirements). `source` is `placeholder | catalog_import | user_submitted | admin`, tracking provenance for future duplicate-detection/merging work. `attributes` is a `jsonb` column reserved for sparse future fields (nutrition, specs) per [ADR-0008](../adrs/0008-database-postgresql.md) — empty `{}` in v1, not yet given first-class columns.

### `reviews`
`rating` is nullable smallint (1-5) — present in the schema from v1 even though the initial product overview's "typical flow" describes text-only reviews, because star ratings are explicitly listed as the first near-term addition and a nullable column now avoids a breaking migration later; the API may choose not to expose/require it until that feature ships. `status` is `published | pending_moderation | removed`. `deleted_at` supports soft delete (edit history / moderation audit trail, both on the roadmap) rather than hard-deleting user content.

### `review_media`, `review_likes`, `review_reports`
Structurally present from v1 (cheap, additive, and referenced directly by the roadmap: photos, likes, abuse reporting) but not wired into the API until their respective features ship — see the relevant future ADRs before building against them.

### `rating_summaries`
Denormalized, one row per product, recomputed transactionally (or via a lightweight trigger/background job — implementation detail, not an ADR-level decision yet) whenever a review is created/updated/deleted, so `GET /v1/products/{barcode}` never needs to aggregate `reviews` on the read path.

## Indexes (initial)
- `products(barcode)` — unique.
- `reviews(product_id, created_at desc)` — paginated review lists per product.
- `reviews(user_id)` — a user's own reviews.
- `identities(provider, provider_subject)` — unique, login lookup.
- GIN index on `products` full-text/trigram columns once search ([ADR-0013](../adrs/0013-search-strategy.md)) is implemented.

## Not yet modeled (tracked here so future migrations have a landing spot)
Categories, manufacturer/region, product history, price tracking, comments-on-reviews (would need a `parent_review_id` or separate `review_comments` table), review edit history (would need a `review_revisions` table).
