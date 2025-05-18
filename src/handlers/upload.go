package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/CustomCloudStorage/middleware"
	"github.com/CustomCloudStorage/types"
	"github.com/CustomCloudStorage/utils"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

func (h *uploadHandler) InitSessionHandler(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	var session types.UploadSession
	if err := json.NewDecoder(r.Body).Decode(&session); err != nil {
		return utils.ErrBadRequest.Wrap(err, "decode upload session payload")
	}

	if err := h.uploadService.InitSession(ctx, &session); err != nil {
		return err
	}

	middleware.WriteJSONResponse(w, http.StatusCreated, map[string]interface{}{
		"session_id": session.ID.String(),
		"message":    "upload session initialized",
	})
	return nil
}

func (h *uploadHandler) UploadPartHandler(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	params := mux.Vars(r)
	sessionID, err := uuid.Parse(params["sessionID"])
	if err != nil {
		return utils.ErrBadRequest.Wrap(err, "invalid session ID")
	}
	partNum, err := strconv.Atoi(params["partNumber"])
	if err != nil {
		return utils.ErrBadRequest.Wrap(err, "invalid part number")
	}

	if err := h.uploadService.UploadPart(ctx, sessionID, partNum, r.Body); err != nil {
		return err
	}

	middleware.WriteJSONResponse(w, http.StatusOK, map[string]interface{}{
		"part_number": partNum,
		"message":     "part uploaded successfully",
	})
	return nil
}

func (h *uploadHandler) ProgressHandler(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	sessionID, err := uuid.Parse(mux.Vars(r)["sessionID"])
	if err != nil {
		return utils.ErrBadRequest.Wrap(err, "invalid session ID")
	}

	uploaded, total, err := h.uploadService.GetProgress(ctx, sessionID)
	if err != nil {
		return err
	}

	middleware.WriteJSONResponse(w, http.StatusOK, map[string]interface{}{
		"uploaded": uploaded,
		"total":    total,
	})
	return nil
}

func (h *uploadHandler) CompleteHandler(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	sessionID, err := uuid.Parse(mux.Vars(r)["sessionID"])
	if err != nil {
		return utils.ErrBadRequest.Wrap(err, "invalid session ID")
	}
	fileMeta, err := h.uploadService.Complete(ctx, sessionID)
	if err != nil {
		return err
	}

	return middleware.WriteJSONResponse(w, http.StatusOK, fileMeta)
}

func (h *uploadHandler) AbortHandler(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	sessionID, err := uuid.Parse(mux.Vars(r)["sessionID"])
	if err != nil {
		return utils.ErrBadRequest.Wrap(err, "invalid session ID")
	}

	if err := h.uploadService.Abort(ctx, sessionID); err != nil {
		return err
	}

	return middleware.WriteJSONResponse(w, http.StatusOK, map[string]interface{}{
		"message": "upload session aborted",
	})
}
