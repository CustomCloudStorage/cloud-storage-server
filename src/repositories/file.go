package repositories

import (
	"context"

	"github.com/CustomCloudStorage/types"
	"github.com/CustomCloudStorage/utils"
)

func (r *fileRepository) Create(ctx context.Context, file *types.File) error {
	if err := r.db.WithContext(ctx).Create(file).Error; err != nil {
		return utils.DetermineSQLError(err, "create data")
	}
	return nil
}

func (r *fileRepository) GetByID(ctx context.Context, id int, userID int) (*types.File, error) {
	var file types.File
	if err := r.db.WithContext(ctx).
		Where("id = ? AND user_id = ?", id, userID).
		First(&file).Error; err != nil {
		return nil, utils.DetermineSQLError(err, "get data")
	}
	return &file, nil
}

func (r *fileRepository) Update(ctx context.Context, file *types.File) error {
	if err := r.db.Save(file).Error; err != nil {
		return utils.DetermineSQLError(err, "update data")
	}
	return nil
}

func (r *fileRepository) ListByUserID(ctx context.Context, userID int) ([]types.File, error) {
	var files []types.File
	if err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Find(&files).Error; err != nil {
		return nil, utils.DetermineSQLError(err, "get data")
	}
	return files, nil
}

func (r *fileRepository) UpdateName(ctx context.Context, id int, userID int, name string) error {
	if err := r.db.WithContext(ctx).
		Model(&types.File{}).
		Where("id = ? AND user_id = ?", id, userID).
		Update("name", name).Error; err != nil {
		return utils.DetermineSQLError(err, "update file name")
	}
	return nil
}

func (r *fileRepository) UpdateFolder(ctx context.Context, id int, userID int, folderID int) error {
	if err := r.db.WithContext(ctx).
		Model(&types.File{}).
		Where("id = & AND user_id = ?", id, userID).
		Update("folder_id", folderID).Error; err != nil {
		return utils.DetermineSQLError(err, "update file folder")
	}
	return nil
}

func (r *fileRepository) ListFilesRecursive(ctx context.Context, userID, folderID int) ([]*types.FileWithPath, error) {
	const sql = `
WITH RECURSIVE folders_cte AS (
    SELECT id, name, parent_id, name AS path
	FROM folders
	WHERE user_id = @user AND id = @root

UNION ALL

    SELECT f.id, f.name, f.parent_id, folders_cte.path || '/' || f.name
	FROM folders f
	JOIN folders_cte ON f.parent_id = folders_cte.id
	WHERE f.user_id = @user
)
SELECT
    fi.physical_name,
    folders_cte.path || '/' || fi.name || fi.extension AS relative_path
	FROM folders_cte
	JOIN files fi ON fi.folder_id = folders_cte.id
	WHERE fi.user_id = @user;
`
	var out []*types.FileWithPath
	if err := r.db.WithContext(ctx).
		Raw(sql,
			map[string]interface{}{"user": userID, "root": folderID},
		).
		Scan(&out).Error; err != nil {
		return nil, err
	}
	return out, nil
}
