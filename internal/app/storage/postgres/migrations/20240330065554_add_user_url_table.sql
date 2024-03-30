-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS users(
    id uuid PRIMARY KEY
);

CREATE TABLE IF NOT EXISTS users_url(
    users_id uuid NOT NULL,
    url_id INTEGER NOT NULL UNIQUE,
    CONSTRAINT fk_users FOREIGN KEY(users_id) REFERENCES users(id),
    CONSTRAINT fk_url FOREIGN KEY(url_id) REFERENCES url(id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS users_url;
-- +goose StatementEnd
