-- User-submitted details for products that were scanned before anyone described them.
-- price/currency are split so the app can format per locale; currency is ISO 4217.
ALTER TABLE products
    ADD COLUMN IF NOT EXISTS description  text,
    ADD COLUMN IF NOT EXISTS manufacturer text,
    ADD COLUMN IF NOT EXISTS price        numeric(12, 2) CHECK (price >= 0),
    ADD COLUMN IF NOT EXISTS currency     char(3);
