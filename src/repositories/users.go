package repositories

import (
	"context"

	"github.com/CustomCloudStorage/types"
	"github.com/CustomCloudStorage/utils"
)

func (postgres *Postgres) GetUser(ctx context.Context, id int) (*types.User, error) {
	var user types.User
	if err := postgres.Db.WithContext(ctx).
		Preload("Profile").
		Preload("Account").
		Preload("Credentials").
		Where("id = ?", id).
		First(&user).Error; err != nil {
		return nil, utils.DetermineSQLError(err, "get data by id")
	}
	return &user, nil
}

func (postgres *Postgres) GetAllUsers(ctx context.Context) ([]types.User, error) {
	var users []types.User
	if err := postgres.Db.WithContext(ctx).
		Preload("Profile").
		Preload("Account").
		Preload("Credentials").
		Find(&users).Error; err != nil {
		return nil, utils.DetermineSQLError(err, "get all data")
	}
	return users, nil
}

func (postgres *Postgres) CreateUser(ctx context.Context, user *types.User) error {
	if err := postgres.Db.WithContext(ctx).Create(user).Error; err != nil {
		return utils.DetermineSQLError(err, "create data")
	}
	return nil
}

func (postgres *Postgres) UpdateProfile(ctx context.Context, profile *types.Profile, id int) error {
	if err := postgres.Db.WithContext(ctx).
		Model(&types.Profile{}).
		Where("user_id = ?", id).
		Updates(profile).Error; err != nil {
		return utils.DetermineSQLError(err, "update profile data")
	}
	return nil
}
