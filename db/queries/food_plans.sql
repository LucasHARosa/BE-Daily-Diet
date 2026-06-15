-- name: CreateFoodPlan :one
INSERT INTO food_plans (id, user_id, title, description, is_active)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetFoodPlanByIDAndUserID :one
SELECT * FROM food_plans
WHERE id = $1 AND user_id = $2;

-- name: ListFoodPlansByUser :many
SELECT * FROM food_plans
WHERE user_id = $1
ORDER BY created_at DESC;

-- name: GetActiveFoodPlan :one
SELECT * FROM food_plans
WHERE user_id = $1 AND is_active = true
LIMIT 1;

-- name: UpdateFoodPlan :one
UPDATE food_plans
SET title = $3, description = $4, updated_at = now()
WHERE id = $1 AND user_id = $2
RETURNING *;

-- name: DeactivateAllUserFoodPlans :exec
UPDATE food_plans
SET is_active = false, updated_at = now()
WHERE user_id = $1;

-- name: SetFoodPlanActive :one
UPDATE food_plans
SET is_active = $2, updated_at = now()
WHERE id = $1
RETURNING *;

-- name: DeleteFoodPlan :exec
DELETE FROM food_plans
WHERE id = $1 AND user_id = $2;
