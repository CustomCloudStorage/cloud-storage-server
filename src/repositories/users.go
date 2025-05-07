package repositories

import (
	"context"

	"github.com/CustomCloudStorage/types"
	"github.com/CustomCloudStorage/utils"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (u *userRepository) GetByID(ctx context.Context, id int) (*types.User, error) {
	var user types.User
	if err := u.db.WithContext(ctx).
		Preload("Profile").
		Preload("Account").
		Preload("Credentials").
		Where("id = ?", id).
		First(&user).Error; err != nil {
		return nil, utils.DetermineSQLError(err, "get user by id")
	}
	return &user, nil
}

func (u *userRepository) Create(ctx context.Context, user *types.User) error {
	if err := u.db.WithContext(ctx).
		Create(user).Error; err != nil {
		return utils.DetermineSQLError(err, "create user")
	}
	return nil
}

func (u *userRepository) UpdateProfile(ctx context.Context, profile *types.Profile, id int) error {
	if err := u.db.WithContext(ctx).
		Model(&types.Profile{}).
		Where("user_id = ?", id).
		Updates(profile).Error; err != nil {
		return utils.DetermineSQLError(err, "update profile")
	}
	return nil
}

func (u *userRepository) UpdateAccount(ctx context.Context, account *types.Account, id int) error {
	if err := u.db.WithContext(ctx).
		Model(&types.Account{}).
		Where("user_id = ?", id).
		Updates(account).Error; err != nil {
		return utils.DetermineSQLError(err, "update account")
	}
	return nil
}

func (u *userRepository) UpdateCredentials(ctx context.Context, credentials *types.Credentials, id int) error {
	if err := u.db.WithContext(ctx).
		Model(&types.Credentials{}).
		Where("user_id = ?", id).
		Updates(credentials).Error; err != nil {
		return utils.DetermineSQLError(err, "update credentials")
	}
	return nil
}

func (u *userRepository) Delete(ctx context.Context, id int) error {
	if err := u.db.WithContext(ctx).
		Delete(&types.User{}, id).Error; err != nil {
		return utils.DetermineSQLError(err, "delete user")
	}
	return nil
}

func (u *userRepository) List(ctx context.Context) ([]types.User, error) {
	var users []types.User
	if err := u.db.WithContext(ctx).
		Preload("Profile").
		Preload("Account").
		Preload("Credentials").
		Find(&users).Error; err != nil {
		return nil, utils.DetermineSQLError(err, "list users")
	}
	return users, nil
}

func (u *userRepository) UpdateUsedStorage(ctx context.Context, id int, newUsedStorage int64) error {
	if err := u.db.WithContext(ctx).
		Model(&types.Account{}).
		Where("user_id = ?", id).
		Update("used_storage", newUsedStorage).Error; err != nil {
		return utils.DetermineSQLError(err, "update used storage")
	}
	return nil
}

func (u *userRepository) ReserveStorage(ctx context.Context, userID int, size int64) error {
	err := u.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var user types.User
		if err := tx.
			Clauses(clause.Locking{Strength: "UPDATE"}).
			Preload("Account").
			First(&user, "id = ?", userID).Error; err != nil {
			return utils.DetermineSQLError(err, "locking user for reserve storage")
		}
		acct := user.Account
		if acct.UsedStorage+size > acct.StorageLimit {
			return utils.ErrConflict.New(
				"quota exceeded: used=%d, limit=%d",
				acct.UsedStorage, acct.StorageLimit,
			)
		}
		if err := tx.Model(&types.Account{}).
			Where("user_id = ?", userID).
			UpdateColumn("used_storage", gorm.Expr("used_storage + ?", size)).
			Error; err != nil {
			return utils.DetermineSQLError(err, "reserve storage update")
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func (u *userRepository) ReleaseStorage(ctx context.Context, userID int, size int64) error {
	err := u.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var user types.User
		if err := tx.
			Clauses(clause.Locking{Strength: "UPDATE"}).
			Preload("Account").
			First(&user, "id = ?", userID).Error; err != nil {
			return utils.DetermineSQLError(err, "locking user for release storage")
		}
		acct := user.Account
		newUsed := acct.UsedStorage - size
		if newUsed < 0 {
			newUsed = 0
		}
		if err := tx.Model(&types.Account{}).
			Where("user_id = ?", userID).
			UpdateColumn("used_storage", newUsed).
			Error; err != nil {
			return utils.DetermineSQLError(err, "release storage update")
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}
