-- +goose Up
-- +goose StatementBegin
ALTER TABLE url DROP CONSTRAINT url_short_url_key;
ALTER TABLE url ADD user_id uuid;
ALTER TABLE url ALTER COLUMN user_id SET NOT NULL;
ALTER TABLE url DROP CONSTRAINT url_pkey;
ALTER TABLE url ADD CONSTRAINT url_pkey PRIMARY KEY (user_id, short_url);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
CREATE INDEX IF NOT EXISTS idx_short_url ON url(short_url);
ALTER TABLE url DROP COLUMN user_id;
ALTER TABLE url DROP CONSTRAINT url_pkey;
ALTER TABLE url ADD PRIMARY KEY (id);
-- +goose StatementEnd
