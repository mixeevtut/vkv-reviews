-- +goose Up
CREATE TABLE movies
(
    id            SERIAL PRIMARY KEY,
    title         VARCHAR(255) UNIQUE NOT NULL,
    description   VARCHAR(255)        NOT NULL,
    average_score REAL                NOT NULL,
    release_date  DATE                NOT NULL
);

-- +goose Down
DROP TABLE movies