package repositories

import (
	"context"

	"github.com/CustomCloudStorage/types"
	"github.com/google/uuid"
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

type uploadSession struct {
	db *gorm.DB
}

type uploadPart struct {
	db *gorm.DB
}

type Repository struct {
	User          *user
	Folder        *folder
	File          *file
	UploadSession *uploadSession
	UploadPart    *uploadPart
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
		UploadSession: &uploadSession{
			db: db,
		},
		UploadPart: &uploadPart{
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
	ReserveStorage(ctx context.Context, userID int, size int64) error
	ReleaseStorage(ctx context.Context, userID int, size int64) error
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

type UploadSessionRepository interface {
	Create(ctx context.Context, session *types.UploadSession) error
	GetByID(ctx context.Context, id uuid.UUID) (*types.UploadSession, error)
	Delete(ctx context.Context, id uuid.UUID) error
	ListOlderThan(ctx context.Context, olderThanMinutes int) ([]types.UploadSession, error)
}

type UploadPartRepository interface {
	Create(ctx context.Context, part *types.UploadPart) error
	ListBySession(ctx context.Context, sessionID uuid.UUID) ([]types.UploadPart, error)
	DeleteBySession(ctx context.Context, sessionID uuid.UUID) error
}
