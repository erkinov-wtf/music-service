// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0

package database

import (
	"github.com/jackc/pgx/v5/pgtype"
)

type Group struct {
	ID        pgtype.UUID
	Name      string
	CreatedAt pgtype.Timestamptz
	UpdatedAt pgtype.Timestamptz
	DeletedAt pgtype.Timestamptz
}

type Song struct {
	ID          pgtype.UUID
	GroupID     pgtype.UUID
	Title       string
	Runtime     int32
	Lyrics      []byte
	ReleaseDate pgtype.Timestamptz
	Link        string
	CreatedAt   pgtype.Timestamptz
	UpdatedAt   pgtype.Timestamptz
	DeletedAt   pgtype.Timestamptz
}
