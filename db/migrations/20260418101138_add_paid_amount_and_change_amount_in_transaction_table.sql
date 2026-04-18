-- migrate:up
ALTER TABLE transactions
    ADD COLUMN paid_amount NUMERIC(15, 2) NOT NULL DEFAULT 0,
    ADD COLUMN change_amount NUMERIC(15, 2) NOT NULL DEFAULT 0;

-- migrate:down
ALTER TABLE transactions
DROP COLUMN paid_amount,
DROP COLUMN change_amount;