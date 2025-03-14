package main

import (
	"github.com/CustomCloudStorage/config"
	"github.com/CustomCloudStorage/databases"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		return
	}

	postgresDB, err := databases.GetDB(cfg.Postgres)
	if err != nil {
		return
	}
	defer postgresDB.Close()
}
