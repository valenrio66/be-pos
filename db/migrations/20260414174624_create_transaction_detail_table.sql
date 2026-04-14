-- migrate:up
CREATE TABLE transaction_details (
    id BIGSERIAL PRIMARY KEY,
    transaction_id BIGINT NOT NULL REFERENCES transactions(id) ON DELETE CASCADE,
    product_id BIGINT NOT NULL REFERENCES products(id),
    qty INT NOT NULL,
    price_at_buy NUMERIC(15, 2) NOT NULL, -- Harga saat dibeli (antisipasi harga produk berubah di masa depan)
    subtotal NUMERIC(15, 2) NOT NULL
);

-- migrate:down
DROP TABLE transaction_details;