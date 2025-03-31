package repositories

import (
	"context"

	"github.com/CustomCloudStorage/types"
	"gorm.io/gorm"
)

type Repository struct {
	Postgres Postgres
}

type Postgres struct {
	Db *gorm.DB
}

type UserRepository interface {
	GetUser(context.Context, string) (*types.User, error)
	GetAllUsers(context.Context) ([]types.User, error)
	CreateUser(context.Context, *types.User) (string, error)
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		Postgres: Postgres{
			Db: db,
		},
	}
}
