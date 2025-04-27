package handlers

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"

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

func (h *fileHandler) StreamFileHandler(w http.ResponseWriter, r *http.Request) error {
	params := mux.Vars(r)
	fileID, _ := strconv.Atoi(params["fileID"])
	// надо будет извлечь userID из токена

	dfile, err := h.fileService.DownloadFile(r.Context(), userID, fileID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return err
	}
	defer dfile.Reader.(io.Closer).Close()

	ifModifiedSince := r.Header.Get("If-Modified-Since")
	if t, err := time.Parse(http.TimeFormat, ifModifiedSince); err == nil {
		if !dfile.ModTime.After(t) {
			w.WriteHeader(http.StatusNotModified)
			return err
		}
	}

	w.Header().Set("Content-Type", dfile.ContentType)
	w.Header().Set("Content-Disposition", "inline; filename=\""+dfile.FileName+"\"")
	w.Header().Set("Last-Modified", dfile.ModTime.UTC().Format(http.TimeFormat))
	w.Header().Set("Accept-Ranges", "bytes")

	http.ServeContent(w, r, dfile.FileName, dfile.ModTime, dfile.Reader)

	return nil
}

func (h *fileHandler) PreviewFileHandler(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	fileID, _ := strconv.Atoi(vars["fileID"])
	// надо будет извлечь userID из токена

	modTime, err := h.fileService.PreviewFile(r.Context(), userID, fileID, io.Discard)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return err
	}
	ifModifiedSince := r.Header.Get("If-Modified-Since")
	if t, err := time.Parse(http.TimeFormat, ifModifiedSince); err == nil {
		if !modTime.After(t) {
			w.WriteHeader(http.StatusNotModified)
			return err
		}
	}

	w.Header().Set("Content-Type", "image/jpeg")
	w.Header().Set("Cache-Control", "public, max-age=86400")
	w.Header().Set("Last-Modified", modTime.UTC().Format(http.TimeFormat))

	_, err = h.fileService.PreviewFile(r.Context(), userID, fileID, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	return nil
}
