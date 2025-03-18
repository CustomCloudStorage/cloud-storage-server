package types

import "time"

type User struct {
	Id           string    `json:"id"`
	Name         string    `json:"name"`
	Email        string    `json:"email"`
	Password     string    `json:"password"`
	Role         string    `json:"role"`
	StorageLimit int       `json:"storage-limit"`
	LastUpdate   time.Time `json:"last-update"`
}
