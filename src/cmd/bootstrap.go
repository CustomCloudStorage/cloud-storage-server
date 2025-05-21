package main

import (
	"log"

	"github.com/CustomCloudStorage/config"
	"github.com/CustomCloudStorage/types"
	"github.com/CustomCloudStorage/utils"
	"gorm.io/gorm"
)

func InitSuperuser(db *gorm.DB, suCfg config.SuperuserConfig) {
	var prof types.Profile
	err := db.Where("email = ?", suCfg.Email).First(&prof).Error

	switch {
	case err == nil:
		log.Printf("Superuser already exists (user_id=%d)", prof.UserID)
		return

	case err == gorm.ErrRecordNotFound:
		hashed, err := utils.HashPassword(suCfg.Password)
		if err != nil {
			log.Fatalf("bcrypt.GenerateFromPassword: %v", err)
		}

		user := types.User{
			Profile: types.Profile{
				Name:  "Superuser",
				Email: suCfg.Email,
			},
			Credentials: types.Credentials{
				Password: string(hashed),
			},
			Account: types.Account{
				Role:         "superuser",
				StorageLimit: 0,
				UsedStorage:  0,
			},
		}
		if err := db.Create(&user).Error; err != nil {
			log.Fatalf("bootstrap superuser: cannot create user: %v", err)
		}
		log.Printf("→ Superuser created: id=%d, email=%s", user.Id, suCfg.Email)
		return

	default:
		log.Fatalf("bootstrap superuser: DB error: %v", err)
	}
}
