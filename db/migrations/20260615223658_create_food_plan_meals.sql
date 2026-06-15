-- +goose Up
CREATE TABLE food_plan_meals (
    id               UUID PRIMARY KEY,
    food_plan_day_id UUID NOT NULL REFERENCES food_plan_days(id) ON DELETE CASCADE,
    name             TEXT NOT NULL,
    description      TEXT,
    scheduled_time   TEXT,
    calories         INTEGER,
    sort_order       INTEGER NOT NULL DEFAULT 0,
    created_at       TIMESTAMP NOT NULL DEFAULT now(),
    updated_at       TIMESTAMP NOT NULL DEFAULT now()
);

CREATE INDEX idx_food_plan_meals_day_id ON food_plan_meals(food_plan_day_id);

-- +goose Down
DROP TABLE food_plan_meals;
