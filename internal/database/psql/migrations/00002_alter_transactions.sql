-- +goose Up
-- +goose StatementBegin
ALTER TABLE IF EXISTS transactions
ADD COLUMN bank VARCHAR(100);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE IF EXISTS transactions
DROP COLUMN IF EXISTS bank;
-- +goose StatementEnd
