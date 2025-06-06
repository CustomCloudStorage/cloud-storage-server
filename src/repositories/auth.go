package repositories

import (
	"context"

	"github.com/CustomCloudStorage/types"
	"github.com/CustomCloudStorage/utils"
	"golang.org/x/crypto/bcrypt"
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
	AuthMe(ctx context.Context, id int) (*types.User, error)
}

func (r *authRepository) LogIn(ctx context.Context, email, password string) (*types.User, error) {
	var user types.User

	err := r.db.WithContext(ctx).
		Preload("Profile").
		Preload("Account").
		Preload("Credentials").
		Joins("JOIN profiles p ON p.user_id = users.id").
		Where("p.email = ?", email).
		First(&user).Error
	if err != nil {
		return nil, utils.ErrUnauthorized.Wrap(err, "get user by email")
	}
	err = bcrypt.CompareHashAndPassword(
		[]byte(user.Credentials.Password),
		[]byte(password),
	)
	if err != nil {
		return nil, utils.ErrUnauthorized.New("invalid email or password")
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

func (r *authRepository) AuthMe(ctx context.Context, id int) (*types.User, error) {
	var user types.User
	if err := r.db.WithContext(ctx).
		Preload("Profile").
		Preload("Account").
		Preload("Credentials").
		Where("id = ?", id).
		First(&user).Error; err != nil {
		return nil, utils.DetermineSQLError(err, "get user by id")
	}
	return &user, nil
}
