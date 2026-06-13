-- name: CreateShortUrl :one
INSERT INTO short_urls (id, created_at, code, original_url, clicks)
VALUES (
  gen_random_id(),
  NOW(),
  $1,
  $2,
  0
)
RETURNING *;

-- name: GetByCode :one
SELECT *
FROM short_urls
WHERE code = $1;

-- name: ListAllUnfiltered :many
SELECT *
FROM short_urls
ORDER BY created_at ASC;

-- name: ListAllFiltered :many
SELECT *
FROM short_urls
WHERE original_url LIKE $1
ORDER BY created_at ASC;

-- name: ListSomeUnfiltered :many
SELECT *
FROM short_urls
ORDER BY created_at ASC
OFFSET $1
LIMIT $2;

-- name: ListSomeFiltered :many
SELECT *
FROM short_urls
WHERE original_url LIKE $3
ORDER BY created_at ASC
OFFSET $1
LIMIT $2;

-- name: IncrementClicks :exec
UPDATE short_urls
SET clicks = $2
WHERE code = $1;

-- name: DeleteByCode :exec
DELETE FROM short_urls
WHERE code = $1;