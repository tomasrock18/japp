CREATE TABLE IF NOT EXISTS products (
    barcode          VARCHAR(32) PRIMARY KEY,
    name             VARCHAR(255) NOT NULL,
    kcal_per_100g    DOUBLE PRECISION NOT NULL CHECK (kcal_per_100g >= 0),
    protein_per_100g DOUBLE PRECISION NOT NULL CHECK (protein_per_100g >= 0),
    fat_per_100g     DOUBLE PRECISION NOT NULL CHECK (fat_per_100g >= 0),
    carbs_per_100g   DOUBLE PRECISION NOT NULL CHECK (carbs_per_100g >= 0),
    created_by       BIGINT,
    created_at       TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
