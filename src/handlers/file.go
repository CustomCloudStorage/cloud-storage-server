package handlers

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func (h *fileHandler) HandleGetFile(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	params := mux.Vars(r)
	userID, err := strconv.Atoi(params["id"])
	if err != nil {
		return err
	}
	fileID, err := strconv.Atoi(params["fileID"])
	if err != nil {
		return err
	}

	log.Println("[GET] Fetching file with id:", fileID, " by user:", userID)

	file, err := h.fileRepository.GetByID(ctx, fileID, userID)
	if err != nil {
		return err
	}

	if err := json.NewEncoder(w).Encode(file); err != nil {
		return err
	}

	return nil
}

func (h *fileHandler) HandleDeleteFile(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	params := mux.Vars(r)
	userID, err := strconv.Atoi(params["id"])
	if err != nil {
		return err
	}
	fileID, err := strconv.Atoi(params["fileID"])
	if err != nil {
		return err
	}

	log.Println("[DELETE] Deleting file with id:", fileID, " by user:", userID)

	if err := h.fileService.DeleteFile(ctx, fileID, userID); err != nil {
		return err
	}

	if err := h.fileRepository.Delete(ctx, fileID, userID); err != nil {
		return err
	}

	return nil
}

func (h *fileHandler) HandleListFiles(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	params := mux.Vars(r)
	userID, err := strconv.Atoi(params["id"])
	if err != nil {
		return err
	}

	log.Println("[GET] Fetching all files by user:", userID)

	files, err := h.fileRepository.ListByUserID(ctx, userID)
	if err != nil {
		return err
	}

	if err := json.NewEncoder(w).Encode(files); err != nil {
		return err
	}

	return nil
}

func (h *fileHandler) HandleUpdateName(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	params := mux.Vars(r)
	userID, err := strconv.Atoi(params["id"])
	if err != nil {
		return err
	}
	fileID, err := strconv.Atoi(params["fileID"])
	if err != nil {
		return err
	}

	var patchName string
	if err := json.NewDecoder(r.Body).Decode(&patchName); err != nil {
		return err
	}

	if err := h.fileRepository.UpdateName(ctx, fileID, userID, patchName); err != nil {
		return err
	}

	return nil
}

func (h *fileHandler) HandleUpdateFolderID(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	params := mux.Vars(r)
	userID, err := strconv.Atoi(params["id"])
	if err != nil {
		return err
	}
	fileID, err := strconv.Atoi(params["fileID"])
	if err != nil {
		return err
	}

	var patchFolderID int
	if err := json.NewDecoder(r.Body).Decode(&patchFolderID); err != nil {
		return err
	}

	if err := h.fileRepository.UpdateFolder(ctx, fileID, userID, patchFolderID); err != nil {
		return err
	}

	return nil
}

func (h *fileHandler) DownloadURLHandler(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	userID, _ := strconv.Atoi(vars["userID"])
	fileID, _ := strconv.Atoi(vars["fileID"])

	url, err := h.fileService.GenerateDownloadURL(r.Context(), userID, fileID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return err
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"download_url":"` + url + `"}`))

	return nil
}

func (h *fileHandler) DownloadByTokenHandler(w http.ResponseWriter, r *http.Request) error {
	token := r.URL.Query().Get("token")
	userID, fileID, err := h.fileService.ValidateDownloadToken(token)
	if err != nil {
		return err
	}

	dfile, err := h.fileService.DownloadFile(r.Context(), userID, fileID)
	if err != nil {
		return err
	}
	defer dfile.Reader.(io.Closer).Close()

	w.Header().Set("Content-Type", dfile.ContentType)
	w.Header().Set("Content-Disposition", "attachment; filename=\""+dfile.FileName+"\"")
	w.Header().Set("Content-Length", strconv.FormatInt(dfile.FileSize, 10))

	return nil
}
