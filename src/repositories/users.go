package repositories

import (
	"context"

	"github.com/CustomCloudStorage/types"
	"github.com/CustomCloudStorage/utils"
)

func (u *user) GetByID(ctx context.Context, id int) (*types.User, error) {
	var user types.User
	if err := u.db.WithContext(ctx).
		Preload("Profile").
		Preload("Account").
		Preload("Credentials").
		Where("id = ?", id).
		First(&user).Error; err != nil {
		return nil, utils.DetermineSQLError(err, "get data by id")
	}
	return &user, nil
}

func (u *user) Create(ctx context.Context, user *types.User) error {
	if err := u.db.WithContext(ctx).Create(user).Error; err != nil {
		return utils.DetermineSQLError(err, "create data")
	}
	return nil
}

func (u *user) UpdateProfile(ctx context.Context, profile *types.Profile, id int) error {
	if err := u.db.WithContext(ctx).
		Model(&types.Profile{}).
		Where("user_id = ?", id).
		Updates(profile).Error; err != nil {
		return utils.DetermineSQLError(err, "update profile data")
	}
	return nil
}

func (u *user) UpdateAccount(ctx context.Context, account *types.Account, id int) error {
	if err := u.db.WithContext(ctx).
		Model(&types.Account{}).
		Where("user_id = ?", id).
		Updates(account).Error; err != nil {
		return utils.DetermineSQLError(err, "update account data")
	}
	return nil
}

func (u *user) UpdateCredentials(ctx context.Context, credentials *types.Credentials, id int) error {
	if err := u.db.WithContext(ctx).
		Model(&types.Credentials{}).
		Where("user_id = ?", id).
		Updates(credentials).Error; err != nil {
		return utils.DetermineSQLError(err, "update credentials data")
	}
	return nil
}

func (u *user) Delete(ctx context.Context, id int) error {
	if err := u.db.WithContext(ctx).Delete(&types.User{}, id).Error; err != nil {
		return utils.DetermineSQLError(err, "delete user data")
	}
	return nil
}

func (u *user) List(ctx context.Context) ([]types.User, error) {
	var users []types.User
	if err := u.db.WithContext(ctx).
		Preload("Profile").
		Preload("Account").
		Preload("Credentials").
		Find(&users).Error; err != nil {
		return nil, utils.DetermineSQLError(err, "get all data")
	}
	return users, nil
}

func (u *user) UpdateUsedStorage(ctx context.Context, id int, newUsedStorage int64) error {
	if err := u.db.WithContext(ctx).
		Model(&types.Account{}).
		Where("user_id = ?", id).
		Update("used_storage", newUsedStorage).Error; err != nil {
		return utils.DetermineSQLError(err, "update used storage data")
	}
	return nil
}
