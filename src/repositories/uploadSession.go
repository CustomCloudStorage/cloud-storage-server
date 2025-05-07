package repositories

import (
	"context"
	"time"

	"github.com/CustomCloudStorage/types"
	"github.com/CustomCloudStorage/utils"
	"github.com/google/uuid"
)

func (r *uploadSessionRepository) Create(ctx context.Context, session *types.UploadSession) error {
	if err := r.db.WithContext(ctx).
		Create(session).
		Error; err != nil {
		return utils.DetermineSQLError(err, "create upload session")
	}
	return nil
}

func (r *uploadSessionRepository) GetByID(ctx context.Context, id uuid.UUID) (*types.UploadSession, error) {
	var session types.UploadSession
	if err := r.db.WithContext(ctx).
		First(&session, "id = ?", id).
		Error; err != nil {
		return nil, utils.DetermineSQLError(err, "get upload session")
	}
	return &session, nil
}

func (r *uploadSessionRepository) Delete(ctx context.Context, id uuid.UUID) error {
	if err := r.db.WithContext(ctx).
		Unscoped().
		Where("id = ?", id).
		Delete(&types.UploadSession{}).
		Error; err != nil {
		return utils.DetermineSQLError(err, "delete upload session")
	}
	return nil
}

func (r *uploadSessionRepository) ListOlderThan(ctx context.Context, olderThan time.Duration) ([]types.UploadSession, error) {
	cutoff := time.Now().Add(-olderThan)
	var sessions []types.UploadSession
	if err := r.db.WithContext(ctx).
		Where("created_at < ?", cutoff).
		Find(&sessions).
		Error; err != nil {
		return nil, utils.DetermineSQLError(err, "list old upload sessions")
	}
	return sessions, nil
}
