package services

import (
	"context"
	"github.com/google/uuid"
	"music-service/internal/storage/database"
	"music-service/internal/storage/database/repository"
)

// SongService handles business logic for songs
type SongService struct {
	songRepo repository.SongRepositoryInterface
}

// NewSongService creates a new song service
func NewSongService(songRepo repository.SongRepositoryInterface) *SongService {
	return &SongService{
		songRepo: songRepo,
	}
}

func (s *SongService) CreateSong(ctx context.Context, params repository.SongCreateParams) (database.Song, error) {
	return s.songRepo.CreateSong(ctx, params)
}

func (s *SongService) GetSong(ctx context.Context, id uuid.UUID) (database.Song, error) {
	return s.songRepo.GetSong(ctx, id)
}

func (s *SongService) GetSongsCount(ctx context.Context) (int64, error) {
	return s.songRepo.GetSongsCount(ctx)
}

func (s *SongService) GetSongsWithPagination(ctx context.Context, limit, offset int32) ([]database.GetSongsWithPaginationRow, error) {
	return s.songRepo.GetSongsWithPagination(ctx, limit, offset)
}

func (s *SongService) UpdateSong(ctx context.Context, params repository.SongUpdateParams) (database.Song, error) {
	return s.songRepo.UpdateSong(ctx, params)
}

func (s *SongService) GetSongsByGroup(ctx context.Context, groupID uuid.UUID, limit, offset int32) ([]database.Song, error) {
	return s.songRepo.GetSongsByGroup(ctx, groupID, limit, offset)
}

func (s *SongService) GetSongsWithFilters(ctx context.Context, params repository.SongFilterParams) ([]database.GetSongsWithPaginationRow, error) {
	return s.songRepo.GetSongsWithFilters(ctx, params)
}

func (s *SongService) GetSongsCountWithFilters(ctx context.Context, groupName, songTitle string) (int64, error) {
	return s.songRepo.GetSongsCountWithFilters(ctx, groupName, songTitle)
}

func (s *SongService) DeleteSong(ctx context.Context, id uuid.UUID) error {
	return s.songRepo.DeleteSong(ctx, id)
}
