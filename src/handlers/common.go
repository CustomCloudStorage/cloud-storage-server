package handlers

import (
	"net/http"

	"github.com/CustomCloudStorage/repositories"
	"github.com/CustomCloudStorage/services"
)

type userHandler struct {
	userRepository repositories.UserRepository
	fileRepository repositories.FileRepository
	fileService    services.FileService
}

type fileHandler struct {
	fileRepository repositories.FileRepository
	fileService    services.FileService
}

type folderHandler struct {
	folderRepository repositories.FolderRepository
	folderService    services.FolderService
}

type uploadHandler struct {
	uploadService services.UploadService
}

type trashHandler struct {
	trashRepository repositories.TrashRepository
	trashService    services.TrashService
}

func NewUserHandler(userRepository repositories.UserRepository, fileRepository repositories.FileRepository, fileService services.FileService) *userHandler {
	return &userHandler{
		userRepository: userRepository,
		fileRepository: fileRepository,
		fileService:    fileService,
	}
}

func NewFileHandler(fileRepository repositories.FileRepository, fileService services.FileService) *fileHandler {
	return &fileHandler{
		fileRepository: fileRepository,
		fileService:    fileService,
	}
}

func NewFolderHandler(folderRepository repositories.FolderRepository, folderService services.FolderService) *folderHandler {
	return &folderHandler{
		folderRepository: folderRepository,
		folderService:    folderService,
	}
}

func NewUploadHandler(uploadService services.UploadService) *uploadHandler {
	return &uploadHandler{
		uploadService: uploadService,
	}
}

func NewTrashHandler(trashRepo repositories.TrashRepository, trashService services.TrashService) *trashHandler {
	return &trashHandler{
		trashRepository: trashRepo,
		trashService:    trashService,
	}
}

type UserHandler interface {
	HandleGetUser(w http.ResponseWriter, r *http.Request) error
	HandleListUsers(w http.ResponseWriter, r *http.Request) error
	HandleCreateUser(w http.ResponseWriter, r *http.Request) error
	HandleUpdateProfile(w http.ResponseWriter, r *http.Request) error
	HandleUpdateAccount(w http.ResponseWriter, r *http.Request) error
	HandleUpdateCredentials(w http.ResponseWriter, r *http.Request) error
	HandleDeleteUser(w http.ResponseWriter, r *http.Request) error
}

type FileHandler interface {
	HandleGetFile(w http.ResponseWriter, r *http.Request) error
	HandleDeleteFile(w http.ResponseWriter, r *http.Request) error
	HandleListFiles(w http.ResponseWriter, r *http.Request) error
	HandleUpdateName(w http.ResponseWriter, r *http.Request) error
	HandleUpdateFolderID(w http.ResponseWriter, r *http.Request) error
	DownloadURLHandler(w http.ResponseWriter, r *http.Request) error
	DownloadByTokenHandler(w http.ResponseWriter, r *http.Request) error
	StreamFileHandler(w http.ResponseWriter, r *http.Request) error
	PreviewFileHandler(w http.ResponseWriter, r *http.Request) error
}

type FolderHandler interface {
	HandleCreateFolder(w http.ResponseWriter, r *http.Request) error
	HandleGetFolder(w http.ResponseWriter, r *http.Request) error
	HandleUpdateFolder(w http.ResponseWriter, r *http.Request) error
	HandleDeleteFolder(w http.ResponseWriter, r *http.Request) error
	HandleListFolders(w http.ResponseWriter, r *http.Request) error
	DownloadFolderHandler(w http.ResponseWriter, r *http.Request) error
}

type UploadHandler interface {
	InitSessionHandler(w http.ResponseWriter, r *http.Request) error
	UploadPartHandler(w http.ResponseWriter, r *http.Request) error
	ProgressHandler(w http.ResponseWriter, r *http.Request) error
	CompleteHandler(w http.ResponseWriter, r *http.Request) error
	AbortHandler(w http.ResponseWriter, r *http.Request) error
}

type TrashHandler interface {
	ListFilesHandler(w http.ResponseWriter, r *http.Request) error
	DeleteFileHandler(w http.ResponseWriter, r *http.Request) error
	RestoreFileHandler(w http.ResponseWriter, r *http.Request) error
	PermanentDeleteFileHandler(w http.ResponseWriter, r *http.Request) error

	ListFoldersHandler(w http.ResponseWriter, r *http.Request) error
	DeleteFolderHandler(w http.ResponseWriter, r *http.Request) error
	RestoreFolderHandler(w http.ResponseWriter, r *http.Request) error
	PermanentDeleteFolderHandler(w http.ResponseWriter, r *http.Request) error
}
