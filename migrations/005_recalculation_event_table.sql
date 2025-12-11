-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS recalculation_event (
                                                   id         BIGSERIAL PRIMARY KEY,
                                                   movie_id   BIGINT      NOT NULL,
                                                   created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
                                                   sent_at    TIMESTAMPTZ NULL
);

CREATE UNIQUE INDEX IF NOT EXISTS recalc_event_new_unique
    ON recalculation_event (movie_id)
    WHERE sent_at IS NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS recalculation_event;
-- +goose StatementEnd
