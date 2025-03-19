/* Groups Table */

-- name: CreateGroup :one
INSERT INTO groups (name)
VALUES ($1)
RETURNING *;

-- name: GetGroup :one
SELECT id, name, created_at, updated_at, deleted_at FROM groups
WHERE id = $1 LIMIT 1;

-- name: GetGroupsWithPagination :many
SELECT id, name, created_at, updated_at FROM groups
WHERE deleted_at IS NULL
ORDER BY created_at DESC LIMIT $1 OFFSET $2;

-- name: GetGroupsCount :one
SELECT count(*) FROM groups
WHERE deleted_at IS NULL;

-- name: UpdateGroup :one
UPDATE groups
SET name = $2
WHERE id = $1
RETURNING *;;

-- name: DeleteGroup :exec
UPDATE groups
SET deleted_at = NOW()
WHERE id = $1 AND deleted_at IS NULL;

/* Songs Table */

-- name: CreateSong :one
INSERT INTO songs (group_id, title, runtime, lyrics, release_date, link)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;;

-- name: GetSong :one
SELECT id, group_id, title, runtime,  lyrics, release_date, link, created_at, updated_at, deleted_at
FROM songs
WHERE id = $1 LIMIT 1;

-- name: GetSongsWithPagination :many
SELECT id, group_id, title, runtime, lyrics, release_date, link, created_at, updated_at FROM songs
WHERE deleted_at IS NULL
ORDER BY created_at DESC LIMIT $1 OFFSET $2;

-- name: GetSongsCount :one
SELECT count(*) FROM songs
WHERE deleted_at IS NULL;

-- name: UpdateSong :one
UPDATE songs
SET
    group_id = $2,
    title = $3,
    runtime = $4,
    lyrics = $5,
    release_date = $6,
    link = $7
WHERE id = $1
RETURNING *;;

-- name: DeleteSong :execresult
UPDATE songs
SET deleted_at = NOW()
WHERE id = $1 AND deleted_at IS NULL;

-- name: GetSongsByGroup :many
SELECT id, group_id, title, runtime, lyrics, release_date, link, created_at, updated_at, deleted_at
FROM songs
WHERE group_id = $1 AND deleted_at IS NULL
ORDER BY release_date DESC LIMIT $2 OFFSET $3;

-- name: GetSongsWithFilters :many
SELECT s.id, s.group_id, s.title, s.runtime, s.lyrics, s.release_date, s.link, s.created_at,  s.updated_at
FROM songs s
         JOIN groups g ON s.group_id = g.id
WHERE s.deleted_at IS NULL
  AND (LOWER(g.name) LIKE LOWER('%' || NULLIF($3, '')::VARCHAR || '%') OR $3 = '')
  AND (LOWER(s.title) LIKE LOWER('%' || NULLIF($4, '')::VARCHAR || '%') OR $4 = '')
ORDER BY s.created_at DESC
    LIMIT $1 OFFSET $2;

-- name: GetSongsCountWithFilters :one
SELECT count(*)
FROM songs s
JOIN groups g ON s.group_id = g.id
WHERE s.deleted_at IS NULL
  AND (LOWER(g.name) LIKE LOWER('%' || NULLIF(@group_name, '')::VARCHAR || '%') OR @group_name = '')
  AND (LOWER(s.title) LIKE LOWER('%' || NULLIF(@song_title, '')::VARCHAR || '%') OR @song_title = '');