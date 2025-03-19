package repository

import (
	"context"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"music-service/internal/storage/database"
)

type GroupRepositoryInterface interface {
	CreateGroup(ctx context.Context, name string) error
	GetGroup(ctx context.Context, id uuid.UUID) (database.Group, error)
	GetGroupsCount(ctx context.Context) (int64, error)
	GetGroupsWithPagination(ctx context.Context, limit, offset int32) ([]database.GetGroupsWithPaginationRow, error)
	UpdateGroup(ctx context.Context, id uuid.UUID, name string) error
	DeleteGroup(ctx context.Context, id uuid.UUID) error
}

type GroupRepository struct {
	q *database.Queries
}

func NewGroupRepository(db database.DBTX) GroupRepositoryInterface {
	return &GroupRepository{
		q: database.New(db),
	}
}

func (r *GroupRepository) CreateGroup(ctx context.Context, name string) error {
	_, err := r.q.CreateGroup(ctx, name)
	return err
}

func (r *GroupRepository) DeleteGroup(ctx context.Context, id uuid.UUID) error {
	pgID := pgtype.UUID{Bytes: id, Valid: true}
	return r.q.DeleteGroup(ctx, pgID)
}

func (r *GroupRepository) GetGroup(ctx context.Context, id uuid.UUID) (database.Group, error) {
	pgID := pgtype.UUID{Bytes: id, Valid: true}
	return r.q.GetGroup(ctx, pgID)
}

func (r *GroupRepository) GetGroupsCount(ctx context.Context) (int64, error) {
	return r.q.GetGroupsCount(ctx)
}

func (r *GroupRepository) GetGroupsWithPagination(ctx context.Context, limit, offset int32) ([]database.GetGroupsWithPaginationRow, error) {
	return r.q.GetGroupsWithPagination(ctx, database.GetGroupsWithPaginationParams{
		Limit:  limit,
		Offset: offset,
	})
}

func (r *GroupRepository) UpdateGroup(ctx context.Context, id uuid.UUID, name string) error {
	pgID := pgtype.UUID{Bytes: id, Valid: true}
	_, err := r.q.UpdateGroup(ctx, database.UpdateGroupParams{
		ID:   pgID,
		Name: name,
	})
	return err
}
