package repositories

import (
	"context"

	"github.com/CustomCloudStorage/types"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

type fileRepository struct {
	db *gorm.DB
}

type folderRepository struct {
	db *gorm.DB
}

type uploadSessionRepository struct {
	db *gorm.DB
}

type uploadPartRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *userRepository {
	return &userRepository{
		db: db,
	}
}

func NewFileRepository(db *gorm.DB) *fileRepository {
	return &fileRepository{
		db: db,
	}
}

func NewFolderRepository(db *gorm.DB) *folderRepository {
	return &folderRepository{
		db: db,
	}
}

func NewUploadSessionRepository(db *gorm.DB) *uploadSessionRepository {
	return &uploadSessionRepository{
		db: db,
	}
}

func NewUploadPartRepository(db *gorm.DB) *uploadPartRepository {
	return &uploadPartRepository{
		db: db,
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
	UpdateUsedStorage(ctx context.Context, id int, newUsedStorage int64) error
	ReserveStorage(ctx context.Context, userID int, size int64) error
	ReleaseStorage(ctx context.Context, userID int, size int64) error
}

type FileRepository interface {
	Create(ctx context.Context, file *types.File) error
	GetByID(ctx context.Context, id int, userID int) (*types.File, error)
	Delete(ctx context.Context, id int, userID int) error
	ListByUserID(ctx context.Context, userID int) ([]types.File, error)
	UpdateName(ctx context.Context, id int, userID int, name string) error
	UpdateFolder(ctx context.Context, id int, userID int, folderID int) error
	ListFilesRecursive(ctx context.Context, userID, folderID int) ([]*types.FileWithPath, error)
}

type FolderRepository interface {
	Create(ctx context.Context, file *types.Folder) error
	GetByID(ctx context.Context, id int, userID int) (*types.Folder, error)
	Update(ctx context.Context, folder *types.Folder) error
	Delete(ctx context.Context, id int, userID int) error
	ListByUserID(ctx context.Context, userID int) ([]types.Folder, error)
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
