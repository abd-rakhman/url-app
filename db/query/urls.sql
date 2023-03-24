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
SELECT * FROM urls WHERE hash_id=$1;

-- name: GetUrlByHashIdForUpdate :one
SELECT * FROM urls
WHERE hash_id=$1
FOR UPDATE;