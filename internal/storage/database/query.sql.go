// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: query.sql

package database

import (
	"context"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
)

const createGroup = `-- name: CreateGroup :execresult

INSERT INTO groups (name)
VALUES ($1)
`

// Groups Table
func (q *Queries) CreateGroup(ctx context.Context, name string) (pgconn.CommandTag, error) {
	return q.db.Exec(ctx, createGroup, name)
}

const createSong = `-- name: CreateSong :execresult

INSERT INTO songs (group_id, title, runtime, lyrics, release_date, link)
VALUES ($1, $2, $3, $4, $5, $6)
`

type CreateSongParams struct {
	GroupID     pgtype.UUID
	Title       string
	Runtime     int32
	Lyrics      []byte
	ReleaseDate pgtype.Timestamptz
	Link        string
}

// Songs Table
func (q *Queries) CreateSong(ctx context.Context, arg CreateSongParams) (pgconn.CommandTag, error) {
	return q.db.Exec(ctx, createSong,
		arg.GroupID,
		arg.Title,
		arg.Runtime,
		arg.Lyrics,
		arg.ReleaseDate,
		arg.Link,
	)
}

const deleteGroup = `-- name: DeleteGroup :exec
UPDATE groups
SET deleted_at = NOW()
WHERE id = $1 AND deleted_at IS NULL
`

func (q *Queries) DeleteGroup(ctx context.Context, id pgtype.UUID) error {
	_, err := q.db.Exec(ctx, deleteGroup, id)
	return err
}

const deleteSong = `-- name: DeleteSong :execresult
UPDATE songs
SET deleted_at = NOW()
WHERE id = $1 AND deleted_at IS NULL
`

func (q *Queries) DeleteSong(ctx context.Context, id pgtype.UUID) (pgconn.CommandTag, error) {
	return q.db.Exec(ctx, deleteSong, id)
}

const getGroup = `-- name: GetGroup :one
SELECT id, name, created_at, updated_at, deleted_at FROM groups
WHERE id = $1 LIMIT 1
`

func (q *Queries) GetGroup(ctx context.Context, id pgtype.UUID) (Group, error) {
	row := q.db.QueryRow(ctx, getGroup, id)
	var i Group
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.DeletedAt,
	)
	return i, err
}

const getGroupsCount = `-- name: GetGroupsCount :one
SELECT count(*) FROM groups
WHERE deleted_at IS NULL
`

func (q *Queries) GetGroupsCount(ctx context.Context) (int64, error) {
	row := q.db.QueryRow(ctx, getGroupsCount)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const getGroupsWithPagination = `-- name: GetGroupsWithPagination :many
SELECT id, name, created_at, updated_at FROM groups
WHERE deleted_at IS NULL
ORDER BY created_at DESC LIMIT $1 OFFSET $2
`

type GetGroupsWithPaginationParams struct {
	Limit  int32
	Offset int32
}

type GetGroupsWithPaginationRow struct {
	ID        pgtype.UUID
	Name      string
	CreatedAt pgtype.Timestamptz
	UpdatedAt pgtype.Timestamptz
}

func (q *Queries) GetGroupsWithPagination(ctx context.Context, arg GetGroupsWithPaginationParams) ([]GetGroupsWithPaginationRow, error) {
	rows, err := q.db.Query(ctx, getGroupsWithPagination, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetGroupsWithPaginationRow
	for rows.Next() {
		var i GetGroupsWithPaginationRow
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getSong = `-- name: GetSong :one
SELECT id, group_id, title, runtime,  lyrics, release_date, link, created_at, updated_at, deleted_at
FROM songs
WHERE id = $1 LIMIT 1
`

func (q *Queries) GetSong(ctx context.Context, id pgtype.UUID) (Song, error) {
	row := q.db.QueryRow(ctx, getSong, id)
	var i Song
	err := row.Scan(
		&i.ID,
		&i.GroupID,
		&i.Title,
		&i.Runtime,
		&i.Lyrics,
		&i.ReleaseDate,
		&i.Link,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.DeletedAt,
	)
	return i, err
}

const getSongsByGroup = `-- name: GetSongsByGroup :many
SELECT id, group_id, title, runtime, lyrics, release_date, link, created_at, updated_at, deleted_at
FROM songs
WHERE group_id = $1 AND deleted_at IS NULL
ORDER BY release_date DESC LIMIT $2 OFFSET $3
`

type GetSongsByGroupParams struct {
	GroupID pgtype.UUID
	Limit   int32
	Offset  int32
}

func (q *Queries) GetSongsByGroup(ctx context.Context, arg GetSongsByGroupParams) ([]Song, error) {
	rows, err := q.db.Query(ctx, getSongsByGroup, arg.GroupID, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Song
	for rows.Next() {
		var i Song
		if err := rows.Scan(
			&i.ID,
			&i.GroupID,
			&i.Title,
			&i.Runtime,
			&i.Lyrics,
			&i.ReleaseDate,
			&i.Link,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.DeletedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getSongsCount = `-- name: GetSongsCount :one
SELECT count(*) FROM songs
WHERE deleted_at IS NULL
`

func (q *Queries) GetSongsCount(ctx context.Context) (int64, error) {
	row := q.db.QueryRow(ctx, getSongsCount)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const getSongsCountWithFilters = `-- name: GetSongsCountWithFilters :one
SELECT count(*)
FROM songs s
JOIN groups g ON s.group_id = g.id
WHERE s.deleted_at IS NULL
  AND (LOWER(g.name) LIKE LOWER('%' || NULLIF($1, '')::VARCHAR || '%') OR $1 = '')
  AND (LOWER(s.title) LIKE LOWER('%' || NULLIF($2, '')::VARCHAR || '%') OR $2 = '')
`

type GetSongsCountWithFiltersParams struct {
	GroupName interface{}
	SongTitle interface{}
}

func (q *Queries) GetSongsCountWithFilters(ctx context.Context, arg GetSongsCountWithFiltersParams) (int64, error) {
	row := q.db.QueryRow(ctx, getSongsCountWithFilters, arg.GroupName, arg.SongTitle)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const getSongsWithFilters = `-- name: GetSongsWithFilters :many
SELECT s.id, s.group_id, s.title, s.runtime, s.lyrics, s.release_date, s.link, s.created_at,  s.updated_at
FROM songs s
JOIN groups g ON s.group_id = g.id
WHERE s.deleted_at IS NULL
  AND (LOWER(g.name) LIKE LOWER('%' || NULLIF($3, '')::VARCHAR || '%') OR $3 = '')
  AND (LOWER(s.title) LIKE LOWER('%' || NULLIF($4, '')::VARCHAR || '%') OR $4 = '')
ORDER BY s.created_at DESC
    LIMIT $1 OFFSET $2
`

type GetSongsWithFiltersParams struct {
	Limit     int32
	Offset    int32
	GroupName interface{}
	SongTitle interface{}
}

type GetSongsWithFiltersRow struct {
	ID          pgtype.UUID
	GroupID     pgtype.UUID
	Title       string
	Runtime     int32
	Lyrics      []byte
	ReleaseDate pgtype.Timestamptz
	Link        string
	CreatedAt   pgtype.Timestamptz
	UpdatedAt   pgtype.Timestamptz
}

func (q *Queries) GetSongsWithFilters(ctx context.Context, arg GetSongsWithFiltersParams) ([]GetSongsWithFiltersRow, error) {
	rows, err := q.db.Query(ctx, getSongsWithFilters,
		arg.Limit,
		arg.Offset,
		arg.GroupName,
		arg.SongTitle,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetSongsWithFiltersRow
	for rows.Next() {
		var i GetSongsWithFiltersRow
		if err := rows.Scan(
			&i.ID,
			&i.GroupID,
			&i.Title,
			&i.Runtime,
			&i.Lyrics,
			&i.ReleaseDate,
			&i.Link,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getSongsWithPagination = `-- name: GetSongsWithPagination :many
SELECT id, group_id, title, runtime, lyrics, release_date, link, created_at, updated_at FROM songs
WHERE deleted_at IS NULL
ORDER BY created_at DESC LIMIT $1 OFFSET $2
`

type GetSongsWithPaginationParams struct {
	Limit  int32
	Offset int32
}

type GetSongsWithPaginationRow struct {
	ID          pgtype.UUID
	GroupID     pgtype.UUID
	Title       string
	Runtime     int32
	Lyrics      []byte
	ReleaseDate pgtype.Timestamptz
	Link        string
	CreatedAt   pgtype.Timestamptz
	UpdatedAt   pgtype.Timestamptz
}

func (q *Queries) GetSongsWithPagination(ctx context.Context, arg GetSongsWithPaginationParams) ([]GetSongsWithPaginationRow, error) {
	rows, err := q.db.Query(ctx, getSongsWithPagination, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetSongsWithPaginationRow
	for rows.Next() {
		var i GetSongsWithPaginationRow
		if err := rows.Scan(
			&i.ID,
			&i.GroupID,
			&i.Title,
			&i.Runtime,
			&i.Lyrics,
			&i.ReleaseDate,
			&i.Link,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateGroup = `-- name: UpdateGroup :execresult
UPDATE groups
SET name = $2
WHERE id = $1
`

type UpdateGroupParams struct {
	ID   pgtype.UUID
	Name string
}

func (q *Queries) UpdateGroup(ctx context.Context, arg UpdateGroupParams) (pgconn.CommandTag, error) {
	return q.db.Exec(ctx, updateGroup, arg.ID, arg.Name)
}

const updateSong = `-- name: UpdateSong :execresult
UPDATE songs
SET
    group_id = $2,
    title = $3,
    runtime = $4,
    lyrics = $5,
    release_date = $6,
    link = $7
WHERE id = $1
`

type UpdateSongParams struct {
	ID          pgtype.UUID
	GroupID     pgtype.UUID
	Title       string
	Runtime     int32
	Lyrics      []byte
	ReleaseDate pgtype.Timestamptz
	Link        string
}

func (q *Queries) UpdateSong(ctx context.Context, arg UpdateSongParams) (pgconn.CommandTag, error) {
	return q.db.Exec(ctx, updateSong,
		arg.ID,
		arg.GroupID,
		arg.Title,
		arg.Runtime,
		arg.Lyrics,
		arg.ReleaseDate,
		arg.Link,
	)
}
