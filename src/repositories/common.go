package repositories

import "database/sql"

type Repository struct {
	Postgres Postgres
}

type Postgres struct {
	Db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		Postgres: Postgres{
			Db: db,
		},
	}
}
