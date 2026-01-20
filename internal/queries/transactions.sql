-- name: CreateTransaction :one
INSERT INTO transactions (
    date, description, amount, currency, bank, category
) VALUES (
    $1, $2, $3, $4, $5, $6
) RETURNING *;

-- name: GetTransaction :one
SELECT * FROM transactions
WHERE id = $1;

-- name: ListTransactions :many
SELECT * FROM transactions
ORDER BY date DESC;

-- name: UpdateTransaction :one
UPDATE transactions
SET date = $2,
    description = $3,
    amount = $4,
    currency = $5,
    bank = $6,
    category = $7,
    updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: DeleteTransaction :exec
DELETE FROM transactions
WHERE id = $1;
