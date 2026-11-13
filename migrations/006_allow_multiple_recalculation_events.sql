-- +goose Up
DROP INDEX recalc_event_new_unique;
CREATE INDEX recalc_event_unsent_idx ON recalculation_event (id) WHERE sent_at IS NULL;

-- +goose Down
DROP INDEX recalc_event_unsent_idx;
DELETE FROM recalculation_event AS duplicate
USING recalculation_event AS first
WHERE duplicate.movie_id = first.movie_id
  AND duplicate.id > first.id
  AND duplicate.sent_at IS NULL
  AND first.sent_at IS NULL;
CREATE UNIQUE INDEX recalc_event_new_unique ON recalculation_event (movie_id) WHERE sent_at IS NULL;
