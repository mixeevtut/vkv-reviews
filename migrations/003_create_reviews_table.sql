-- +goose Up
CREATE TABLE reviews
(
    id         SERIAL PRIMARY KEY,
    user_id    INT          NOT NULL,
    movie_id   INT          NOT NULL,
    score      INT          NOT NULL,
    content    VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    -- Внешние ключи
    CONSTRAINT fk_reviews_user
        FOREIGN KEY (user_id) REFERENCES users(id)
            ON DELETE CASCADE,

    CONSTRAINT fk_reviews_movie
        FOREIGN KEY (movie_id) REFERENCES movies(id)
            ON DELETE CASCADE,

    -- Один отзыв на пользователя и фильм
    CONSTRAINT unique_user_movie_review
        UNIQUE(user_id, movie_id)
);

CREATE INDEX idx_reviews_user_id ON reviews (user_id);
CREATE INDEX idx_reviews_movie_id ON reviews (movie_id);

-- +goose Down
DROP TABLE reviews