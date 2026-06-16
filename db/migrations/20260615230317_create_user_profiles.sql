-- +goose Up
CREATE TABLE user_profiles (
    id                     UUID PRIMARY KEY,
    user_id                UUID NOT NULL UNIQUE REFERENCES users(id) ON DELETE CASCADE,
    weight_kg              NUMERIC(5,2),
    height_cm              INTEGER,
    birth_date             DATE,
    body_fat_percentage    NUMERIC(5,2),
    basal_calories         INTEGER,
    activity_level         TEXT,
    gym_frequency_per_week INTEGER,
    created_at             TIMESTAMP NOT NULL DEFAULT now(),
    updated_at             TIMESTAMP NOT NULL DEFAULT now()
);

-- +goose Down
DROP TABLE user_profiles;
