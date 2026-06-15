-- +goose Up
CREATE TABLE food_plan_days (
    id           UUID PRIMARY KEY,
    food_plan_id UUID NOT NULL REFERENCES food_plans(id) ON DELETE CASCADE,
    weekday      SMALLINT NOT NULL CHECK (weekday BETWEEN 0 AND 6),
    created_at   TIMESTAMP NOT NULL DEFAULT now(),
    updated_at   TIMESTAMP NOT NULL DEFAULT now(),
    UNIQUE (food_plan_id, weekday)
);

-- +goose Down
DROP TABLE food_plan_days;
