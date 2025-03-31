package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/CustomCloudStorage/types"
	"github.com/CustomCloudStorage/utils"
	"github.com/gorilla/mux"
)

func (handler *Handler) HandleGetUser(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	params := mux.Vars(r)

	log.Println("[GET] Fetching user with ID:", params["id"])

	id, err := strconv.Atoi(params["id"])
	if err != nil {
		return utils.ErrConversion.Wrap(err, "failed to convert ID to int")
	}

	user, err := handler.Repository.Postgres.GetUser(ctx, id)
	if err != nil {
		return err
	}

	publicUser := types.NewPublicUser(user)

	if err := json.NewEncoder(w).Encode(publicUser); err != nil {
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

	publicUsers := types.NewPublicUsers(users)

	if err := json.NewEncoder(w).Encode(publicUsers); err != nil {
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

	user.Profile.LastUpdateProfile = time
	user.Account.LastUpdateAccount = time
	user.Credentials.LastUpdateCredentials = time

	if err = handler.Repository.Postgres.CreateUser(ctx, &user); err != nil {
		return err
	}

	writeJSONResponse(w, http.StatusCreated, map[string]string{
		"success": "User created successfully",
	})

	return nil
}

func (handler *Handler) HandleUpdateProfile(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	params := mux.Vars(r)

	log.Println("[PUT] Updating user`s profile with id:", params["id"])

	id, err := strconv.Atoi(params["id"])
	if err != nil {
		return utils.ErrConversion.Wrap(err, "failed to convert ID to int")
	}

	var profile types.Profile
	if err := json.NewDecoder(r.Body).Decode(&profile); err != nil {
		return utils.ErrJsonDecode.Wrap(err, "failed to decode json into the profile's struct")
	}

	user, err := handler.Repository.Postgres.GetUser(ctx, id)
	if err != nil {
		return err
	}

	if profile.LastUpdateProfile != user.Profile.LastUpdateProfile {
		return utils.ErrDataConflict.New("The profile was changed by another user")
	}

	time, err := utils.FormatDateTime(time.Now())
	if err != nil {
		return utils.ErrFormat.Wrap(err, "failed to get the current time in the format")
	}

	profile.LastUpdateProfile = time

	if err := handler.Repository.Postgres.UpdateProfile(ctx, &profile, id); err != nil {
		return err
	}

	writeJSONResponse(w, http.StatusOK, map[string]string{
		"success": "Profile updated successfully",
	})

	return nil
}
