package databases

import (
	"database/sql"
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/CustomCloudStorage/utils"
	migrate "github.com/golang-migrate/migrate/v4"
	migratepg "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

type PostgresConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
}

func GetDB(cfg PostgresConfig) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Error),
	})
	if err != nil {
		return nil, utils.ErrConnection.Wrap(err, "failed to open database connection")
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, utils.ErrPingFailed.Wrap(err, "failed to ping database")
	}

	if err = sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	if err := RunMigrations(sqlDB, "./migrations"); err != nil {
		log.Fatalf("Could not run migrations: %v", err)
	}

	return db, nil
}

func RunMigrations(db *sql.DB, migrationsPath string) error {
	driver, err := migratepg.WithInstance(db, &migratepg.Config{})
	if err != nil {
		return utils.ErrDriverCreate.Wrap(err, "failed to create migration driver")
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://"+migrationsPath,
		"postgres",
		driver,
	)
	if err != nil {
		return utils.ErrMigration.Wrap(err, "failed to create migration instance")
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return utils.ErrMigration.Wrap(err, "failed to run migrations")
	}

	ver, dirty, err := GetCurrentMigrationVersion(db, "./migrations")
	if err != nil {
		return err
	}
	log.Printf("current version: %d, dirty: %v", ver, dirty)

	return nil
}

func GetCurrentMigrationVersion(db *sql.DB, migrationsPath string) (version uint, dirty bool, err error) {
	driver, err := migratepg.WithInstance(db, &migratepg.Config{})
	if err != nil {
		return 0, false, utils.ErrDriverCreate.Wrap(err, "failed to create migration driver")
	}
	m, err := migrate.NewWithDatabaseInstance(
		"file://"+migrationsPath,
		"postgres",
		driver,
	)
	if err != nil {
		return 0, false, utils.ErrMigration.Wrap(err, "failed to create migration instance")
	}
	version, dirty, err = m.Version()
	if err == migrate.ErrNilVersion {
		return 0, false, nil
	}
	return version, dirty, err
}
