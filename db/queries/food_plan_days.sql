-- name: CreateFoodPlanDay :one
INSERT INTO food_plan_days (id, food_plan_id, weekday)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetFoodPlanDayByWeekday :one
SELECT * FROM food_plan_days
WHERE food_plan_id = $1 AND weekday = $2;

-- name: ListFoodPlanDays :many
SELECT * FROM food_plan_days
WHERE food_plan_id = $1
ORDER BY weekday;
