ALTER TABLE products
    DROP COLUMN IF EXISTS currency,
    DROP COLUMN IF EXISTS price,
    DROP COLUMN IF EXISTS manufacturer,
    DROP COLUMN IF EXISTS description;
