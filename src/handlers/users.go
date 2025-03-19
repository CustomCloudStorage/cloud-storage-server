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

func (handler *Handler) HandleGetAllUsers(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	log.Println("GET Request to get all users")

	users, err := handler.Repository.Postgres.GetAllUsers(ctx)
	if err != nil {
		return utils.ErrGet.Wrap(err, "failed to get all users")
	}

	if err := json.NewEncoder(w).Encode(users); err != nil {
		return utils.ErrJsonEncode.Wrap(err, "failed to encode users to JSON")
	}

	return nil
}
