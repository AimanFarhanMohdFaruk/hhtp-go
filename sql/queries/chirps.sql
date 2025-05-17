-- name: GetChirp :one
SELECT * FROM chirps
WHERE id = $1 LIMIT 1;

-- name: ListChirps :many
SELECT * FROM chirps
WHERE  (user_id=sqlc.narg('author_id') OR sqlc.narg('author_id') IS NULL)
ORDER BY created_at;

-- name: CreateChirp :one
INSERT INTO chirps (
  body, user_id
) VALUES (
  $1, $2
)
RETURNING *;

-- name: UpdateChirp :exec
UPDATE chirps
  set body = $2
WHERE id = $1;

-- name: DeleteChirp :exec
DELETE FROM chirps
WHERE id = $1;