-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS url(
    id SERIAL PRIMARY KEY,
    short_url TEXT NOT NULL UNIQUE,
    url TEXT NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_short_url ON url(short_url);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS url;
-- +goose StatementEnd
