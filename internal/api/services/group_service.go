package services

import (
	"context"
	"github.com/google/uuid"
	"music-service/internal/storage/database"
	"music-service/internal/storage/database/repository"
)

// GroupService handles business logic for groups
type GroupService struct {
	groupRepo repository.GroupRepositoryInterface
}

// NewGroupService creates a new group service
func NewGroupService(groupRepo repository.GroupRepositoryInterface) *GroupService {
	return &GroupService{
		groupRepo: groupRepo,
	}
}

func (s *GroupService) CreateGroup(ctx context.Context, name string) error {
	return s.groupRepo.CreateGroup(ctx, name)
}

func (s *GroupService) DeleteGroup(ctx context.Context, id uuid.UUID) error {
	return s.groupRepo.DeleteGroup(ctx, id)
}

func (s *GroupService) GetGroup(ctx context.Context, id uuid.UUID) (database.Group, error) {
	return s.groupRepo.GetGroup(ctx, id)
}

func (s *GroupService) GetGroupsCount(ctx context.Context) (int64, error) {
	return s.groupRepo.GetGroupsCount(ctx)
}

func (s *GroupService) GetGroupsWithPagination(ctx context.Context, limit, offset int32) ([]database.GetGroupsWithPaginationRow, error) {
	return s.groupRepo.GetGroupsWithPagination(ctx, limit, offset)
}

func (s *GroupService) UpdateGroup(ctx context.Context, id uuid.UUID, name string) error {
	return s.groupRepo.UpdateGroup(ctx, id, name)
}
