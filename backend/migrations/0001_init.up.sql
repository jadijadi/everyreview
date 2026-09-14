CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- MVP simplification: no users/identities tables yet (see backend/README.md
-- "MVP scope" note) — reviews store a plain author_name instead of a user_id FK.
CREATE TABLE products (
    id               uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    barcode          text NOT NULL,
    name             text,
    brand            text,
    image_object_key text,
    source           text NOT NULL DEFAULT 'placeholder',
    attributes       jsonb NOT NULL DEFAULT '{}',
    created_at       timestamptz NOT NULL DEFAULT now(),
    updated_at       timestamptz NOT NULL DEFAULT now()
);

CREATE UNIQUE INDEX products_barcode_idx ON products (barcode);

CREATE TABLE reviews (
    id          uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    product_id  uuid NOT NULL REFERENCES products (id),
    author_name text NOT NULL DEFAULT 'Anonymous',
    body        text NOT NULL,
    rating      smallint CHECK (rating BETWEEN 1 AND 5),
    status      text NOT NULL DEFAULT 'published',
    created_at  timestamptz NOT NULL DEFAULT now(),
    updated_at  timestamptz NOT NULL DEFAULT now(),
    deleted_at  timestamptz
);

CREATE INDEX reviews_product_id_created_at_idx ON reviews (product_id, created_at DESC);
