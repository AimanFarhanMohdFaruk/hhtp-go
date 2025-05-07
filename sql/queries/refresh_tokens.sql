-- name: GetRefreshToken :one
SELECT * FROM refresh_tokens
WHERE token = $1 LIMIT 1;

-- name: CreateRefreshToken :one
INSERT INTO refresh_tokens (
  token, user_id, expires_at
) VALUES (
  $1, $2, $3
)
RETURNING *;

-- name: UpdateRefreshToken :exec
UPDATE refresh_tokens
  set revoked_at = $2
WHERE token = $1;

-- name: DeleteRefreshToken :exec
DELETE FROM refresh_tokens
WHERE token = $1;