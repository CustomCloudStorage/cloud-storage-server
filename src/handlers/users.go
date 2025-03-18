package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/CustomCloudStorage/utils"
	"github.com/gorilla/mux"
)

func (handler *Handler) HandleGetUser(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	params := mux.Vars(r)

	log.Println("GET Request to get user with ID:", params["id"])

	user, err := handler.Repository.Postgres.GetUser(ctx, params["id"])
	if err != nil {
		return utils.ErrGet.Wrap(err, "failed to get user %s", params["id"])
	}

	if err := json.NewEncoder(w).Encode(user); err != nil {
		return utils.ErrJsonEncode.Wrap(err, "failed to encode user %s to JSON", params["id"])
	}

	return nil
}
