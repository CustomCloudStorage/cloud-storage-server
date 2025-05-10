package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/CustomCloudStorage/utils"
	"github.com/gorilla/mux"
)

func (h *fileHandler) HandleGetFile(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	userID, ok := ctx.Value("userID").(int)
	if !ok {
		return utils.ErrBadRequest.New("user not authenticated")
	}

	params := mux.Vars(r)
	fileID, err := strconv.Atoi(params["fileID"])
	if err != nil {
		return utils.ErrBadRequest.Wrap(err, "invalid file ID")
	}

	file, err := h.fileRepository.GetByID(ctx, fileID, userID)
	if err != nil {
		return err
	}

	return WriteJSONResponse(w, http.StatusOK, map[string]interface{}{
		"file": file,
	})
}

func (h *fileHandler) HandleListFiles(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	userID, ok := ctx.Value("userID").(int)
	if !ok {
		return utils.ErrBadRequest.New("user not authenticated")
	}

	files, err := h.fileRepository.ListByUserID(ctx, userID)
	if err != nil {
		return err
	}

	return WriteJSONResponse(w, http.StatusOK, map[string]interface{}{
		"files": files,
	})
}

func (h *fileHandler) HandleUpdateName(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	userID, ok := ctx.Value("userID").(int)
	if !ok {
		return utils.ErrBadRequest.New("user not authenticated")
	}

	params := mux.Vars(r)
	fileID, err := strconv.Atoi(params["fileID"])
	if err != nil {
		return utils.ErrBadRequest.Wrap(err, "invalid file ID")
	}

	var payload struct{ Name string }
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		return utils.ErrBadRequest.Wrap(err, "invalid JSON payload")
	}

	if err := h.fileRepository.UpdateName(ctx, fileID, userID, payload.Name); err != nil {
		return err
	}

	return WriteJSONResponse(w, http.StatusOK, map[string]interface{}{
		"message": "file name updated successfully",
	})
}

func (h *fileHandler) HandleUpdateFolderID(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	userID, ok := ctx.Value("userID").(int)
	if !ok {
		return utils.ErrBadRequest.New("user not authenticated")
	}

	params := mux.Vars(r)
	fileID, err := strconv.Atoi(params["fileID"])
	if err != nil {
		return utils.ErrBadRequest.Wrap(err, "invalid file ID")
	}

	var payload struct{ FolderID int }
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		return utils.ErrBadRequest.Wrap(err, "invalid JSON payload")
	}

	if err := h.fileRepository.UpdateFolder(ctx, fileID, userID, payload.FolderID); err != nil {
		return err
	}

	return WriteJSONResponse(w, http.StatusOK, map[string]interface{}{
		"message": "file moved successfully",
	})
}

func (h *fileHandler) DownloadURLHandler(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	userID, ok := ctx.Value("userID").(int)
	if !ok {
		return utils.ErrBadRequest.New("user not authenticated")
	}

	params := mux.Vars(r)
	fileID, err := strconv.Atoi(params["fileID"])
	if err != nil {
		return utils.ErrBadRequest.Wrap(err, "invalid file ID")
	}

	url, err := h.fileService.GenerateDownloadURL(ctx, userID, fileID)
	if err != nil {
		return err
	}

	return WriteJSONResponse(w, http.StatusOK, map[string]interface{}{
		"download_url": url,
	})
}

func (h *fileHandler) DownloadByTokenHandler(w http.ResponseWriter, r *http.Request) error {
	token := r.URL.Query().Get("download_url")
	userID, fileID, err := h.fileService.ValidateDownloadToken(token)
	if err != nil {
		return utils.ErrBadRequest.Wrap(err, "invalid download token")
	}

	dfile, err := h.fileService.DownloadFile(r.Context(), userID, fileID)
	if err != nil {
		return err
	}
	defer dfile.Reader.(io.Closer).Close()

	w.Header().Set("Content-Type", dfile.ContentType)
	w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, dfile.FileName))
	w.Header().Set("Content-Length", strconv.FormatInt(dfile.FileSize, 10))
	return nil
}

func (h *fileHandler) StreamFileHandler(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	userID, ok := ctx.Value("userID").(int)
	if !ok {
		return utils.ErrBadRequest.New("user not authenticated")
	}

	fileID, err := strconv.Atoi(mux.Vars(r)["fileID"])
	if err != nil {
		return utils.ErrBadRequest.Wrap(err, "invalid file ID")
	}

	dfile, err := h.fileService.DownloadFile(ctx, userID, fileID)
	if err != nil {
		return err
	}
	defer dfile.Reader.(io.Closer).Close()

	ifModifiedSince := r.Header.Get("If-Modified-Since")
	if t, err := time.Parse(http.TimeFormat, ifModifiedSince); err == nil {
		if !dfile.ModTime.After(t) {
			w.WriteHeader(http.StatusNotModified)
			return nil
		}
	}

	w.Header().Set("Content-Type", dfile.ContentType)
	w.Header().Set("Content-Disposition", fmt.Sprintf(`inline; filename="%s"`, dfile.FileName))
	w.Header().Set("Last-Modified", dfile.ModTime.UTC().Format(http.TimeFormat))
	w.Header().Set("Accept-Ranges", "bytes")

	http.ServeContent(w, r, dfile.FileName, dfile.ModTime, dfile.Reader)
	return nil
}

func (h *fileHandler) PreviewFileHandler(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	userID, ok := ctx.Value("userID").(int)
	if !ok {
		return utils.ErrBadRequest.New("user not authenticated")
	}

	fileID, err := strconv.Atoi(mux.Vars(r)["fileID"])
	if err != nil {
		return utils.ErrBadRequest.Wrap(err, "invalid file ID")
	}

	modTime, err := h.fileService.PreviewFile(ctx, userID, fileID, io.Discard)
	if err != nil {
		return err
	}

	ifModifiedSince := r.Header.Get("If-Modified-Since")
	if t, err := time.Parse(http.TimeFormat, ifModifiedSince); err == nil {
		if !modTime.After(t) {
			w.WriteHeader(http.StatusNotModified)
			return nil
		}
	}

	w.Header().Set("Content-Type", "image/jpeg")
	w.Header().Set("Cache-Control", "public, max-age=86400")
	w.Header().Set("Last-Modified", modTime.UTC().Format(http.TimeFormat))

	_, err = h.fileService.PreviewFile(ctx, userID, fileID, w)
	if err != nil {
		return utils.ErrInternal.Wrap(err, "generate preview")
	}
	return nil
}
