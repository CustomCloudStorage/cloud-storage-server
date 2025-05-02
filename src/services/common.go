package services

import (
	"context"
	"io"
	"time"

	"github.com/CustomCloudStorage/repositories"
	"github.com/CustomCloudStorage/types"
	"github.com/google/uuid"
)

type ServiceConfig struct {
	StorageDir    string
	TmpUpload     string
	Secret        string
	UrlTtlSeconds time.Duration
	Host          string
}

type TrashConfig struct {
	TTL             time.Duration
	CleanupInterval time.Duration
}

type fileService struct {
	userRepository   repositories.UserRepository
	fileRepository   repositories.FileRepository
	folderRepository repositories.FolderRepository
	storageDir       string
	secret           string
	urlTtlSeconds    time.Duration
	host             string
}

type folderService struct {
	folderRepository repositories.FolderRepository
	fileRepository   repositories.FileRepository
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

type trashService struct {
	trashRepository repositories.TrashRepository
	storageDir      string
	cfg             TrashConfig
}

func NewFileService(userRepo repositories.UserRepository, fileRepo repositories.FileRepository, folderRepo repositories.FolderRepository, cfg ServiceConfig) *fileService {
	return &fileService{
		userRepository:   userRepo,
		fileRepository:   fileRepo,
		folderRepository: folderRepo,
		storageDir:       cfg.StorageDir,
		secret:           cfg.Secret,
		urlTtlSeconds:    cfg.UrlTtlSeconds,
		host:             cfg.Host,
	}
}

func NewFolderService(fileRepo repositories.FileRepository, folderRepo repositories.FolderRepository, cfg ServiceConfig) *folderService {
	return &folderService{
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

func NewTrashService(trashRepo repositories.TrashRepository, cfg ServiceConfig, trashCfg TrashConfig) *trashService {
	svc := &trashService{
		trashRepository: trashRepo,
		storageDir:      cfg.StorageDir,
		cfg:             trashCfg,
	}
	go svc.purgeLoop()
	return svc
}

type FileService interface {
	GenerateDownloadURL(ctx context.Context, userID, fileID int) (string, error)
	ValidateDownloadToken(token string) (userID, fileID int, err error)
	DownloadFile(ctx context.Context, userID int, fileID int) (*types.DownloadedFile, error)
	DeleteFile(ctx context.Context, id int, userID int) error
	StreamFile(ctx context.Context, userID, fileID int) (*types.DownloadedFile, error)
	PreviewFile(ctx context.Context, userID, fileID int, w io.Writer) (time.Time, error)
}

type FolderService interface {
	DownloadFolder(ctx context.Context, userID, folderID int) (io.ReadCloser, string, error)
}

type UploadService interface {
	InitSession(ctx context.Context, session *types.UploadSession) error
	UploadPart(ctx context.Context, sessionID uuid.UUID, partNumber int, data io.Reader) error
	GetProgress(ctx context.Context, sessionID uuid.UUID) (int64, int, error)
	Complete(ctx context.Context, sessionID uuid.UUID) (*types.File, error)
	Abort(ctx context.Context, sessionID uuid.UUID) error
}

type TrashService interface {
	purgeLoop()
	PermanentDeleteFile(ctx context.Context, userID, fileID int) error
	PermanentDeleteFolder(ctx context.Context, userID, folderID int) error
}
