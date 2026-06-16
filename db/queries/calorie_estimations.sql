-- name: CreateCalorieEstimation :one
INSERT INTO calorie_estimations (id, user_id, description, estimated_calories, confidence, items, observation)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;

-- name: ListCalorieEstimationsByUser :many
SELECT * FROM calorie_estimations
WHERE user_id = $1
ORDER BY created_at DESC
LIMIT 20;
