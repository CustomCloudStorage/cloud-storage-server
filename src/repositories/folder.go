package repositories

import (
	"context"

	"github.com/CustomCloudStorage/types"
	"github.com/CustomCloudStorage/utils"
	"gorm.io/gorm"
)

type FolderRepository interface {
	Create(ctx context.Context, file *types.Folder) error
	GetByID(ctx context.Context, id int, userID int) (*types.Folder, error)
	Update(ctx context.Context, folder *types.Folder) error
	ListByUserID(ctx context.Context, userID int) ([]types.Folder, error)
}

type folderRepository struct {
	db *gorm.DB
}

func NewFolderRepository(db *gorm.DB) FolderRepository {
	return &folderRepository{
		db: db,
	}
}

func (r *folderRepository) Create(ctx context.Context, folder *types.Folder) error {
	if err := r.db.WithContext(ctx).Create(folder).Error; err != nil {
		return utils.DetermineSQLError(err, "create folder")
	}
	return nil
}

func (r *folderRepository) GetByID(ctx context.Context, id int, userID int) (*types.Folder, error) {
	var folder types.Folder
	if err := r.db.WithContext(ctx).
		Where("id = ? AND user_id = ?", id, userID).
		First(&folder).Error; err != nil {
		return nil, utils.DetermineSQLError(err, "get folder")
	}
	return &folder, nil
}

func (r *folderRepository) Update(ctx context.Context, folder *types.Folder) error {
	if err := r.db.Save(folder).Error; err != nil {
		return utils.DetermineSQLError(err, "update folder")
	}
	return nil
}

func (r *folderRepository) ListByUserID(ctx context.Context, userID int) ([]types.Folder, error) {
	var folders []types.Folder
	if err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Find(&folders).Error; err != nil {
		return nil, utils.DetermineSQLError(err, "list folders")
	}
	return folders, nil
}
