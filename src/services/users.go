package services

import (
	"context"

	"github.com/CustomCloudStorage/repositories"
	"github.com/CustomCloudStorage/types"
)

type UserService interface {
	StatsStorage(ctx context.Context) (*types.StorageStats, error)
}

type userService struct {
	userRepository repositories.UserRepository
	cfg            ServiceConfig
}

func NewUserService(userRepo repositories.UserRepository, cfg ServiceConfig) UserService {
	return &userService{
		userRepository: userRepo,
		cfg:            cfg,
	}
}

func (s *userService) StatsStorage(ctx context.Context) (*types.StorageStats, error) {
	allocated, err := s.userRepository.SumActiveStorageLimit(ctx)
	if err != nil {
		return nil, err
	}
	total := s.cfg.TotalStorageBytes()
	free := total - allocated
	stats := types.StorageStats{
		TotalBytes:     total,
		AllocatedBytes: allocated,
		FreeBytes:      free,
	}
	return &stats, nil
}
