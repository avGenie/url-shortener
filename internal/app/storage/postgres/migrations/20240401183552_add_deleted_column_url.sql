-- +goose Up
-- +goose StatementBegin
ALTER TABLE url ADD COLUMN deleted BOOLEAN NOT NULL DEFAULT false;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE url DROP COLUMN deleted;
-- +goose StatementEnd
