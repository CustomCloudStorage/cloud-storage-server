package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/CustomCloudStorage/repositories"
	"github.com/CustomCloudStorage/services"
	"github.com/CustomCloudStorage/utils"
)

type authHandler struct {
	authRepository repositories.AuthRepository
	authservice    services.AuthService
}

func NewAuthHandler(authRepo repositories.AuthRepository, authService services.AuthService) *authHandler {
	return &authHandler{
		authRepository: authRepo,
		authservice:    authService,
	}
}

type AuthHandler interface {
	HandleLogIn(w http.ResponseWriter, r *http.Request) error
}

func (h *authHandler) HandleLogIn(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return utils.ErrBadRequest.Wrap(err, "decode login payload")
	}

	if req.Email == "" || req.Password == "" {
		return utils.ErrInternal.New("failed to find account")
	}

	if err := h.authservice.LogInService(ctx, req.Email, req.Password); err != nil {
		return nil
	}
	return nil
}
