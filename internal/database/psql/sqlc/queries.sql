-- name: GetAllTransactions :many
SELECT * FROM transactions;

-- name: AddTransaction :one 
INSERT INTO transactions (date, description, amount_out, amount_in, currency, bank)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;