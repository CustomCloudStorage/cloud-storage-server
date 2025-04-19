package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/CustomCloudStorage/types"
	"github.com/CustomCloudStorage/utils"
	"github.com/gorilla/mux"
)

func (h *Handler) HandleGetUser(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	params := mux.Vars(r)

	log.Println("[GET] Fetching user with ID:", params["id"])

	id, err := strconv.Atoi(params["id"])
	if err != nil {
		return utils.ErrConversion.Wrap(err, "failed to convert ID to int")
	}

	user, err := h.Repository.User.GetByID(ctx, id)
	if err != nil {
		return err
	}

	publicUser := types.NewPublicUser(user)

	if err := json.NewEncoder(w).Encode(publicUser); err != nil {
		return utils.ErrJsonEncode.Wrap(err, "failed to encode user %s to JSON", params["id"])
	}

	return nil
}

func (h *Handler) HandleListUsers(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	log.Println("[GET] Fetching all users")

	users, err := h.Repository.User.List(ctx)
	if err != nil {
		return err
	}

	publicUsers := types.NewPublicUsers(users)

	if err := json.NewEncoder(w).Encode(publicUsers); err != nil {
		return utils.ErrJsonEncode.Wrap(err, "failed to encode users to JSON")
	}

	return nil
}

func (h *Handler) HandleCreateUser(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	log.Println("[POST] Creating new user")

	var user types.User

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		return utils.ErrJsonDecode.Wrap(err, "failed to decode json into the user's struct")
	}

	securePass, err := utils.HashPassword(user.Credentials.Password)
	if err != nil {
		return utils.ErrHash.Wrap(err, "Failed to hash the password")
	}

	user.Credentials.Password = securePass
	user.Account.UsedStorage = 0

	if err := h.Repository.User.Create(ctx, &user); err != nil {
		return err
	}

	writeJSONResponse(w, http.StatusCreated, map[string]string{
		"success": "User created successfully",
	})

	return nil
}

func (h *Handler) HandleUpdateProfile(w http.ResponseWriter, r *http.Request) error {
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

	user, err := h.Repository.User.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if profile.UpdatedAt != user.Profile.UpdatedAt {
		return utils.ErrDataConflict.New("The profile was changed by another user")
	}

	if err := h.Repository.User.UpdateProfile(ctx, &profile, id); err != nil {
		return err
	}

	writeJSONResponse(w, http.StatusOK, map[string]string{
		"success": "Profile updated successfully",
	})

	return nil
}

func (h *Handler) HandleUpdateAccount(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	params := mux.Vars(r)

	log.Println("[PUT] Updating user's account with id:", params["id"])

	id, err := strconv.Atoi(params["id"])
	if err != nil {
		return utils.ErrConversion.Wrap(err, "failed to convert ID to int")
	}

	var account types.Account
	if err := json.NewDecoder(r.Body).Decode(&account); err != nil {
		return utils.ErrJsonDecode.Wrap(err, "failed to decode json into the account's struct")
	}

	user, err := h.Repository.User.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if account.UpdatedAt != user.Account.UpdatedAt {
		return utils.ErrDataConflict.New("The account was changed by another user")
	}

	if err := h.Repository.User.UpdateAccount(ctx, &account, id); err != nil {
		return err
	}

	writeJSONResponse(w, http.StatusOK, map[string]string{
		"success": "Account updated successfully",
	})

	return nil
}

func (h *Handler) HandleUpdateCredentials(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	params := mux.Vars(r)

	log.Println("[PUT] Updating user's credentials with id:", params["id"])

	id, err := strconv.Atoi(params["id"])
	if err != nil {
		return utils.ErrConversion.Wrap(err, "failed to convert ID to int")
	}

	var credentials types.Credentials
	if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
		return utils.ErrJsonDecode.Wrap(err, "failed to decode json into the credentials struct")
	}

	securePass, err := utils.HashPassword(credentials.Password)
	if err != nil {
		return utils.ErrHash.Wrap(err, "Failed to hash the password")
	}

	credentials.Password = securePass

	user, err := h.Repository.User.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if credentials.UpdatedAt != user.Credentials.UpdatedAt {
		return utils.ErrDataConflict.New("The credentials were changed by another user")
	}

	if err := h.Repository.User.UpdateCredentials(ctx, &credentials, id); err != nil {
		return err
	}

	writeJSONResponse(w, http.StatusOK, map[string]string{
		"success": "Credentials updated successfully",
	})

	return nil
}

func (h *Handler) HandleDeleteUser(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	params := mux.Vars(r)

	log.Println("[DELETE] Deleting user with id:", params["id"])

	id, err := strconv.Atoi(params["id"])
	if err != nil {
		return utils.ErrConversion.Wrap(err, "failed to convert ID to int")
	}

	files, err := h.Repository.File.ListByUserID(ctx, id)
	if err != nil {
		return err
	}
	if len(files) != 0 {
		for _, file := range files {
			if err := h.Service.File.DeleteFile(ctx, file.ID, id); err != nil {
				return err
			}
		}
	}

	if err := h.Repository.User.Delete(ctx, id); err != nil {
		return err
	}

	writeJSONResponse(w, http.StatusOK, map[string]string{
		"success": "User successfully deleted",
	})

	return nil
}
