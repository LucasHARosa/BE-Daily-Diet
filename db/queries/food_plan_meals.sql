-- name: CreateFoodPlanMeal :one
INSERT INTO food_plan_meals (id, food_plan_day_id, name, description, scheduled_time, calories, sort_order)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;

-- name: GetFoodPlanMealByID :one
SELECT * FROM food_plan_meals
WHERE id = $1;

-- name: ListFoodPlanMealsByDay :many
SELECT * FROM food_plan_meals
WHERE food_plan_day_id = $1
ORDER BY sort_order, scheduled_time;

-- name: UpdateFoodPlanMeal :one
UPDATE food_plan_meals
SET name           = $2,
    description    = $3,
    scheduled_time = $4,
    calories       = $5,
    sort_order     = $6,
    updated_at     = now()
WHERE id = $1
RETURNING *;

-- name: DeleteFoodPlanMeal :exec
DELETE FROM food_plan_meals
WHERE id = $1;
