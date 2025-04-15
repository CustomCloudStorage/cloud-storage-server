package repositories

import (
	"context"

	"github.com/CustomCloudStorage/types"
	"github.com/CustomCloudStorage/utils"
)

func (f *folder) Create(ctx context.Context, folder *types.Folder) error {
	if err := f.db.WithContext(ctx).Create(folder).Error; err != nil {
		return utils.DetermineSQLError(err, "create data")
	}
	return nil
}

func (f *folder) GetByID(ctx context.Context, id int, userID int) (*types.Folder, error) {
	var folder types.Folder
	if err := f.db.WithContext(ctx).
		Where("id = ? AND user_id = ?", id, userID).
		First(&folder).Error; err != nil {
		return nil, utils.DetermineSQLError(err, "get data")
	}
	return &folder, nil
}

func (f *folder) Update(ctx context.Context, folder *types.Folder) error {
	if err := f.db.Save(folder).Error; err != nil {
		return utils.DetermineSQLError(err, "update data")
	}
	return nil
}

func (f *folder) Delete(ctx context.Context, id int, userID int) error {
	if err := f.db.WithContext(ctx).
		Where("id = ? AND user_id = ?", id, userID).
		Delete(&types.Folder{}).Error; err != nil {
		return utils.DetermineSQLError(err, "delete data")
	}
	return nil
}

func (f *folder) ListByUserID(ctx context.Context, userID int) ([]types.Folder, error) {
	var folders []types.Folder
	if err := f.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Find(&folders).Error; err != nil {
		return nil, utils.DetermineSQLError(err, "get data")
	}
	return folders, nil
}
