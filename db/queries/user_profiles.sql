-- name: GetUserProfile :one
SELECT * FROM user_profiles
WHERE user_id = $1;

-- name: UpsertUserProfile :one
INSERT INTO user_profiles (id, user_id, weight_kg, height_cm, birth_date, body_fat_percentage, basal_calories, activity_level, gym_frequency_per_week)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
ON CONFLICT (user_id) DO UPDATE
SET weight_kg              = EXCLUDED.weight_kg,
    height_cm              = EXCLUDED.height_cm,
    birth_date             = EXCLUDED.birth_date,
    body_fat_percentage    = EXCLUDED.body_fat_percentage,
    basal_calories         = EXCLUDED.basal_calories,
    activity_level         = EXCLUDED.activity_level,
    gym_frequency_per_week = EXCLUDED.gym_frequency_per_week,
    updated_at             = now()
RETURNING *;
