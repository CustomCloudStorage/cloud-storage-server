package repositories

import (
	"context"

	"github.com/CustomCloudStorage/types"
	"github.com/CustomCloudStorage/utils"
	"gorm.io/gorm"
)

type authRepository struct {
	db *gorm.DB
}

func NewAuthRepository(db *gorm.DB) *authRepository {
	return &authRepository{
		db: db,
	}
}

type AuthRepository interface {
	LogIn(ctx context.Context, email, password string) (*types.User, error)
	GetRole(ctx context.Context, email string) (string, error)
}

func (r *authRepository) LogIn(ctx context.Context, email, password string) (*types.User, error) {
	var user types.User
	if err := r.db.WithContext(ctx).
		Preload("Profile").
		Preload("Account").
		Preload("Credentials").
		Where("email = ?, password = ?", email, password).
		First(&user).Error; err != nil {
		return nil, utils.DetermineSQLError(err, "get user by email")
	}
	return &user, nil
}

func (r *authRepository) GetRole(ctx context.Context, email string) (string, error) {
	var role string
	if err := r.db.WithContext(ctx).
		Table("profiles p").
		Select("a.role").
		Joins("JOIN accounts a ON a.user_id = p.user_id").
		Where("p.email = ?", email).
		Scan(&role).Error; err != nil {
		return "", utils.DetermineSQLError(err, "get role by email")
	}
	return role, nil
}
