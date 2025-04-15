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

	log.Println("[POST] Creating file")

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
