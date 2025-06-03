package services

import (
	"context"
	"fmt"
	"time"

	"github.com/CustomCloudStorage/repositories"
	"github.com/CustomCloudStorage/types"
	"github.com/CustomCloudStorage/utils"
)

const (
	codeTTL       = 3 * time.Hour
	codeLength    = 6
	purgeInterval = 1 * time.Hour
)

type RegistrationService interface {
	Register(ctx context.Context, email, password string) error
	Confirm(ctx context.Context, email, code string) error
	ResendCode(ctx context.Context, email string) error
}

type registrationService struct {
	registrationRepository repositories.RegistrationRepository
	userRepository         repositories.UserRepository
	emailService           EmailService
	cfg                    ServiceConfig
}

func NewRegistrationService(registrationRepo repositories.RegistrationRepository, userRepo repositories.UserRepository, emailService EmailService, cfg ServiceConfig) RegistrationService {
	svc := &registrationService{
		registrationRepository: registrationRepo,
		userRepository:         userRepo,
		emailService:           emailService,
		cfg:                    cfg,
	}

	go func() {
		ticker := time.NewTicker(purgeInterval)
		defer ticker.Stop()
		for range ticker.C {
			if err := svc.registrationRepository.DeleteExpired(context.Background(), codeTTL); err != nil {
				fmt.Printf("[registrationService] delete expired registrations error: %v\n", err)
			}
		}
	}()

	return svc
}

func (s *registrationService) Register(ctx context.Context, email, password string) error {
	hash, err := utils.HashPassword(password)
	if err != nil {
		return err
	}
	code, err := utils.GenerateCode(codeLength)
	if err != nil {
		return utils.ErrInternal.Wrap(err, "generate confirmation code failed")
	}
	reg := &types.Registration{
		Email:    email,
		Password: hash,
		Code:     code,
	}
	if err := s.registrationRepository.Create(ctx, reg); err != nil {
		return err
	}
	return s.emailService.EnqueueEmail(ctx,
		"registration_confirmation",
		email,
		"Код подтверждения регистрации",
		map[string]interface{}{"Code": code},
	)
}

func (s *registrationService) Confirm(ctx context.Context, email, code string) error {
	reg, err := s.registrationRepository.GetByEmail(ctx, email)
	if err != nil {
		return err
	}

	if code != reg.Code {
		return utils.ErrUnauthorized.Wrap(nil, "confirmation code mismatch")
	}

	allocated, err := s.userRepository.SumActiveStorageLimit(ctx)
	if err != nil {
		return err
	}

	totalBytes := s.cfg.TotalStorageBytes()

	freeBytes := totalBytes - allocated

	var allocBytes int64
	if freeBytes >= s.cfg.UserAllocBytes() {
		allocBytes = s.cfg.UserAllocBytes()
	}

	user := &types.User{
		Profile: types.Profile{
			Name:  utils.UsernameFromEmail(email),
			Email: email,
		},
		Account: types.Account{
			Role:         "",
			StorageLimit: allocBytes,
			UsedStorage:  0,
		},
		Credentials: types.Credentials{
			Password: reg.Password,
		},
	}
	if err := s.userRepository.Create(ctx, user); err != nil {
		return err
	}
	if err := s.registrationRepository.Delete(ctx, email); err != nil {
		return utils.ErrInternal.Wrap(err, "delete registration after confirm failed")
	}
	return nil
}

func (s *registrationService) ResendCode(ctx context.Context, email string) error {
	code, err := utils.GenerateCode(codeLength)
	if err != nil {
		return utils.ErrInternal.Wrap(err, "generate confirmation code failed")
	}
	now := time.Now()
	if err := s.registrationRepository.UpdateCode(ctx, email, code, now); err != nil {
		return err
	}
	return s.emailService.EnqueueEmail(ctx,
		"registration_confirmation",
		email,
		"Код подтверждения регистрации",
		map[string]interface{}{"Code": code},
	)
}
