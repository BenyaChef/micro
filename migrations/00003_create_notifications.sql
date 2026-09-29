-- +goose Up
CREATE TABLE notifications (
    id         uuid        PRIMARY KEY,
    event_id   uuid        NOT NULL UNIQUE,
    user_id    uuid        NOT NULL,
    kind       text        NOT NULL,
    payload    jsonb       NOT NULL,
    sent_at    timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX notifications_user_id_idx ON notifications (user_id);

-- +goose Down
DROP TABLE notifications;
