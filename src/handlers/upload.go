package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/CustomCloudStorage/types"
	"github.com/CustomCloudStorage/utils"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

func (h *Handler) InitSessionHandler(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	var session types.UploadSession
	if err := json.NewDecoder(r.Body).Decode(&session); err != nil {
		return err
	}

	if err := h.Service.MultiPart.InitSession(ctx, &session); err != nil {
		return err
	}

	return nil
}

func (h *Handler) UploadPartHandler(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	params := mux.Vars(r)
	sessionID, err := uuid.Parse(params["sessionID"])
	if err != nil {
		return err
	}
	partNum, err := strconv.Atoi(params["partNumber"])
	if err != nil {
		return err
	}

	if err := h.Service.MultiPart.UploadPart(ctx, sessionID, partNum, r.Body); err != nil {
		return err
	}

	return nil
}

func (h *Handler) ProgressHandler(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	sessionID, err := uuid.Parse(mux.Vars(r)["sessionID"])
	if err != nil {
		return err
	}

	uploaded, total, err := h.Service.MultiPart.GetProgress(ctx, sessionID)
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"uploaded_bytes": uploaded,
		"total_parts":    total,
	})

	return nil
}

func (h *Handler) CompleteHandler(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	sessionID, err := uuid.Parse(mux.Vars(r)["sessionID"])
	if err != nil {
		return err
	}
	fileMeta, err := h.Service.MultiPart.Complete(ctx, sessionID)
	if err != nil {
		return err
	}

	if err := json.NewEncoder(w).Encode(fileMeta); err != nil {
		return utils.ErrJsonEncode.Wrap(err, "failed to encode fileMeta to JSON")
	}

	return nil
}

func (h *Handler) AbortHandler(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	sessionID, err := uuid.Parse(mux.Vars(r)["sessionID"])
	if err != nil {
		return err
	}

	if err := h.Service.MultiPart.Abort(ctx, sessionID); err != nil {
		return err
	}

	return nil
}
