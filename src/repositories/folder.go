package repositories

import (
	"context"

	"github.com/CustomCloudStorage/types"
	"github.com/CustomCloudStorage/utils"
)

func (r *folderRepository) Create(ctx context.Context, folder *types.Folder) error {
	if err := r.db.WithContext(ctx).Create(folder).Error; err != nil {
		return utils.DetermineSQLError(err, "create data")
	}
	return nil
}

func (r *folderRepository) GetByID(ctx context.Context, id int, userID int) (*types.Folder, error) {
	var folder types.Folder
	if err := r.db.WithContext(ctx).
		Where("id = ? AND user_id = ?", id, userID).
		First(&folder).Error; err != nil {
		return nil, utils.DetermineSQLError(err, "get data")
	}
	return &folder, nil
}

func (r *folderRepository) Update(ctx context.Context, folder *types.Folder) error {
	if err := r.db.Save(folder).Error; err != nil {
		return utils.DetermineSQLError(err, "update data")
	}
	return nil
}

func (r *folderRepository) ListByUserID(ctx context.Context, userID int) ([]types.Folder, error) {
	var folders []types.Folder
	if err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Find(&folders).Error; err != nil {
		return nil, utils.DetermineSQLError(err, "get data")
	}
	return folders, nil
}
