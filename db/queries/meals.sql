-- name: CreateMeal :one
INSERT INTO meals (
    id,
    user_id,
    name,
    description,
    eaten_at,
    is_on_diet,
    calories
) VALUES (
    $1, $2, $3, $4, $5, $6, $7
)
RETURNING *;

-- name: GetMealByIDAndUserID :one
SELECT * FROM meals
WHERE id = $1 AND user_id = $2;

-- name: ListMealsByUserFiltered :many
SELECT * FROM meals
WHERE user_id = $1
  AND (sqlc.narg('start_date')::TIMESTAMP IS NULL OR eaten_at >= sqlc.narg('start_date'))
  AND (sqlc.narg('end_date')::TIMESTAMP IS NULL OR eaten_at <= sqlc.narg('end_date'))
  AND (sqlc.narg('is_on_diet')::BOOLEAN IS NULL OR is_on_diet = sqlc.narg('is_on_diet'))
ORDER BY eaten_at DESC;

-- name: UpdateMeal :one
UPDATE meals
SET
    name        = $3,
    description = $4,
    eaten_at    = $5,
    is_on_diet  = $6,
    calories    = $7,
    updated_at  = now()
WHERE id = $1 AND user_id = $2
RETURNING *;

-- name: DeleteMeal :exec
DELETE FROM meals
WHERE id = $1 AND user_id = $2;
