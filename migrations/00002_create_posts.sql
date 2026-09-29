-- +goose Up
CREATE TABLE posts (
    id         uuid        PRIMARY KEY,
    author_id  uuid        NOT NULL,
    title      text        NOT NULL,
    body       text        NOT NULL,
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX posts_created_at_idx ON posts (created_at DESC);
CREATE INDEX posts_author_id_idx ON posts (author_id);

-- +goose Down
DROP TABLE posts;
