package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/CustomCloudStorage/middleware"
	"github.com/CustomCloudStorage/repositories"
	"github.com/CustomCloudStorage/services"
	"github.com/CustomCloudStorage/types"
	"github.com/CustomCloudStorage/utils"
	"github.com/dgrijalva/jwt-go"
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
	HandleLogOut(w http.ResponseWriter, r *http.Request) error
	HandleAuthMe(w http.ResponseWriter, r *http.Request) error
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

func (h *authHandler) HandleLogOut(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	claims := ctx.Value("claims").(jwt.MapClaims)

	email, ok := claims["email"].(string)
	if !ok || email == "" {
		return middleware.WriteJSONResponse(w, http.StatusUnauthorized, map[string]string{
			"error": "invalid or expired token",
		})
	}

	if err := h.authservice.LogOut(ctx, email); err != nil {
		return err
	}

	return middleware.WriteJSONResponse(w, http.StatusOK, map[string]interface{}{
		"message": "successfully logged out",
	})
}

func (h *authHandler) HandleAuthMe(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	claims := ctx.Value("claims").(jwt.MapClaims)

	userID, ok := claims["userID"].(int)
	if !ok {
		return middleware.WriteJSONResponse(w, http.StatusUnauthorized, map[string]string{
			"error": "invalid or expired token",
		})
	}

	user, err := h.authRepository.AuthMe(ctx, userID)
	if err != nil {
		return err
	}

	publicUser := types.NewPublicUser(user)
	return middleware.WriteJSONResponse(w, http.StatusOK, map[string]interface{}{
		"user": publicUser,
	})
}
