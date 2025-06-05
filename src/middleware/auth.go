package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/CustomCloudStorage/repositories"
	"github.com/CustomCloudStorage/services"
	"github.com/CustomCloudStorage/utils"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"github.com/joomcode/errorx"
)

type authMiddleware struct {
	authRepository repositories.AuthRepository
	authService    services.AuthService
	cfg            services.Auth
}

func NewAuthMiddleware(authRepo repositories.AuthRepository, authService services.AuthService, cfg services.Auth) *authMiddleware {
	return &authMiddleware{
		authRepository: authRepo,
		authService:    authService,
		cfg:            cfg,
	}
}

type AuthMiddleware interface {
	AuthMiddleWare() func(http.Handler) http.Handler
	RequireRole(allowedRoles ...string) mux.MiddlewareFunc
}

func (m *authMiddleware) AuthMiddleWare() func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			path := r.URL.Path
			for _, ign := range m.cfg.Ignore {
				if strings.HasPrefix(path, ign) {
					h.ServeHTTP(w, r)
					return
				}
			}
			token := r.Header.Get(m.cfg.Header)
			if token == "" {
				WriteJSONResponse(w, http.StatusUnauthorized, map[string]string{
					"error": "authorization token is required",
				})
				return
			}

			claims, err := m.authService.ValidateToken(ctx, token)
			if err != nil {
				switch {
				case errorx.IsOfType(err, utils.ErrUnauthorized):
					WriteJSONResponse(w, http.StatusUnauthorized, map[string]string{
						"error": "invalid or expired token",
					})
				case errorx.IsOfType(err, utils.ErrForbidden):
					WriteJSONResponse(w, http.StatusForbidden, map[string]string{
						"error": "access denied",
					})
				default:
					WriteJSONResponse(w, http.StatusInternalServerError, map[string]string{
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

func (m *authMiddleware) RequireRole(allowedRoles ...string) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			claims := ctx.Value("claims").(jwt.MapClaims)

			email, ok := claims["email"].(string)
			if !ok || email == "" {
				WriteJSONResponse(w, http.StatusUnauthorized, map[string]string{
					"error": "invalid or expired token",
				})
				return
			}

			role, err := m.authRepository.GetRole(ctx, email)
			if err != nil {
				fmt.Printf("RequireRole: could not get role for %s: %v", email, err)
				WriteJSONResponse(w, http.StatusInternalServerError, map[string]string{
					"error": "internal server error",
				})
				return
			}

			for _, allowed := range allowedRoles {
				if role == allowed {
					next.ServeHTTP(w, r)
					return
				}
			}

			WriteJSONResponse(w, http.StatusForbidden, map[string]string{
				"error": "access denied",
			})
		})
	}
}
