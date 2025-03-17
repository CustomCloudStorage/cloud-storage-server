package handlers

import (
	"encoding/json"
	"log"
	"net/http"
)

type HandlerWithErrorFunc func(w http.ResponseWriter, r *http.Request) error

func HandleError(handler HandlerWithErrorFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := handler(w, r)
		if err != nil {
			switch {
			//Здесь будут описаны остальные ошибки
			default:
				log.Println("Internal server error:", err.Error())
				writeErrorResponse(w, http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
			}
		}
	}
}

func writeErrorResponse(w http.ResponseWriter, httpCode int, message map[string]string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpCode)
	json.NewEncoder(w).Encode(message)
}
