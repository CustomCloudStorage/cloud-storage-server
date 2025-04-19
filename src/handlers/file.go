package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func (h *Handler) HandleUploadFile(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	log.Println("[POST] Upload file")

	params := mux.Vars(r)
	userId, err := strconv.Atoi(params["id"])
	if err != nil {
		return err
	}

	folderIDStr := r.FormValue("folder_id")
	var folderID *int
	if folderIDStr != "" {
		id, err := strconv.Atoi(folderIDStr)
		if err != nil {
			folderID = &id
		}
	}

	fileData, header, err := r.FormFile("file")
	if err != nil {
		return err
	}
	defer fileData.Close()

	fileSize := header.Size
	fileName := header.Filename

	if err := h.Service.File.UploadFile(ctx, userId, folderID, fileName, fileSize, fileData); err != nil {
		return err
	}

	return nil
}

func (h *Handler) HandleDownloadFile(w http.ResponseWriter, r *http.Request) error {
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

	log.Println("[GET] Downloading file with id:", fileID)

	downloadedFile, err := h.Service.File.DownloadFile(ctx, userID, fileID)
	if err != nil {
		return err
	}

	defer func() {
		if closer, ok := downloadedFile.Reader.(interface{ Close() error }); ok {
			closer.Close()
		}
	}()

	w.Header().Set("Content-Disposition", "attachment; filename=\""+downloadedFile.FileName+"\"")
	w.Header().Set("Content-Type", downloadedFile.ContentType)
	w.Header().Set("Content-Length", strconv.FormatInt(downloadedFile.FileSize, 10))

	http.ServeContent(w, r, downloadedFile.FileName, downloadedFile.ModTime, downloadedFile.Reader)

	return nil
}

func (h *Handler) HandleGetFile(w http.ResponseWriter, r *http.Request) error {
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

	file, err := h.Repository.File.GetByID(ctx, fileID, userID)
	if err != nil {
		return err
	}

	if err := json.NewEncoder(w).Encode(file); err != nil {
		return err
	}

	return nil
}

func (h *Handler) HandleDeleteFile(w http.ResponseWriter, r *http.Request) error {
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

	if err := h.Service.File.DeleteFile(ctx, fileID, userID); err != nil {
		return err
	}

	if err := h.Repository.File.Delete(ctx, fileID, userID); err != nil {
		return err
	}

	return nil
}

func (h *Handler) HandleListFiles(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	params := mux.Vars(r)
	userID, err := strconv.Atoi(params["id"])
	if err != nil {
		return err
	}

	log.Println("[GET] Fetching all files by user:", userID)

	files, err := h.Repository.File.ListByUserID(ctx, userID)
	if err != nil {
		return err
	}

	if err := json.NewEncoder(w).Encode(files); err != nil {
		return err
	}

	return nil
}

func (h *Handler) HandleUpdateName(w http.ResponseWriter, r *http.Request) error {
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

	if err := h.Repository.File.UpdateName(ctx, fileID, userID, patchName); err != nil {
		return err
	}

	return nil
}

func (h *Handler) HandleUpdateFolderID(w http.ResponseWriter, r *http.Request) error {
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

	if err := h.Repository.File.UpdateFolder(ctx, fileID, userID, patchFolderID); err != nil {
		return err
	}

	return nil
}
