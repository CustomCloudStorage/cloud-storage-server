package types

type User struct {
	Id          int         `json:"id" gorm:"primaryKey;column:id"`
	Profile     Profile     `json:"profile" gorm:"foreignKey:UserID;references:Id"`
	Account     Account     `json:"account" gorm:"foreignKey:UserID;references:Id"`
	Credentials Credentials `json:"credentials" gorm:"foreignKey:UserID;references:Id"`
}

type Profile struct {
	UserID            int    `json:"user_id" gorm:"primaryKey;column:user_id"`
	Name              string `json:"name" gorm:"column:name"`
	Email             string `json:"email" gorm:"column:email;unique"`
	LastUpdateProfile string `json:"last_update_profile" gorm:"column:last_update_profile"`
}

type Account struct {
	UserID            int    `json:"user_id" gorm:"primaryKey;column:user_id"`
	Role              string `json:"role" gorm:"column:role"`
	StorageLimit      int    `json:"storage_limit" gorm:"column:storage_limit"`
	LastUpdateAccount string `json:"last_update_account" gorm:"column:last_update_account"`
}

type Credentials struct {
	UserID                int    `json:"user_id" gorm:"primaryKey;column:user_id"`
	Password              string `json:"password" gorm:"column:password"`
	LastUpdateCredentials string `json:"last_update_credentials" gorm:"column:last_update_credentials"`
}

type PublicUser struct {
	Id      int     `json:"id" gorm:"primaryKey;column:id"`
	Profile Profile `json:"profile" gorm:"embedded"`
	Account Account `json:"account" gorm:"embedded"`
}

func NewPublicUser(user *User) *PublicUser {
	return &PublicUser{
		Id:      user.Id,
		Profile: user.Profile,
		Account: user.Account,
	}
}

func NewPublicUsers(users []User) []PublicUser {
	publicUsers := make([]PublicUser, len(users))
	for i, user := range users {
		publicUsers[i] = *NewPublicUser(&user)
	}
	return publicUsers
}
