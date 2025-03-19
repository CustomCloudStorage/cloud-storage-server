package repositories

import (
	"context"
	"database/sql"

	"github.com/CustomCloudStorage/types"
)

type Repository struct {
	Postgres Postgres
}

type Postgres struct {
	Db *sql.DB
}

type UserRepository interface {
	GetUser(context.Context, string) (*types.User, error)
	GetAllUsers(context.Context) ([]types.User, error)
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		Postgres: Postgres{
			Db: db,
		},
	}
}
