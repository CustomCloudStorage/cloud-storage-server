package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/CustomCloudStorage/repositories"
	"github.com/CustomCloudStorage/services"
	"github.com/CustomCloudStorage/utils"
	"github.com/joomcode/errorx"
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

type HandlerWithErrorFunc func(w http.ResponseWriter, r *http.Request) error

func HandleError(handler HandlerWithErrorFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := handler(w, r); err != nil {
			log.Println(err.Error())
			var (
				status  int
				message string
			)
			switch {
			case errorx.IsOfType(err, utils.ErrBadRequest):
				status = http.StatusBadRequest
				message = "Bad request"
			case errorx.IsOfType(err, utils.ErrNotFound):
				status = http.StatusNotFound
				message = "Resource not found"
			case errorx.IsOfType(err, utils.ErrConflict):
				status = http.StatusConflict
				message = "Conflict occurred"
			default:
				status = http.StatusInternalServerError
				message = "Internal server error"
			}

			writeJSONResponse(w, status, map[string]interface{}{
				"error": message,
			})
		}
	}
}

func writeJSONResponse(w http.ResponseWriter, httpCode int, payload interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpCode)
	return json.NewEncoder(w).Encode(payload)
}
