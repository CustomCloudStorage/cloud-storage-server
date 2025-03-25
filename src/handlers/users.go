package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/CustomCloudStorage/types"
	"github.com/CustomCloudStorage/utils"
	"github.com/gorilla/mux"
)

func (handler *Handler) HandleGetUser(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	params := mux.Vars(r)

	log.Println("[GET] Fetching user with ID:", params["id"])

	user, err := handler.Repository.Postgres.GetUser(ctx, params["id"])
	if err != nil {
		return err
	}

	if err := json.NewEncoder(w).Encode(user); err != nil {
		return utils.ErrJsonEncode.Wrap(err, "failed to encode user %s to JSON", params["id"])
	}

	return nil
}

func (handler *Handler) HandleGetAllUsers(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	log.Println("[GET] Fetching all users")

	users, err := handler.Repository.Postgres.GetAllUsers(ctx)
	if err != nil {
		return err
	}

	if err := json.NewEncoder(w).Encode(users); err != nil {
		return utils.ErrJsonEncode.Wrap(err, "failed to encode users to JSON")
	}

	return nil
}

func (handler *Handler) HandleCreateUser(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	log.Println("[POST] Creating new user")

	var user types.User

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		return utils.ErrJsonDecode.Wrap(err, "failed to decode json into the user's struct")
	}

	time, err := utils.FormatDateTime(time.Now())
	if err != nil {
		return utils.ErrFormat.Wrap(err, "failed to get the current time in the format")
	}

	user.LastUpdate = time

	id, err := handler.Repository.Postgres.CreateUser(ctx, &user)
	if err != nil {
		return err
	}

	writeErrorResponse(w, http.StatusCreated, map[string]string{
		"success": "User created successfully",
		"user_id": id,
	})

	return nil
}
