-- migrate:up
CREATE TABLE products (
    id BIGSERIAL PRIMARY KEY,
    sku VARCHAR(50) UNIQUE NOT NULL, -- Stock Keeping Unit (Kode Barang, misal: BT-SPLIT-01)
    name VARCHAR(255) NOT NULL,      -- Nama Barang (misal: Batu Split 1/2)
    description TEXT,                -- Deskripsi Opsional
    price NUMERIC(15, 2) NOT NULL,   -- Harga (pakai NUMERIC agar presisi untuk uang)
    unit VARCHAR(50) NOT NULL,       -- Satuan (misal: m3, rit, engkel, sak)
    stock INT NOT NULL DEFAULT 0,    -- Sisa Stok
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Indexing untuk mempercepat pencarian nama barang oleh kasir nanti
CREATE INDEX idx_products_name ON products(name);

-- migrate:down
DROP INDEX idx_products_name;
DROP TABLE products;