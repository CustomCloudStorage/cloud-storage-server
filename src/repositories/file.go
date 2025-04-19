package repositories

import (
	"context"

	"github.com/CustomCloudStorage/types"
	"github.com/CustomCloudStorage/utils"
)

func (f *file) Create(ctx context.Context, file *types.File) error {
	if err := f.db.WithContext(ctx).Create(file).Error; err != nil {
		return utils.DetermineSQLError(err, "create data")
	}
	return nil
}

func (f *file) GetByID(ctx context.Context, id int, userID int) (*types.File, error) {
	var file types.File
	if err := f.db.WithContext(ctx).
		Where("id = ? AND user_id = ?", id, userID).
		First(&file).Error; err != nil {
		return nil, utils.DetermineSQLError(err, "get data")
	}
	return &file, nil
}

func (f *file) Update(ctx context.Context, file *types.File) error {
	if err := f.db.Save(file).Error; err != nil {
		return utils.DetermineSQLError(err, "update data")
	}
	return nil
}

func (f *file) Delete(ctx context.Context, id int, userID int) error {
	if err := f.db.WithContext(ctx).
		Where("id = ? AND user_id = ?", id, userID).
		Delete(&types.File{}).Error; err != nil {
		return utils.DetermineSQLError(err, "delete data")
	}
	return nil
}

func (f *file) ListByUserID(ctx context.Context, userID int) ([]types.File, error) {
	var files []types.File
	if err := f.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Find(&files).Error; err != nil {
		return nil, utils.DetermineSQLError(err, "get data")
	}
	return files, nil
}

func (f *file) UpdateName(ctx context.Context, id int, userID int, name string) error {
	if err := f.db.WithContext(ctx).
		Model(&types.File{}).
		Where("id = ? AND user_id = ?", id, userID).
		Update("name", name).Error; err != nil {
		return utils.DetermineSQLError(err, "update file name")
	}
	return nil
}

func (f *file) UpdateFolder(ctx context.Context, id int, userID int, folderID int) error {
	if err := f.db.WithContext(ctx).
		Model(&types.File{}).
		Where("id = & AND user_id = ?", id, userID).
		Update("folder_id", folderID).Error; err != nil {
		return utils.DetermineSQLError(err, "update file folder")
	}
	return nil
}
