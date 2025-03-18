package databases

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/CustomCloudStorage/utils"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

type PostgresConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
}

func GetDB(cfg PostgresConfig) (*sql.DB, error) {
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, utils.ErrConnection.Wrap(err, "failed to open database connection")
	}

	if err = db.Ping(); err != nil {
		return nil, utils.ErrPingFailed.Wrap(err, "failed to ping database")
	}

	if err := RunMigrations(db, "./migrations"); err != nil {
		log.Fatalf("Could not run migrations: %v", err)
	}

	return db, nil
}

func RunMigrations(db *sql.DB, migrationsPath string) error {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return utils.ErrDriverCreate.Wrap(err, "failed to create migration driver")
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://"+migrationsPath,
		"postgres", driver)
	if err != nil {
		return utils.ErrMigration.Wrap(err, "failed to create migration instance")
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return utils.ErrMigration.Wrap(err, "failed to run migrations")
	}

	log.Println("Migrations applied successfully")
	return nil
}
