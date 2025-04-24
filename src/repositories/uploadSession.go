package repositories

import (
	"context"

	"github.com/CustomCloudStorage/types"
	"github.com/CustomCloudStorage/utils"
	"github.com/google/uuid"
)

func (r *uploadSessionRepository) Create(ctx context.Context, session *types.UploadSession) error {
	if err := r.db.WithContext(ctx).
		Create(session).
		Error; err != nil {
		return utils.DetermineSQLError(err, "create data")
	}
	return nil
}

func (r *uploadSessionRepository) GetByID(ctx context.Context, id uuid.UUID) (*types.UploadSession, error) {
	var session types.UploadSession
	if err := r.db.WithContext(ctx).
		First(&session, "id = ?", id).
		Error; err != nil {
		return nil, utils.DetermineSQLError(err, "get data")
	}
	return &session, nil
}

func (r *uploadSessionRepository) Delete(ctx context.Context, id uuid.UUID) error {
	if err := r.db.WithContext(ctx).
		Delete(&types.UploadSession{}, "id = ?", id).
		Error; err != nil {
		return utils.DetermineSQLError(err, "delete data")
	}
	return nil
}

func (r *uploadSessionRepository) ListOlderThan(ctx context.Context, olderThanMinutes int) ([]types.UploadSession, error) {
	var sessions []types.UploadSession
	if err := r.db.WithContext(ctx).
		Where("created_at < NOW() - INTERVAL '? minutes'", olderThanMinutes).
		Find(&sessions).Error; err != nil {
		return nil, utils.DetermineSQLError(err, "list data")
	}
	return sessions, nil
}
