package repositories

import (
	"context"
	"time"

	"github.com/CustomCloudStorage/types"
	"github.com/CustomCloudStorage/utils"
	"gorm.io/gorm"
)

type TrashRepository interface {
	SoftDeleteFile(ctx context.Context, userID, fileID int, ts time.Time) error
	RestoreFile(ctx context.Context, userID, fileID int) error
	ListTrashedFiles(ctx context.Context, userID int) ([]*types.File, error)
	ListFilesToPurge(ctx context.Context, before time.Time) ([]*types.File, error)
	HardDeleteFileByID(ctx context.Context, fileID int) error
	PermanentDeleteFile(ctx context.Context, userID, fileID int) error

	SoftDeleteFolderCascade(ctx context.Context, userID, folderID int, ts time.Time) error
	RestoreFolderCascade(ctx context.Context, userID, folderID int) error
	ListTrashedFolders(ctx context.Context, userID int) ([]*types.Folder, error)
	ListFoldersToPurge(ctx context.Context, before time.Time) ([]*types.Folder, error)
	HardDeleteFolderByID(ctx context.Context, folderID int) error
	PermanentDeleteFolder(ctx context.Context, userID, folderID int) ([]int, error)
}

type trashRepository struct {
	db *gorm.DB
}

func NewTrashRepository(db *gorm.DB) TrashRepository {
	return &trashRepository{
		db: db,
	}
}

func (r *trashRepository) SoftDeleteFile(ctx context.Context, userID, fileID int, ts time.Time) error {
	if err := r.db.WithContext(ctx).
		Model(&types.File{}).
		Where("id = & AND user_id = ?", fileID, userID).
		Update("deleted_at", ts).
		Error; err != nil {
		return utils.DetermineSQLError(err, "soft delete file")
	}
	return nil
}

func (r *trashRepository) RestoreFile(ctx context.Context, userID, fileID int) error {
	if err := r.db.WithContext(ctx).
		Model(&types.File{}).
		Where("id = ? AND user_id = ?", fileID, userID).
		Update("deleted_at", nil).
		Error; err != nil {
		return utils.DetermineSQLError(err, "restore file")
	}
	return nil
}

func (r *trashRepository) ListTrashedFiles(ctx context.Context, userID int) ([]*types.File, error) {
	var out []*types.File
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND deleted_at IS NOT NULL", userID).
		Find(&out).Error
	if err != nil {
		return nil, utils.DetermineSQLError(err, "list trashed files")
	}
	return out, nil
}

func (r *trashRepository) ListFilesToPurge(ctx context.Context, before time.Time) ([]*types.File, error) {
	var out []*types.File
	err := r.db.WithContext(ctx).
		Where("deleted_at IS NOT NULL AND deleted_at < ?", before).
		Find(&out).Error
	if err != nil {
		return nil, utils.DetermineSQLError(err, "list files to purge")
	}
	return out, err
}

func (r *trashRepository) HardDeleteFileByID(ctx context.Context, fileID int) error {
	if err := r.db.WithContext(ctx).
		Unscoped().
		Where("id = ?", fileID).
		Delete(&types.File{}).
		Error; err != nil {
		return utils.DetermineSQLError(err, "hard delete file")
	}
	return nil
}

func (r *trashRepository) PermanentDeleteFile(ctx context.Context, userID, fileID int) error {
	var file types.File
	if err := r.db.WithContext(ctx).
		Unscoped().
		Where("id = ? AND user_id = ? AND deleted_at IS NOT NULL", fileID, userID).
		First(&file).Error; err != nil {
		return utils.DetermineSQLError(err, "get trashed file")
	}
	if err := r.db.WithContext(ctx).
		Unscoped().
		Where("id = ?", fileID).
		Delete(&types.File{}).Error; err != nil {
		return utils.DetermineSQLError(err, "hard delete trashed file")
	}
	return nil
}

func (r *trashRepository) SoftDeleteFolderCascade(ctx context.Context, userID, folderID int, ts time.Time) error {
	const sql = `
WITH RECURSIVE cte AS (
    SELECT id FROM folders WHERE user_id = @user AND id = @root
	UNION ALL
    SELECT f.id FROM folders f JOIN cte ON f.parent_id = cte.id WHERE f.user_id = @user
)
UPDATE folders
	SET deleted_at = @ts
	WHERE id IN (SELECT id FROM cte);

UPDATE files
	SET deleted_at = @ts
	WHERE user_id = @user
	AND folder_id IN (SELECT id FROM cte);
`
	if err := r.db.WithContext(ctx).
		Raw(sql, map[string]interface{}{"user": userID, "root": folderID, "ts": ts}).
		Error; err != nil {
		return utils.DetermineSQLError(err, "soft delete folder cascade")
	}
	return nil
}

func (r *trashRepository) RestoreFolderCascade(ctx context.Context, userID, folderID int) error {
	const sql = `
WITH RECURSIVE cte AS (
    SELECT id FROM folders WHERE user_id = @user AND id = @root
	UNION ALL
    SELECT f.id FROM folders f JOIN cte ON f.parent_id = cte.id WHERE f.user_id = @user
)
UPDATE folders
	SET deleted_at = NULL
	WHERE id IN (SELECT id FROM cte);

UPDATE files
	SET deleted_at = NULL
	WHERE user_id = @user
	AND folder_id IN (SELECT id FROM cte);
`
	if err := r.db.WithContext(ctx).
		Raw(sql, map[string]interface{}{"user": userID, "root": folderID}).
		Error; err != nil {
		return utils.DetermineSQLError(err, "restore folder cascade")
	}
	return nil
}

func (r *trashRepository) ListTrashedFolders(ctx context.Context, userID int) ([]*types.Folder, error) {
	var out []*types.Folder
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND deleted_at IS NOT NULL", userID).
		Find(&out).Error
	if err != nil {
		return nil, utils.DetermineSQLError(err, "list trashed folders")
	}
	return out, nil
}

func (r *trashRepository) ListFoldersToPurge(ctx context.Context, before time.Time) ([]*types.Folder, error) {
	var out []*types.Folder
	err := r.db.WithContext(ctx).
		Where("deleted_at IS NOT NULL AND deleted_at < ?", before).
		Find(&out).Error
	if err != nil {
		return nil, utils.DetermineSQLError(err, "list folders to purge")
	}
	return out, nil
}

func (r *trashRepository) HardDeleteFolderByID(ctx context.Context, folderID int) error {
	if err := r.db.WithContext(ctx).
		Unscoped().
		Where("id = ?", folderID).
		Delete(&types.Folder{}).Error; err != nil {
		return utils.DetermineSQLError(err, "hard delete folder")
	}
	return nil
}

func (r *trashRepository) PermanentDeleteFolder(ctx context.Context, userID, folderID int) ([]int, error) {
	const sqlGetFiles = `
WITH RECURSIVE cte AS (
    SELECT id FROM folders WHERE user_id = @user AND id = @root
	UNION ALL
    SELECT f.id FROM folders f JOIN cte ON f.parent_id = cte.id WHERE f.user_id = @user
)
SELECT fi.id
	FROM files fi
	WHERE fi.user_id = @user AND fi.folder_id IN (SELECT id FROM cte);
`
	rows, err := r.db.WithContext(ctx).
		Raw(sqlGetFiles, map[string]interface{}{"user": userID, "root": folderID}).
		Rows()
	if err != nil {
		return nil, utils.DetermineSQLError(err, "list files for permanent delete")
	}
	defer rows.Close()

	var filesID []int
	for rows.Next() {
		var id int
		rows.Scan(&id)
		filesID = append(filesID, id)
	}

	const sqlDel = `
WITH RECURSIVE cte AS (
    SELECT id FROM folders WHERE user_id = @user AND id = @root
	UNION ALL
    SELECT f.id FROM folders f JOIN cte ON f.parent_id = cte.id WHERE f.user_id = @user
)
DELETE FROM files USING cte WHERE files.folder_id = cte.id;
DELETE FROM folders USING cte WHERE folders.id = cte.id;
`
	if err := r.db.WithContext(ctx).
		Raw(sqlDel, map[string]interface{}{"user": userID, "root": folderID}).
		Error; err != nil {
		return filesID, utils.DetermineSQLError(err, "hard delete folder cascade")
	}

	return filesID, nil
}
