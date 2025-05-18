package middleware

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/CustomCloudStorage/utils"
	"github.com/joomcode/errorx"
)

type HandlerWithErrorFunc func(w http.ResponseWriter, r *http.Request) error

func HandleError(handler HandlerWithErrorFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := handler(w, r); err != nil {
			log.Println(err.Error())
			var (
				status  int
				message string
			)
			switch {
			case errorx.IsOfType(err, utils.ErrBadRequest):
				status = http.StatusBadRequest
				message = "Bad request"
			case errorx.IsOfType(err, utils.ErrNotFound):
				status = http.StatusNotFound
				message = "Resource not found"
			case errorx.IsOfType(err, utils.ErrConflict):
				status = http.StatusConflict
				message = "Conflict occurred"
			case errorx.IsOfType(err, utils.ErrUnauthorized):
				status = http.StatusUnauthorized
				message = "Authentication required"
			case errorx.IsOfType(err, utils.ErrForbidden):
				status = http.StatusForbidden
				message = "Access denied"
			default:
				status = http.StatusInternalServerError
				message = "Internal server error"
			}

			WriteJSONResponse(w, status, map[string]interface{}{
				"error": message,
			})
		}
	}
}

func WriteJSONResponse(w http.ResponseWriter, httpCode int, payload interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpCode)
	return json.NewEncoder(w).Encode(payload)
}
