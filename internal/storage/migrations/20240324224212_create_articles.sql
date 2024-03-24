-- +goose Up
-- +goose StatementBegin
CREATE TABLE articles (
    article_id int GENERATED ALWAYS AS IDENTITY,
    source_id int NOT NULL,
    title varchar(255) NOT NULL,
    link varchar(255) NOT NULL,
    summary text NOT NULL,
    published_at timestamp NOT NULL,
    created_at timestamp NOT NULL,
    posted_at timestamp,

    PRIMARY KEY(article_id),
    CONSTRAINT fk_articles_source_id
      FOREIGN KEY(source_id) 
        REFERENCES sources(source_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS articles;
-- +goose StatementEnd
