-- name: CreateUrl :one
INSERT INTO urls (
  hash_id,
  url
) VALUES (
  $1,
  $2
) RETURNING *;

-- name: GetUrlByHashId :one
SELECT * FROM urls WHERE hash_id=$1;