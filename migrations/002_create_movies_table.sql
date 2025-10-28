-- +goose Up
CREATE TABLE movies
(
    id            SERIAL PRIMARY KEY,
    title         VARCHAR(255)        NOT NULL,
    description   VARCHAR(255)        NOT NULL,
    average_score REAL                NOT NULL,
    release_date  DATE                NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- +goose Down
DROP TABLE movies