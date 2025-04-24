package services

import (
	"context"
	"io"

	"github.com/CustomCloudStorage/repositories"
	"github.com/CustomCloudStorage/types"
	"github.com/google/uuid"
)

type ServiceConfig struct {
	StorageDir string
	TmpUpload  string
}

type fileService struct {
	userRepository   repositories.UserRepository
	fileRepository   repositories.FileRepository
	folderRepository repositories.FolderRepository
	storageDir       string
}

type uploadService struct {
	userRepository          repositories.UserRepository
	fileRepository          repositories.FileRepository
	uploadSessionRepository repositories.UploadSessionRepository
	uploadPartRepository    repositories.UploadPartRepository
	storageDir              string
	tmpUpload               string
}

func NewFileService(userRepo repositories.UserRepository, fileRepo repositories.FileRepository, folderRepo repositories.FolderRepository, cfg ServiceConfig) *fileService {
	return &fileService{
		userRepository:   userRepo,
		fileRepository:   fileRepo,
		folderRepository: folderRepo,
		storageDir:       cfg.StorageDir,
	}
}

func NewUploadService(userRepo repositories.UserRepository, fileRepo repositories.FileRepository, uploadSessionRepo repositories.UploadSessionRepository, uploadPartRepo repositories.UploadPartRepository, cfg ServiceConfig) *uploadService {
	return &uploadService{
		userRepository:          userRepo,
		fileRepository:          fileRepo,
		uploadSessionRepository: uploadSessionRepo,
		uploadPartRepository:    uploadPartRepo,
		storageDir:              cfg.StorageDir,
		tmpUpload:               cfg.TmpUpload,
	}
}

type FileService interface {
	DeleteFile(ctx context.Context, id int, userID int) error
}

type UploadService interface {
	InitSession(ctx context.Context, session *types.UploadSession) error
	UploadPart(ctx context.Context, sessionID uuid.UUID, partNumber int, data io.Reader) error
	GetProgress(ctx context.Context, sessionID uuid.UUID) (int64, int, error)
	Complete(ctx context.Context, sessionID uuid.UUID) (*types.File, error)
	Abort(ctx context.Context, sessionID uuid.UUID) error
}
