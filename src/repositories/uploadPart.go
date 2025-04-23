package repositories

import (
	"context"

	"github.com/CustomCloudStorage/types"
	"github.com/CustomCloudStorage/utils"
	"github.com/google/uuid"
)

func (r *uploadPart) Create(ctx context.Context, part *types.UploadPart) error {
	if err := r.db.WithContext(ctx).
		Create(part).
		Error; err != nil {
		return utils.DetermineSQLError(err, "create data")
	}
	return nil
}

func (r *uploadPart) ListBySession(ctx context.Context, sessionID uuid.UUID) ([]types.UploadPart, error) {
	var parts []types.UploadPart
	if err := r.db.WithContext(ctx).
		Find(&parts, "session_id = ?", sessionID).
		Error; err != nil {
		return nil, utils.DetermineSQLError(err, "create data")
	}
	return parts, nil
}

func (r *uploadPart) DeleteBySession(ctx context.Context, sessionID uuid.UUID) error {
	if err := r.db.WithContext(ctx).
		Where("session_id = ?", sessionID).
		Delete(&types.UploadPart{}).
		Error; err != nil {
		return utils.DetermineSQLError(err, "delete data")
	}
	return nil
}
