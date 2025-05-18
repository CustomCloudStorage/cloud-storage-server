package repositories

import (
	"context"

	"github.com/CustomCloudStorage/types"
	"github.com/CustomCloudStorage/utils"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UploadPartRepository interface {
	Create(ctx context.Context, part *types.UploadPart) error
	ListBySession(ctx context.Context, sessionID uuid.UUID) ([]types.UploadPart, error)
	DeleteBySession(ctx context.Context, sessionID uuid.UUID) error
}

type uploadPartRepository struct {
	db *gorm.DB
}

func NewUploadPartRepository(db *gorm.DB) UploadPartRepository {
	return &uploadPartRepository{
		db: db,
	}
}

func (r *uploadPartRepository) Create(ctx context.Context, part *types.UploadPart) error {
	if err := r.db.WithContext(ctx).
		Create(part).
		Error; err != nil {
		return utils.DetermineSQLError(err, "create upload part")
	}
	return nil
}

func (r *uploadPartRepository) ListBySession(ctx context.Context, sessionID uuid.UUID) ([]types.UploadPart, error) {
	var parts []types.UploadPart
	if err := r.db.WithContext(ctx).
		Find(&parts, "session_id = ?", sessionID).
		Error; err != nil {
		return nil, utils.DetermineSQLError(err, "list upload parts")
	}
	return parts, nil
}

func (r *uploadPartRepository) DeleteBySession(ctx context.Context, sessionID uuid.UUID) error {
	if err := r.db.WithContext(ctx).
		Where("session_id = ?", sessionID).
		Delete(&types.UploadPart{}).
		Error; err != nil {
		return utils.DetermineSQLError(err, "delete upload parts")
	}
	return nil
}
