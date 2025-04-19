package repositories

import (
	"context"

	"github.com/CustomCloudStorage/types"
	"gorm.io/gorm"
)

type user struct {
	db *gorm.DB
}

type folder struct {
	db *gorm.DB
}

type file struct {
	db *gorm.DB
}

type Repository struct {
	User   *user
	Folder *folder
	File   *file
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		User: &user{
			db: db,
		},
		Folder: &folder{
			db: db,
		},
		File: &file{
			db: db,
		},
	}
}

type UserRepository interface {
	GetByID(context.Context, int) (*types.User, error)
	Create(context.Context, *types.User) error
	UpdateProfile(context.Context, *types.Profile, int) error
	UpdateAccount(context.Context, *types.Account, int) error
	UpdateCredentials(context.Context, *types.Credentials, int) error
	Delete(context.Context, int) error
	List(context.Context) ([]types.User, error)
}

type FolderRepository interface {
	Create(ctx context.Context, file *types.Folder) error
	GetByID(ctx context.Context, id int, userID int) (*types.Folder, error)
	Update(ctx context.Context, folder *types.Folder) error
	Delete(ctx context.Context, id int, userID int) error
	ListByUserID(ctx context.Context, userID int) ([]types.Folder, error)
}

type FileRepository interface {
	Create(ctx context.Context, file *types.File) error
	GetByID(ctx context.Context, id int, userID int) (*types.File, error)
	Delete(ctx context.Context, id int, userID int) error
	ListByUserID(ctx context.Context, userID int) ([]types.File, error)
	UpdateName(ctx context.Context, id int, userID int, name string)
	UpdateFolder(ctx context.Context, id int, userID int, folderID int)
}
