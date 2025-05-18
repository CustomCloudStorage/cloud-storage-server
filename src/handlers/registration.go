package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/CustomCloudStorage/middleware"
	"github.com/CustomCloudStorage/services"
	"github.com/CustomCloudStorage/utils"
)

type RegistrationHandler interface {
	Register(w http.ResponseWriter, r *http.Request) error
	Confirm(w http.ResponseWriter, r *http.Request) error
	ResendCode(w http.ResponseWriter, r *http.Request) error
}

type registrationHandler struct {
	registrationService services.RegistrationService
}

func NewRegistrationHandler(registrationService services.RegistrationService) RegistrationHandler {
	return &registrationHandler{registrationService: registrationService}
}

func (h *registrationHandler) Register(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return utils.ErrBadRequest.Wrap(err, "decode reg payload")
	}

	err := h.registrationService.Register(ctx, req.Email, req.Password)
	if err != nil {
		return err
	}

	return middleware.WriteJSONResponse(w, http.StatusCreated, map[string]interface{}{"message": "user created successfully"})
}

func (h *registrationHandler) Confirm(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	var req struct {
		Email string `json:"email"`
		Code  string `json:"code"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return utils.ErrBadRequest.Wrap(err, "decode confirm payload")
	}

	err := h.registrationService.Confirm(ctx, req.Email, req.Code)
	if err != nil {
		return err
	}

	return middleware.WriteJSONResponse(w, http.StatusOK, map[string]interface{}{"message": "registration confirmed"})
}

func (h *registrationHandler) ResendCode(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	var email string
	if err := json.NewDecoder(r.Body).Decode(&email); err != nil {
		return utils.ErrBadRequest.Wrap(err, "decode email")
	}

	err := h.registrationService.ResendCode(ctx, email)
	if err != nil {
		return err
	}

	return middleware.WriteJSONResponse(w, http.StatusOK, map[string]interface{}{"message": "confirmation code resent"})
}
