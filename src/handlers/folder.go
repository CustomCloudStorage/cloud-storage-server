package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/CustomCloudStorage/types"
	"github.com/CustomCloudStorage/utils"
	"github.com/gorilla/mux"
)

func (h *folderHandler) HandleCreateFolder(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	log.Println("[POST] Creating folder")

	var folder types.Folder

	if err := json.NewDecoder(r.Body).Decode(&folder); err != nil {
		return err
	}

	if err := h.folderRepository.Create(ctx, &folder); err != nil {
		return err
	}

	return nil
}

func (h *folderHandler) HandleGetFolder(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	params := mux.Vars(r)
	userID, err := strconv.Atoi(params["id"])
	if err != nil {
		return err
	}
	folderID, err := strconv.Atoi(params["folderID"])
	if err != nil {
		return err
	}

	log.Println("[GET] Fetching folder with id:", folderID, " by user:", userID)

	folder, err := h.folderRepository.GetByID(ctx, folderID, userID)
	if err != nil {
		return err
	}

	if err := json.NewEncoder(w).Encode(folder); err != nil {
		return err
	}

	return nil
}

func (h *folderHandler) HandleUpdateFolder(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	params := mux.Vars(r)
	userID, err := strconv.Atoi(params["id"])
	if err != nil {
		return nil
	}
	folderID, err := strconv.Atoi(params["folderID"])
	if err != nil {
		return err
	}

	log.Println("[PUT] Updating folder with id:", folderID, " by user:", userID)

	var req types.Folder

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return err
	}

	folder, err := h.folderRepository.GetByID(ctx, folderID, userID)
	if err != nil {
		return err
	}

	if req.UpdatedAt != folder.UpdatedAt {
		return utils.ErrDataConflict.New("The folder was changed by another user")
	}

	if err := h.folderRepository.Update(ctx, folder); err != nil {
		return err
	}

	return nil
}

func (h *folderHandler) HandleDeleteFolder(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	params := mux.Vars(r)
	userID, err := strconv.Atoi(params["id"])
	if err != nil {
		return err
	}
	folderID, err := strconv.Atoi(params["folderID"])
	if err != nil {
		return err
	}

	log.Println("[DELETE] Deleting folder with id:", folderID, " by user:", userID)

	if err := h.folderRepository.Delete(ctx, folderID, userID); err != nil {
		return err
	}

	return nil
}

func (h *folderHandler) HandleListFolders(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	params := mux.Vars(r)
	userID, err := strconv.Atoi(params["id"])
	if err != nil {
		return err
	}

	log.Println("[GET] Fetching all folders by user:", userID)

	folders, err := h.folderRepository.ListByUserID(ctx, userID)
	if err != nil {
		return err
	}

	if err := json.NewEncoder(w).Encode(folders); err != nil {
		return err
	}

	return nil
}
