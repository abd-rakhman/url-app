-- name: CreateUrl :one
INSERT INTO urls (
  hash_id,
  url
) VALUES (
  $1,
  $2
) RETURNING *;

-- name: CreateUrlWithExpiresAt :one
INSERT INTO urls (
  hash_id,
  url,
  expires_at
) VALUES (
  $1,
  $2,
  $3
) RETURNING *;

-- name: GetUrlByHashId :one
SELECT * FROM urls WHERE hash_id=$1 and now() < expires_at;

-- name: GetUrlByHashIdForUpdate :one
SELECT * FROM urls
WHERE hash_id=$1
FOR UPDATE;

-- name: DeleteExpiredUrls :exec
DELETE FROM urls WHERE now() > expires_at;