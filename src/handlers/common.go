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
}

type uploadHandler struct {
	uploadService services.UploadService
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

func NewFolderHandler(folderRepository repositories.FolderRepository) *folderHandler {
	return &folderHandler{
		folderRepository: folderRepository,
	}
}

func NewUploadHandler(uploadService services.UploadService) *uploadHandler {
	return &uploadHandler{
		uploadService: uploadService,
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
}

type FolderHandler interface {
	HandleCreateFolder(w http.ResponseWriter, r *http.Request) error
	HandleGetFolder(w http.ResponseWriter, r *http.Request) error
	HandleUpdateFolder(w http.ResponseWriter, r *http.Request) error
	HandleDeleteFolder(w http.ResponseWriter, r *http.Request) error
	HandleListFolders(w http.ResponseWriter, r *http.Request) error
}

type UploadHandler interface {
	InitSessionHandler(w http.ResponseWriter, r *http.Request) error
	UploadPartHandler(w http.ResponseWriter, r *http.Request) error
	ProgressHandler(w http.ResponseWriter, r *http.Request) error
	CompleteHandler(w http.ResponseWriter, r *http.Request) error
	AbortHandler(w http.ResponseWriter, r *http.Request) error
}

type HandlerWithErrorFunc func(w http.ResponseWriter, r *http.Request) error

func HandleError(handler HandlerWithErrorFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := handler(w, r)
		if err != nil {
			switch {
			case errorx.IsOfType(err, utils.ErrNotFound):
				log.Println(err.Error())
				writeJSONResponse(w, http.StatusNotFound, map[string]string{"error": "Data not found"})
			case errorx.IsOfType(err, utils.ErrAlreadyExist):
				log.Println(err.Error())
				writeJSONResponse(w, http.StatusBadRequest, map[string]string{"error": "The request body is invalid"})
			case errorx.IsOfType(err, utils.ErrDataConflict):
				log.Println(err.Error())
				writeJSONResponse(w, http.StatusConflict, map[string]string{
					"сonflict": "The data was changed",
				})
			default:
				log.Println("Internal server error:", err.Error())
				writeJSONResponse(w, http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
			}
		}
	}
}

func writeJSONResponse(w http.ResponseWriter, httpCode int, message map[string]string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpCode)
	json.NewEncoder(w).Encode(message)
}
