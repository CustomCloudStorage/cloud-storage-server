package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/CustomCloudStorage/middleware"
	"github.com/CustomCloudStorage/repositories"
	"github.com/CustomCloudStorage/services"
	"github.com/CustomCloudStorage/utils"
)

type RegistrationHandler interface {
	Register(w http.ResponseWriter, r *http.Request) error
	Confirm(w http.ResponseWriter, r *http.Request) error
	ResendCode(w http.ResponseWriter, r *http.Request) error
}

type registrationHandler struct {
	registrationRepository repositories.RegistrationRepository
	registrationService    services.RegistrationService
}

func NewRegistrationHandler(registrationRepo repositories.RegistrationRepository, registrationService services.RegistrationService) RegistrationHandler {
	return &registrationHandler{
		registrationRepository: registrationRepo,
		registrationService:    registrationService,
	}
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

	var req struct {
		Email string `json:"email"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return utils.ErrBadRequest.Wrap(err, "decode email")
	}

	reg, err := h.registrationRepository.GetByEmail(ctx, req.Email)
	if err != nil {
		return err
	}
	if time.Since(reg.LastSentAt) < (1 * time.Minute) {
		nextAllowed := reg.LastSentAt.Add(1 * time.Minute)
		wait := time.Until(nextAllowed)
		fmt.Printf("Подождите %d секунд", wait)
		return middleware.WriteJSONResponse(
			w,
			http.StatusTooManyRequests,
			map[string]interface{}{
				"error":      fmt.Sprintf("Подождите %d секунд", wait),
				"retryAfter": wait,
			},
		)
	}

	if err = h.registrationService.ResendCode(ctx, reg.Email); err != nil {
		return err
	}

	return middleware.WriteJSONResponse(w, http.StatusOK, map[string]interface{}{"message": "confirmation code resent"})
}
