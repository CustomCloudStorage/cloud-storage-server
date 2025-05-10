package services

import (
	"context"
	"fmt"
	"time"

	"github.com/CustomCloudStorage/repositories"
	"github.com/CustomCloudStorage/utils"
	"github.com/dgrijalva/jwt-go"
)

type authService struct {
	authRepository repositories.AuthRepository
	redis          repositories.RedisCache
	cfg            Auth
}

type Auth struct {
	Secret string
	Header string
	Ignore []string
}

func NewAuthService(authRepo repositories.AuthRepository, redis repositories.RedisCache, cfg Auth) *authService {
	return &authService{
		authRepository: authRepo,
		redis:          redis,
		cfg:            cfg,
	}
}

type AuthService interface {
	LogInService(ctx context.Context, email, password string) error
	ValidateToken(ctx context.Context, signedToken string) (jwt.MapClaims, error)
}

func (s *authService) LogInService(ctx context.Context, email, password string) error {
	hashPass, err := utils.HashPassword(password)
	if err != nil {
		return err
	}

	user, err := s.authRepository.LogIn(ctx, email, hashPass)
	if err != nil {
		return err
	}

	claims := jwt.MapClaims{
		"userID":  user.Id,
		"email":   email,
		"exp":     time.Now().Add(time.Hour * 72).Unix(),
		"expired": "false",
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(s.cfg.Secret))
	if err != nil {
		return utils.ErrInternal.Wrap(err, "failed to sign JWT token")
	}

	key := fmt.Sprintf("TOKEN_%s", email)

	if err := s.redis.Set(ctx, key, signedToken, 72*time.Hour); err != nil {
		return err
	}

	return nil
}

func (s *authService) ValidateToken(ctx context.Context, signedToken string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(signedToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, utils.ErrBadRequest.New("unexpected signing method")
		}
		return []byte(s.cfg.Secret), nil
	})
	if err != nil {
		return nil, utils.ErrUnauthorized.Wrap(err, "invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, utils.ErrUnauthorized.New("invalid or expired token")
	}

	rawEmail, ok := claims["email"]
	email, okCast := rawEmail.(string)
	if !ok || !okCast || email == "" {
		return nil, utils.ErrUnauthorized.New("token missing email claim")
	}

	key := fmt.Sprintf("TOKEN_%s", email)
	ok, err = s.redis.Exists(ctx, key)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, utils.ErrUnauthorized.New("session expired")
	}

	return claims, nil
}
