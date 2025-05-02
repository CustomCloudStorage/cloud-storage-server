package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/CustomCloudStorage/utils"
	"github.com/gorilla/mux"
)

func (h *trashHandler) ListFilesHandler(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	//нужно достать userID из токена
	files, err := h.trashRepository.ListTrashedFiles(ctx, userID)
	if err != nil {
		return err
	}

	if err := json.NewEncoder(w).Encode(files); err != nil {
		return utils.ErrJsonEncode.Wrap(err, "failed to encode trashed files %s to JSON")
	}

	return nil
}

func (h *trashHandler) DeleteFileHandler(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	params := mux.Vars(r)
	fileID, err := strconv.Atoi(params["fileID"])
	if err != nil {
		return err
	}
	//нужно достать userID из токена

	if err := h.trashRepository.SoftDeleteFile(ctx, userID, fileID, time.Now()); err != nil {
		return err
	}

	return nil
}

func (h *trashHandler) RestoreFileHandler(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	params := mux.Vars(r)
	fileID, err := strconv.Atoi(params["fileID"])
	if err != nil {
		return err
	}
	//нужно достать userID из токена

	if err := h.trashRepository.RestoreFile(ctx, userID, fileID); err != nil {
		return err
	}

	return err
}

func (h *trashHandler) PermanentDeleteFileHandler(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	params := mux.Vars(r)
	fileID, err := strconv.Atoi(params["fileID"])
	if err != nil {
		return err
	}
	//нужно достать userID из токена

	if err := h.trashService.PermanentDeleteFile(ctx, userID, fileID); err != nil {
		return err
	}

	return err
}

func (h *trashHandler) ListFoldersHandler(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	//нужно достать userID из токена
	folders, err := h.trashRepository.ListTrashedFolders(ctx, userID)
	if err != nil {
		return err
	}

	if err := json.NewEncoder(w).Encode(folders); err != nil {
		return utils.ErrJsonEncode.Wrap(err, "failed to encode trashed files %s to JSON")
	}

	return nil
}

func (h *trashHandler) DeleteFolderHandler(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	params := mux.Vars(r)
	folderID, err := strconv.Atoi(params["folderID"])
	if err != nil {
		return err
	}
	//нужно достать userID из токена

	if err := h.trashRepository.SoftDeleteFolderCascade(ctx, userID, folderID, time.Now()); err != nil {
		return err
	}

	return nil
}

func (h *trashHandler) RestoreFolderHandler(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	params := mux.Vars(r)
	folderID, err := strconv.Atoi(params["folderID"])
	if err != nil {
		return err
	}
	//нужно достать userID из токена

	if err := h.trashRepository.RestoreFolderCascade(ctx, userID, folderID); err != nil {
		return err
	}

	return err
}

func (h *trashHandler) PermanentDeleteFolderHandler(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	params := mux.Vars(r)
	folderID, err := strconv.Atoi(params["folderID"])
	if err != nil {
		return err
	}
	//нужно достать userID из токена

	if err := h.trashService.PermanentDeleteFolder(ctx, userID, folderID); err != nil {
		return err
	}

	return err
}
