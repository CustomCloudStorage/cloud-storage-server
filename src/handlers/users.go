package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/CustomCloudStorage/middleware"
	"github.com/CustomCloudStorage/repositories"
	"github.com/CustomCloudStorage/services"
	"github.com/CustomCloudStorage/types"
	"github.com/CustomCloudStorage/utils"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
)

type UserHandler interface {
	HandleGetUser(w http.ResponseWriter, r *http.Request) error
	HandleListUsers(w http.ResponseWriter, r *http.Request) error
	HandleUpdateProfile(w http.ResponseWriter, r *http.Request) error
	HandleUpdateAccount(w http.ResponseWriter, r *http.Request) error
	HandleUpdateCredentials(w http.ResponseWriter, r *http.Request) error
	HandleDeleteUser(w http.ResponseWriter, r *http.Request) error
	HandleStorageStats(w http.ResponseWriter, r *http.Request) error
}

type userHandler struct {
	userRepository repositories.UserRepository
	fileRepository repositories.FileRepository
	fileService    services.FileService
	userService    services.UserService
}

func NewUserHandler(userRepository repositories.UserRepository, fileRepository repositories.FileRepository, fileService services.FileService, userService services.UserService) UserHandler {
	return &userHandler{
		userRepository: userRepository,
		fileRepository: fileRepository,
		fileService:    fileService,
		userService:    userService,
	}
}

func (h *userHandler) HandleGetUser(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	params := mux.Vars(r)

	id, err := strconv.Atoi(params["id"])
	if err != nil {
		return utils.ErrBadRequest.Wrap(err, "invalid user ID %q", params["id"])
	}

	user, err := h.userRepository.GetByID(ctx, id)
	if err != nil {
		return err
	}

	publicUser := types.NewPublicUser(user)
	return middleware.WriteJSONResponse(w, http.StatusOK, map[string]interface{}{
		"user": publicUser,
	})
}

func (h *userHandler) HandleListUsers(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	users, err := h.userRepository.List(ctx)
	if err != nil {
		return err
	}

	publicUsers := types.NewPublicUsers(users)
	return middleware.WriteJSONResponse(w, http.StatusOK, map[string]interface{}{
		"users": publicUsers,
	})
}

func (h *userHandler) HandleUpdateProfile(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	claims := ctx.Value("claims").(jwt.MapClaims)
	id, ok := claims["userID"].(int)
	if !ok {
		return middleware.WriteJSONResponse(w, http.StatusUnauthorized, map[string]string{
			"error": "invalid or expired token",
		})
	}

	var profile types.Profile
	if err := json.NewDecoder(r.Body).Decode(&profile); err != nil {
		return utils.ErrBadRequest.Wrap(err, "decode profile JSON")
	}

	user, err := h.userRepository.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if !profile.UpdatedAt.Equal(user.Profile.UpdatedAt) {
		return utils.ErrConflict.New("profile was changed by another user")
	}

	if err := h.userRepository.UpdateProfile(ctx, &profile, id); err != nil {
		return err
	}

	return middleware.WriteJSONResponse(w, http.StatusOK, map[string]interface{}{
		"message": "profile updated successfully",
	})
}

func (h *userHandler) HandleUpdateAccount(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	params := mux.Vars(r)

	id, err := strconv.Atoi(params["id"])
	if err != nil {
		return utils.ErrBadRequest.Wrap(err, "invalid user ID")
	}

	var account types.Account
	if err := json.NewDecoder(r.Body).Decode(&account); err != nil {
		return utils.ErrBadRequest.Wrap(err, "decode account JSON")
	}

	user, err := h.userRepository.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if !account.UpdatedAt.Equal(user.Account.UpdatedAt) {
		return utils.ErrConflict.New("account was changed by another user")
	}

	if err := h.userRepository.UpdateAccount(ctx, &account, id); err != nil {
		return err
	}

	return middleware.WriteJSONResponse(w, http.StatusOK, map[string]interface{}{
		"message": "account updated successfully",
	})
}

func (h *userHandler) HandleUpdateCredentials(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	claims := ctx.Value("claims").(jwt.MapClaims)
	id, ok := claims["userID"].(int)
	if !ok {
		return middleware.WriteJSONResponse(w, http.StatusUnauthorized, map[string]string{
			"error": "invalid or expired token",
		})
	}

	var creds types.Credentials
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		return utils.ErrBadRequest.Wrap(err, "decode credentials JSON")
	}

	securePass, err := utils.HashPassword(creds.Password)
	if err != nil {
		return utils.ErrInternal.Wrap(err, "hash password")
	}
	creds.Password = securePass

	user, err := h.userRepository.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if !creds.UpdatedAt.Equal(user.Credentials.UpdatedAt) {
		return utils.ErrConflict.New("credentials were changed by another user")
	}

	if err := h.userRepository.UpdateCredentials(ctx, &creds, id); err != nil {
		return err
	}

	return middleware.WriteJSONResponse(w, http.StatusOK, map[string]interface{}{
		"message": "credentials updated successfully",
	})
}

func (h *userHandler) HandleDeleteUser(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	params := mux.Vars(r)

	id, err := strconv.Atoi(params["id"])
	if err != nil {
		return utils.ErrBadRequest.Wrap(err, "invalid user ID")
	}

	files, err := h.fileRepository.ListByUserID(ctx, id)
	if err != nil {
		return err
	}
	for _, f := range files {
		if err := h.fileService.DeleteFile(ctx, f.ID, id); err != nil {
			return err
		}
	}

	if err := h.userRepository.Delete(ctx, id); err != nil {
		return err
	}

	return middleware.WriteJSONResponse(w, http.StatusOK, map[string]interface{}{
		"message": "user deleted successfully",
	})
}

func (h *userHandler) HandleStorageStats(w http.ResponseWriter, r *http.Request) error {
	stats, err := h.userService.StatsStorage(r.Context())
	if err != nil {
		return err
	}

	const toGB = 1.0 / (1024 * 1024 * 1024)
	out := struct {
		Total     float64 `json:"total_gb"`
		Allocated float64 `json:"allocated_gb"`
		Free      float64 `json:"free_gb"`
	}{
		Total:     float64(stats.TotalBytes) * toGB,
		Allocated: float64(stats.AllocatedBytes) * toGB,
		Free:      float64(stats.FreeBytes) * toGB,
	}
	return middleware.WriteJSONResponse(w, http.StatusOK, out)
}
