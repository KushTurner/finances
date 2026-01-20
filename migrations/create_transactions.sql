-- +goose Up
CREATE TABLE IF NOT EXISTS transactions (
    id SERIAL PRIMARY KEY,
    date DATE NOT NULL,
    description VARCHAR(500) NOT NULL,
    amount BIGINT NOT NULL,
    currency CHAR(3) NOT NULL DEFAULT 'GBP',
    bank VARCHAR(100) NOT NULL,
    category VARCHAR(100),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- +goose Down
DROP TABLE IF EXISTS transactions;
