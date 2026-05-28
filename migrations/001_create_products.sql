-- Jalankan file ini sekali untuk menyiapkan tabel di database
-- Perintah: psql -U postgres -d go_api_db -f migrations/001_create_products.sql

CREATE TABLE IF NOT EXISTS products (
    id    SERIAL PRIMARY KEY,
    name  VARCHAR(255) NOT NULL,
    price INTEGER      NOT NULL DEFAULT 0,
    stock INTEGER      NOT NULL DEFAULT 0
);

-- Opsional: Masukkan data awal (seed data) untuk testing
INSERT INTO products (name, price, stock) VALUES
    ('Kopi Hitam', 15000, 100),
    ('Susu Murni', 20000, 50);
