package middleware

import (
	"context"
	"net/http"
	"slices"

	"github.com/CustomCloudStorage/handlers"
	"github.com/CustomCloudStorage/services"
	"github.com/CustomCloudStorage/utils"
	"github.com/joomcode/errorx"
)

type authMiddleware struct {
	authService services.AuthService
	cfg         services.Auth
}

func NewAuthMiddleware(authService services.AuthService, cfg services.Auth) *authMiddleware {
	return &authMiddleware{
		authService: authService,
		cfg:         cfg,
	}
}

type AuthMiddleware interface {
	AuthMiddleWare() func(http.Handler) http.Handler
}

func (m *authMiddleware) AuthMiddleWare() func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			path := r.URL.Path
			if slices.Contains(m.cfg.Ignore, path) {
				h.ServeHTTP(w, r)
				return
			}
			token := r.Header.Get(m.cfg.Header)
			if token == "" {
				handlers.WriteJSONResponse(w, http.StatusUnauthorized, map[string]string{
					"error": "authorization token is required",
				})
				return
			}

			claims, err := m.authService.ValidateToken(ctx, token)
			if err != nil {
				switch {
				case errorx.IsOfType(err, utils.ErrUnauthorized):
					handlers.WriteJSONResponse(w, http.StatusUnauthorized, map[string]string{
						"error": "invalid or expired token",
					})
				case errorx.IsOfType(err, utils.ErrForbidden):
					handlers.WriteJSONResponse(w, http.StatusForbidden, map[string]string{
						"error": "access denied",
					})
				default:
					handlers.WriteJSONResponse(w, http.StatusInternalServerError, map[string]string{
						"error": "internal server error",
					})
				}
				return
			}

			context := context.WithValue(r.Context(), "claims", claims)
			h.ServeHTTP(w, r.WithContext(context))
		})
	}
}
