-- name: GetMetricsSummary :one
SELECT
    COUNT(*)                              AS total_meals,
    COUNT(*) FILTER (WHERE is_on_diet)    AS total_on_diet,
    COUNT(*) FILTER (WHERE NOT is_on_diet) AS total_off_diet
FROM meals
WHERE user_id = $1;

-- name: ListMealsOnDietStatus :many
SELECT is_on_diet FROM meals
WHERE user_id = $1
  AND (sqlc.narg('start_date')::TIMESTAMP IS NULL OR eaten_at >= sqlc.narg('start_date'))
  AND (sqlc.narg('end_date')::TIMESTAMP IS NULL OR eaten_at <= sqlc.narg('end_date'))
ORDER BY eaten_at ASC;

-- name: ListMealsGroupedByDay :many
SELECT
    DATE(eaten_at)                                  AS date,
    COUNT(*)                                        AS total_meals,
    COUNT(*) FILTER (WHERE is_on_diet)              AS total_on_diet,
    COUNT(*) FILTER (WHERE NOT is_on_diet)          AS total_off_diet,
    COALESCE(SUM(calories), 0)::BIGINT              AS total_calories
FROM meals
WHERE user_id = $1
  AND eaten_at >= $2
  AND eaten_at <= $3
GROUP BY DATE(eaten_at)
ORDER BY date;

-- name: ListMealsGroupedByMonth :many
SELECT
    DATE_TRUNC('month', eaten_at)::DATE             AS date,
    COUNT(*)                                        AS total_meals,
    COUNT(*) FILTER (WHERE is_on_diet)              AS total_on_diet,
    COUNT(*) FILTER (WHERE NOT is_on_diet)          AS total_off_diet,
    COALESCE(SUM(calories), 0)::BIGINT              AS total_calories
FROM meals
WHERE user_id = $1
  AND eaten_at >= $2
  AND eaten_at <= $3
GROUP BY DATE_TRUNC('month', eaten_at)
ORDER BY date;
