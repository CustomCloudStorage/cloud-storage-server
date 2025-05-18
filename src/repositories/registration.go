package repositories

import (
	"context"
	"time"

	"github.com/CustomCloudStorage/types"
	"github.com/CustomCloudStorage/utils"
	"gorm.io/gorm"
)

type RegistrationRepository interface {
	Create(ctx context.Context, r *types.Registration) error
	GetByEmail(ctx context.Context, email string) (*types.Registration, error)
	UpdateCode(ctx context.Context, email, newCode string, now time.Time) error
	Delete(ctx context.Context, email string) error
	DeleteExpired(ctx context.Context, olderThan time.Duration) error
}

type registrationRepository struct {
	db *gorm.DB
}

func NewRegistrationRepository(db *gorm.DB) RegistrationRepository {
	return &registrationRepository{db: db}
}

func (r *registrationRepository) Create(ctx context.Context, reg *types.Registration) error {
	if err := r.db.WithContext(ctx).Create(reg).Error; err != nil {
		return utils.DetermineSQLError(err, "create registration")
	}
	return nil
}

func (r *registrationRepository) GetByEmail(ctx context.Context, email string) (*types.Registration, error) {
	var reg types.Registration
	err := r.db.WithContext(ctx).
		Where("email = ?", email).
		First(&reg).Error
	if err != nil {
		return nil, utils.DetermineSQLError(err, "fetch registration")
	}
	return &reg, nil
}

func (r *registrationRepository) UpdateCode(ctx context.Context, email, newCode string, now time.Time) error {
	result := r.db.WithContext(ctx).
		Model(&types.Registration{}).
		Where("email = ?", email).
		Updates(map[string]interface{}{
			"code":         newCode,
			"last_sent_at": now,
		})
	if err := result.Error; err != nil {
		return utils.DetermineSQLError(err, "update confirmation code")
	}
	return nil
}

func (r *registrationRepository) Delete(ctx context.Context, email string) error {
	err := r.db.WithContext(ctx).
		Where("email = ?", email).
		Delete(&types.Registration{}).
		Error
	if err != nil {
		return utils.DetermineSQLError(err, "delete registration")
	}
	return nil
}

func (r *registrationRepository) DeleteExpired(ctx context.Context, olderThan time.Duration) error {
	cutoff := time.Now().Add(-olderThan)
	if err := r.db.WithContext(ctx).
		Where("created_at <= ?", cutoff).
		Delete(&types.Registration{}).
		Error; err != nil {
		return utils.DetermineSQLError(err, "delete expired registrations")
	}
	return nil
}
