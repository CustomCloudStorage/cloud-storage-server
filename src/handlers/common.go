package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/CustomCloudStorage/repositories"
	"github.com/CustomCloudStorage/service"
	"github.com/CustomCloudStorage/utils"
	"github.com/joomcode/errorx"
)

type Handler struct {
	Repository *repositories.Repository
	Service    *service.Service
}

func NewHandler(repository *repositories.Repository, service *service.Service) *Handler {
	return &Handler{
		Repository: repository,
		Service:    service,
	}
}

type HandlerWithErrorFunc func(w http.ResponseWriter, r *http.Request) error

func HandleError(handler HandlerWithErrorFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := handler(w, r)
		if err != nil {
			switch {
			case errorx.IsOfType(err, utils.ErrNotFound):
				log.Println(err.Error())
				writeJSONResponse(w, http.StatusNotFound, map[string]string{"error": "Data not found"})
			case errorx.IsOfType(err, utils.ErrAlreadyExist):
				log.Println(err.Error())
				writeJSONResponse(w, http.StatusBadRequest, map[string]string{"error": "The request body is invalid"})
			case errorx.IsOfType(err, utils.ErrDataConflict):
				log.Println(err.Error())
				writeJSONResponse(w, http.StatusConflict, map[string]string{
					"сonflict": "The data was changed",
				})
			default:
				log.Println("Internal server error:", err.Error())
				writeJSONResponse(w, http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
			}
		}
	}
}

func writeJSONResponse(w http.ResponseWriter, httpCode int, message map[string]string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpCode)
	json.NewEncoder(w).Encode(message)
}
