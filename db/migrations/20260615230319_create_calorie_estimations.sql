-- +goose Up
CREATE TABLE calorie_estimations (
    id                 UUID PRIMARY KEY,
    user_id            UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    description        TEXT NOT NULL,
    estimated_calories INTEGER NOT NULL,
    confidence         TEXT NOT NULL,
    items              JSONB NOT NULL DEFAULT '[]',
    observation        TEXT,
    created_at         TIMESTAMP NOT NULL DEFAULT now()
);

CREATE INDEX idx_calorie_estimations_user_id ON calorie_estimations(user_id);

-- +goose Down
DROP TABLE calorie_estimations;
