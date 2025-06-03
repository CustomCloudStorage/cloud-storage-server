package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/CustomCloudStorage/middleware"
	"github.com/CustomCloudStorage/repositories"
	"github.com/CustomCloudStorage/services"
	"github.com/CustomCloudStorage/utils"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
)

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

type trashHandler struct {
	trashRepository repositories.TrashRepository
	trashService    services.TrashService
}

func NewTrashHandler(trashRepo repositories.TrashRepository, trashService services.TrashService) TrashHandler {
	return &trashHandler{
		trashRepository: trashRepo,
		trashService:    trashService,
	}
}

func (h *trashHandler) ListFilesHandler(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	claims := ctx.Value("claims").(jwt.MapClaims)
	userID, ok := claims["userID"].(float64)
	if !ok {
		return utils.ErrUnauthorized.New("invalid userID")
	}

	files, err := h.trashRepository.ListTrashedFiles(ctx, int(userID))
	if err != nil {
		return err
	}

	return middleware.WriteJSONResponse(w, http.StatusOK, map[string]interface{}{
		"files": files,
	})
}

func (h *trashHandler) DeleteFileHandler(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	claims := ctx.Value("claims").(jwt.MapClaims)
	userID, ok := claims["userID"].(float64)
	if !ok {
		return utils.ErrUnauthorized.New("invalid userID")
	}

	fileID, err := strconv.Atoi(mux.Vars(r)["fileID"])
	if err != nil {
		return utils.ErrBadRequest.Wrap(err, "invalid file ID")
	}

	if err := h.trashRepository.SoftDeleteFile(ctx, int(userID), fileID, time.Now()); err != nil {
		return err
	}

	return middleware.WriteJSONResponse(w, http.StatusOK, map[string]interface{}{
		"message": "file moved to trash",
	})
}

func (h *trashHandler) RestoreFileHandler(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	claims := ctx.Value("claims").(jwt.MapClaims)
	userID, ok := claims["userID"].(float64)
	if !ok {
		return utils.ErrUnauthorized.New("invalid userID")
	}

	fileID, err := strconv.Atoi(mux.Vars(r)["fileID"])
	if err != nil {
		return utils.ErrBadRequest.Wrap(err, "invalid file ID")
	}

	if err := h.trashRepository.RestoreFile(ctx, int(userID), fileID); err != nil {
		return err
	}

	return middleware.WriteJSONResponse(w, http.StatusOK, map[string]interface{}{
		"message": "file restored",
	})
}

func (h *trashHandler) PermanentDeleteFileHandler(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	claims := ctx.Value("claims").(jwt.MapClaims)
	userID, ok := claims["userID"].(float64)
	if !ok {
		return utils.ErrUnauthorized.New("invalid userID")
	}

	fileID, err := strconv.Atoi(mux.Vars(r)["fileID"])
	if err != nil {
		return utils.ErrBadRequest.Wrap(err, "invalid file ID")
	}

	if err := h.trashService.PermanentDeleteFile(ctx, int(userID), fileID); err != nil {
		return err
	}

	return middleware.WriteJSONResponse(w, http.StatusOK, map[string]interface{}{
		"message": "file permanently deleted",
	})
}

func (h *trashHandler) ListFoldersHandler(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	claims := ctx.Value("claims").(jwt.MapClaims)
	userID, ok := claims["userID"].(float64)
	if !ok {
		return utils.ErrUnauthorized.New("invalid userID")
	}

	folders, err := h.trashRepository.ListTrashedFolders(ctx, int(userID))
	if err != nil {
		return err
	}

	return middleware.WriteJSONResponse(w, http.StatusOK, map[string]interface{}{
		"folders": folders,
	})
}

func (h *trashHandler) DeleteFolderHandler(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	claims := ctx.Value("claims").(jwt.MapClaims)
	userID, ok := claims["userID"].(float64)
	if !ok {
		return utils.ErrUnauthorized.New("invalid userID")
	}

	folderID, err := strconv.Atoi(mux.Vars(r)["folderID"])
	if err != nil {
		return utils.ErrBadRequest.Wrap(err, "invalid folder ID")
	}

	if err := h.trashRepository.SoftDeleteFolderCascade(ctx, int(userID), folderID, time.Now()); err != nil {
		return err
	}

	return middleware.WriteJSONResponse(w, http.StatusOK, map[string]interface{}{
		"message": "folder moved to trash",
	})
}

func (h *trashHandler) RestoreFolderHandler(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	claims := ctx.Value("claims").(jwt.MapClaims)
	userID, ok := claims["userID"].(float64)
	if !ok {
		return utils.ErrUnauthorized.New("invalid userID")
	}
	folderID, err := strconv.Atoi(mux.Vars(r)["folderID"])
	if err != nil {
		return utils.ErrBadRequest.Wrap(err, "invalid folder ID")
	}

	if err := h.trashRepository.RestoreFolderCascade(ctx, int(userID), folderID); err != nil {
		return err
	}

	return middleware.WriteJSONResponse(w, http.StatusOK, map[string]interface{}{
		"message": "folder restored",
	})
}

func (h *trashHandler) PermanentDeleteFolderHandler(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	claims := ctx.Value("claims").(jwt.MapClaims)
	userID, ok := claims["userID"].(float64)
	if !ok {
		return utils.ErrUnauthorized.New("invalid userID")
	}

	folderID, err := strconv.Atoi(mux.Vars(r)["folderID"])
	if err != nil {
		return utils.ErrBadRequest.Wrap(err, "invalid folder ID")
	}

	if err := h.trashService.PermanentDeleteFolder(ctx, int(userID), folderID); err != nil {
		return err
	}

	return middleware.WriteJSONResponse(w, http.StatusOK, map[string]interface{}{
		"message": "folder permanently deleted",
	})
}
