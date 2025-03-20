package types

type User struct {
	Id           string `json:"id"`
	Name         string `json:"name"`
	Email        string `json:"email"`
	Password     string `json:"password"`
	Role         string `json:"role"`
	StorageLimit int    `json:"storage_limit"`
	LastUpdate   string `json:"last_update"`
}
