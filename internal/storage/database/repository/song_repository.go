package repository

import (
	"context"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"music-service/internal/storage/database"
	"time"
)

type SongRepositoryInterface interface {
	CreateSong(ctx context.Context, params SongCreateParams) (database.Song, error)
	GetSong(ctx context.Context, id uuid.UUID) (database.Song, error)
	GetSongsCount(ctx context.Context) (int64, error)
	GetSongsWithPagination(ctx context.Context, limit, offset int32) ([]database.GetSongsWithPaginationRow, error)
	UpdateSong(ctx context.Context, params SongUpdateParams) (database.Song, error)
	GetSongsByGroup(ctx context.Context, groupID uuid.UUID, limit, offset int32) ([]database.Song, error)
	GetSongsWithFilters(ctx context.Context, params SongFilterParams) ([]database.GetSongsWithPaginationRow, error)
	GetSongsCountWithFilters(ctx context.Context, groupName, songTitle string) (int64, error)
	DeleteSong(ctx context.Context, id uuid.UUID) error
}

type SongCreateParams struct {
	GroupID     uuid.UUID
	Title       string
	Runtime     int32
	Lyrics      []byte
	ReleaseDate time.Time
	Link        string
}

type SongUpdateParams struct {
	ID          uuid.UUID
	GroupID     uuid.UUID
	Title       string
	Runtime     int32
	Lyrics      []byte
	ReleaseDate time.Time
	Link        string
}

type SongFilterParams struct {
	Limit     int32
	Offset    int32
	GroupName string
	SongTitle string
}

type SongRepository struct {
	q *database.Queries
}

func NewSongRepository(db database.DBTX) SongRepositoryInterface {
	return &SongRepository{
		q: database.New(db),
	}
}

func (r *SongRepository) CreateSong(ctx context.Context, params SongCreateParams) (database.Song, error) {
	pgGroupID := pgtype.UUID{Bytes: params.GroupID, Valid: true}
	pgReleaseDate := pgtype.Timestamptz{Time: params.ReleaseDate, Valid: true}

	return r.q.CreateSong(ctx, database.CreateSongParams{
		GroupID:     pgGroupID,
		Title:       params.Title,
		Runtime:     params.Runtime,
		Lyrics:      params.Lyrics,
		ReleaseDate: pgReleaseDate,
		Link:        params.Link,
	})

}

func (r *SongRepository) GetSong(ctx context.Context, id uuid.UUID) (database.Song, error) {
	pgID := pgtype.UUID{Bytes: id, Valid: true}
	return r.q.GetSong(ctx, pgID)
}

func (r *SongRepository) GetSongsCount(ctx context.Context) (int64, error) {
	return r.q.GetSongsCount(ctx)
}

func (r *SongRepository) GetSongsWithPagination(ctx context.Context, limit, offset int32) ([]database.GetSongsWithPaginationRow, error) {
	return r.q.GetSongsWithPagination(ctx, database.GetSongsWithPaginationParams{
		Limit:  limit,
		Offset: offset,
	})
}

func (r *SongRepository) UpdateSong(ctx context.Context, params SongUpdateParams) (database.Song, error) {
	pgID := pgtype.UUID{Bytes: params.ID, Valid: true}
	pgGroupID := pgtype.UUID{Bytes: params.GroupID, Valid: true}
	pgReleaseDate := pgtype.Timestamptz{Time: params.ReleaseDate, Valid: true}

	return r.q.UpdateSong(ctx, database.UpdateSongParams{
		ID:          pgID,
		GroupID:     pgGroupID,
		Title:       params.Title,
		Runtime:     params.Runtime,
		Lyrics:      params.Lyrics,
		ReleaseDate: pgReleaseDate,
		Link:        params.Link,
	})
}

func (r *SongRepository) DeleteSong(ctx context.Context, id uuid.UUID) error {
	pgID := pgtype.UUID{Bytes: id, Valid: true}
	_, err := r.q.DeleteSong(ctx, pgID)
	return err
}

func (r *SongRepository) GetSongsByGroup(ctx context.Context, groupID uuid.UUID, limit, offset int32) ([]database.Song, error) {
	pgGroupID := pgtype.UUID{Bytes: groupID, Valid: true}
	return r.q.GetSongsByGroup(ctx, database.GetSongsByGroupParams{
		GroupID: pgGroupID,
		Limit:   limit,
		Offset:  offset,
	})
}

func (r *SongRepository) GetSongsWithFilters(ctx context.Context, params SongFilterParams) ([]database.GetSongsWithPaginationRow, error) {
	rows, err := r.q.GetSongsWithFilters(ctx, database.GetSongsWithFiltersParams{
		Limit:   params.Limit,
		Offset:  params.Offset,
		Column3: params.GroupName,
		Column4: params.SongTitle,
	})
	if err != nil {
		return nil, err
	}

	var songs []database.GetSongsWithPaginationRow
	for _, row := range rows {
		songs = append(songs, database.GetSongsWithPaginationRow(row))
	}
	return songs, nil
}

func (r *SongRepository) GetSongsCountWithFilters(ctx context.Context, groupName, songTitle string) (int64, error) {
	return r.q.GetSongsCountWithFilters(ctx, database.GetSongsCountWithFiltersParams{
		GroupName: groupName,
		SongTitle: songTitle,
	})
}
