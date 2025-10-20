-- +goose Up
-- +goose StatementBegin
ALTER TABLE reviews ADD is_deleted BOOLEAN DEFAULT FALSE;
-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
ALTER TABLE reviews DROP COLUMN is_deleted
-- +goose StatementEnd
