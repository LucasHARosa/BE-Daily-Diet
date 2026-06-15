-- +goose Up
CREATE TABLE meals (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    food_plan_meal_id UUID,
    food_plan_meal_snapshot JSONB,
    name TEXT NOT NULL,
    description TEXT,
    eaten_at TIMESTAMP NOT NULL,
    is_on_diet BOOLEAN NOT NULL,
    calories INTEGER,
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    updated_at TIMESTAMP NOT NULL DEFAULT now()
);

CREATE INDEX idx_meals_user_id ON meals(user_id);
CREATE INDEX idx_meals_user_eaten_at ON meals(user_id, eaten_at);

-- +goose Down
DROP TABLE meals;
