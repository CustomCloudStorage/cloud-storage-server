package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/CustomCloudStorage/utils"
	"github.com/gorilla/mux"
)

func (h *trashHandler) ListFilesHandler(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	// извлекаем userID, предположим, что middleware положил его в контекст под ключ "userID"
	userID, ok := ctx.Value("userID").(int)
	if !ok {
		return utils.ErrBadRequest.New("user not authenticated")
	}

	files, err := h.trashRepository.ListTrashedFiles(ctx, userID)
	if err != nil {
		return err
	}

	return writeJSONResponse(w, http.StatusOK, map[string]interface{}{
		"files": files,
	})
}

func (h *trashHandler) DeleteFileHandler(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	userID, ok := ctx.Value("userID").(int)
	if !ok {
		return utils.ErrBadRequest.New("user not authenticated")
	}

	fileID, err := strconv.Atoi(mux.Vars(r)["fileID"])
	if err != nil {
		return utils.ErrBadRequest.Wrap(err, "invalid file ID")
	}

	if err := h.trashRepository.SoftDeleteFile(ctx, userID, fileID, time.Now()); err != nil {
		return err
	}

	return writeJSONResponse(w, http.StatusOK, map[string]interface{}{
		"message": "file moved to trash",
	})
}

func (h *trashHandler) RestoreFileHandler(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	userID, ok := ctx.Value("userID").(int)
	if !ok {
		return utils.ErrBadRequest.New("user not authenticated")
	}

	fileID, err := strconv.Atoi(mux.Vars(r)["fileID"])
	if err != nil {
		return utils.ErrBadRequest.Wrap(err, "invalid file ID")
	}

	if err := h.trashRepository.RestoreFile(ctx, userID, fileID); err != nil {
		return err
	}

	return writeJSONResponse(w, http.StatusOK, map[string]interface{}{
		"message": "file restored",
	})
}

func (h *trashHandler) PermanentDeleteFileHandler(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	userID, ok := ctx.Value("userID").(int)
	if !ok {
		return utils.ErrBadRequest.New("user not authenticated")
	}

	fileID, err := strconv.Atoi(mux.Vars(r)["fileID"])
	if err != nil {
		return utils.ErrBadRequest.Wrap(err, "invalid file ID")
	}

	if err := h.trashService.PermanentDeleteFile(ctx, userID, fileID); err != nil {
		return err
	}

	return writeJSONResponse(w, http.StatusOK, map[string]interface{}{
		"message": "file permanently deleted",
	})
}

func (h *trashHandler) ListFoldersHandler(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	userID, ok := ctx.Value("userID").(int)
	if !ok {
		return utils.ErrBadRequest.New("user not authenticated")
	}

	folders, err := h.trashRepository.ListTrashedFolders(ctx, userID)
	if err != nil {
		return err
	}

	return writeJSONResponse(w, http.StatusOK, map[string]interface{}{
		"folders": folders,
	})
}

func (h *trashHandler) DeleteFolderHandler(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	userID, ok := ctx.Value("userID").(int)
	if !ok {
		return utils.ErrBadRequest.New("user not authenticated")
	}

	folderID, err := strconv.Atoi(mux.Vars(r)["folderID"])
	if err != nil {
		return utils.ErrBadRequest.Wrap(err, "invalid folder ID")
	}

	if err := h.trashRepository.SoftDeleteFolderCascade(ctx, userID, folderID, time.Now()); err != nil {
		return err
	}

	return writeJSONResponse(w, http.StatusOK, map[string]interface{}{
		"message": "folder moved to trash",
	})
}

func (h *trashHandler) RestoreFolderHandler(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	userID, ok := ctx.Value("userID").(int)
	if !ok {
		return utils.ErrBadRequest.New("user not authenticated")
	}

	folderID, err := strconv.Atoi(mux.Vars(r)["folderID"])
	if err != nil {
		return utils.ErrBadRequest.Wrap(err, "invalid folder ID")
	}

	if err := h.trashRepository.RestoreFolderCascade(ctx, userID, folderID); err != nil {
		return err
	}

	return writeJSONResponse(w, http.StatusOK, map[string]interface{}{
		"message": "folder restored",
	})
}

func (h *trashHandler) PermanentDeleteFolderHandler(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	userID, ok := ctx.Value("userID").(int)
	if !ok {
		return utils.ErrBadRequest.New("user not authenticated")
	}

	folderID, err := strconv.Atoi(mux.Vars(r)["folderID"])
	if err != nil {
		return utils.ErrBadRequest.Wrap(err, "invalid folder ID")
	}

	if err := h.trashService.PermanentDeleteFolder(ctx, userID, folderID); err != nil {
		return err
	}

	return writeJSONResponse(w, http.StatusOK, map[string]interface{}{
		"message": "folder permanently deleted",
	})
}
