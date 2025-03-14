package main

import (
	"log"
	"net/http"

	"github.com/CustomCloudStorage/config"
	"github.com/CustomCloudStorage/databases"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
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

	router := mux.NewRouter()

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{cfg.Cors.AllowedOrigin},
		AllowCredentials: true,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Authorization", "Content-Type", "Token"},
		Debug:            false,
	})

	h := c.Handler(router)

	log.Println("Server is up and running!")
	log.Fatal(http.ListenAndServe(":"+cfg.Port, h))
}
