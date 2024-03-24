-- +goose Up
-- +goose StatementBegin
CREATE TABLE sources(
    source_id int GENERATED ALWAYS AS IDENTITY,
    name varchar(255) NOT NULL,
    feed_url varchar(255) NOT NULL,
    created_at timestamp NOT NULL DEFAULT NOW(),
    updated_at timestamp NOT NULL DEFAULT NOW(),

    PRIMARY KEY (source_id)
);

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS sources;

-- +goose StatementEnd
