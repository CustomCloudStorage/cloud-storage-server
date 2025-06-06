package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/CustomCloudStorage/middleware"
	"github.com/CustomCloudStorage/repositories"
	"github.com/CustomCloudStorage/services"
	"github.com/CustomCloudStorage/types"
	"github.com/CustomCloudStorage/utils"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
)

type FolderHandler interface {
	HandleCreateFolder(w http.ResponseWriter, r *http.Request) error
	HandleGetFolder(w http.ResponseWriter, r *http.Request) error
	HandleUpdateFolder(w http.ResponseWriter, r *http.Request) error
	HandleListFolders(w http.ResponseWriter, r *http.Request) error
	DownloadFolderHandler(w http.ResponseWriter, r *http.Request) error
}

type folderHandler struct {
	folderRepository repositories.FolderRepository
	folderService    services.FolderService
}

func NewFolderHandler(folderRepository repositories.FolderRepository, folderService services.FolderService) FolderHandler {
	return &folderHandler{
		folderRepository: folderRepository,
		folderService:    folderService,
	}
}

func (h *folderHandler) HandleCreateFolder(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	claims := ctx.Value("claims").(jwt.MapClaims)
	userID, ok := claims["userID"].(float64)
	if !ok {
		return utils.ErrUnauthorized.New("invalid userID")
	}

	var folder types.Folder
	if err := json.NewDecoder(r.Body).Decode(&folder); err != nil {
		return utils.ErrBadRequest.Wrap(err, "decode folder JSON")
	}

	folder.UserID = int(userID)

	if err := h.folderRepository.Create(ctx, &folder); err != nil {
		return err
	}

	return middleware.WriteJSONResponse(w, http.StatusCreated, map[string]interface{}{
		"folder_id": folder.ID,
		"message":   "folder created successfully",
	})
}

func (h *folderHandler) HandleGetFolder(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	params := mux.Vars(r)

	claims := ctx.Value("claims").(jwt.MapClaims)
	userID, ok := claims["userID"].(float64)
	if !ok {
		return utils.ErrUnauthorized.New("invalid userID")
	}
	folderID, err := strconv.Atoi(params["folderID"])
	if err != nil {
		return utils.ErrBadRequest.Wrap(err, "invalid folder ID")
	}

	folder, err := h.folderRepository.GetByID(ctx, folderID, int(userID))
	if err != nil {
		return err
	}

	return middleware.WriteJSONResponse(w, http.StatusOK, map[string]interface{}{
		"folder": folder,
	})
}

func (h *folderHandler) HandleUpdateFolder(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	params := mux.Vars(r)

	claims := ctx.Value("claims").(jwt.MapClaims)
	userID, ok := claims["userID"].(float64)
	if !ok {
		return utils.ErrUnauthorized.New("invalid userID")
	}
	folderID, err := strconv.Atoi(params["folderID"])
	if err != nil {
		return utils.ErrBadRequest.Wrap(err, "invalid folder ID")
	}

	var req types.Folder
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return utils.ErrBadRequest.Wrap(err, "decode folder JSON")
	}

	current, err := h.folderRepository.GetByID(ctx, folderID, int(userID))
	if err != nil {
		return err
	}
	if !req.UpdatedAt.Equal(current.UpdatedAt) {
		return utils.ErrConflict.New("the folder was modified by another process")
	}

	current.Name = req.Name
	current.ParentID = req.ParentID

	if err := h.folderRepository.Update(ctx, current); err != nil {
		return err
	}

	return middleware.WriteJSONResponse(w, http.StatusOK, map[string]interface{}{
		"message": "folder updated successfully",
	})
}

func (h *folderHandler) HandleListFolders(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	claims := ctx.Value("claims").(jwt.MapClaims)
	userID, ok := claims["userID"].(float64)
	if !ok {
		return utils.ErrUnauthorized.New("invalid userID")
	}

	folders, err := h.folderRepository.ListByUserID(ctx, int(userID))
	if err != nil {
		return err
	}

	return middleware.WriteJSONResponse(w, http.StatusOK, map[string]interface{}{
		"folders": folders,
	})
}

func (h *folderHandler) DownloadFolderHandler(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	params := mux.Vars(r)

	claims := ctx.Value("claims").(jwt.MapClaims)
	userID, ok := claims["userID"].(float64)
	if !ok {
		return utils.ErrUnauthorized.New("invalid userID")
	}
	folderID, err := strconv.Atoi(params["folderID"])
	if err != nil {
		return utils.ErrBadRequest.Wrap(err, "invalid folder ID")
	}

	reader, archiveName, err := h.folderService.DownloadFolder(ctx, int(userID), folderID)
	if err != nil {
		return err
	}
	defer reader.Close()

	w.Header().Set("Content-Type", "application/zip")
	w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, archiveName))

	if _, err := io.Copy(w, reader); err != nil {
		return utils.ErrInternal.Wrap(err, "stream zip archive")
	}
	return nil
}
