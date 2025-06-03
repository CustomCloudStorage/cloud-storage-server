package types

import "time"

type User struct {
	Id          int         `json:"id" gorm:"primaryKey;column:id"`
	Profile     Profile     `json:"profile" gorm:"foreignKey:UserID;references:Id"`
	Account     Account     `json:"account" gorm:"foreignKey:UserID;references:Id"`
	Credentials Credentials `json:"credentials" gorm:"foreignKey:UserID;references:Id"`
	CreatedAt   time.Time   `json:"created_at" gorm:"column:created_at;autoCreateTime"`
}

type Profile struct {
	UserID    int       `json:"user_id" gorm:"primaryKey;column:user_id"`
	Name      string    `json:"name" gorm:"column:name"`
	Email     string    `json:"email" gorm:"column:email;unique"`
	UpdatedAt time.Time `json:"updated_at" gorm:"column:updated_at;autoUpdateTime"`
}

type Account struct {
	UserID       int       `json:"user_id" gorm:"primaryKey;column:user_id"`
	Role         string    `json:"role" gorm:"column:role"`
	StorageLimit int64     `json:"storage_limit" gorm:"column:storage_limit"type:bigint`
	UsedStorage  int64     `json:"used_storage" gorm:"not null;column:used_storage"type:bigint`
	UpdatedAt    time.Time `json:"updated_at" gorm:"column:updated_at;autoUpdateTime"`
}

type Credentials struct {
	UserID    int       `json:"user_id" gorm:"primaryKey;column:user_id"`
	Password  string    `json:"password" gorm:"column:password"`
	UpdatedAt time.Time `json:"updated_at" gorm:"column:updated_at;autoUpdateTime"`
}

type PublicUser struct {
	Id        int       `json:"id"`
	Profile   Profile   `json:"profile"`
	Account   Account   `json:"account"`
	CreatedAt time.Time `json:"created_at"`
}

type StorageStats struct {
	TotalBytes     int64 `json:"total_bytes"`
	AllocatedBytes int64 `json:"allocated_bytes"`
	FreeBytes      int64 `json:"free_bytes"`
}

func NewPublicUser(user *User) *PublicUser {
	return &PublicUser{
		Id:        user.Id,
		Profile:   user.Profile,
		Account:   user.Account,
		CreatedAt: user.CreatedAt,
	}
}

func NewPublicUsers(users []User) []PublicUser {
	publicUsers := make([]PublicUser, len(users))
	for i, user := range users {
		publicUsers[i] = *NewPublicUser(&user)
	}
	return publicUsers
}
