package main

import (
	"log"
	"net/http"

	"github.com/CustomCloudStorage/config"
	"github.com/CustomCloudStorage/databases"
	"github.com/CustomCloudStorage/handlers"
	"github.com/CustomCloudStorage/repositories"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load the config: %v", err)
	}

	postgresDB, err := databases.GetDB(cfg.Postgres)
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}

	repository := repositories.NewRepository(postgresDB)
	handler := handlers.NewHandler(repository)

	router := mux.NewRouter()

	router.HandleFunc("/users/{id}", handlers.HandleError(handler.HandleGetUser)).Methods("GET")
	router.HandleFunc("/users", handlers.HandleError(handler.HandleGetAllUsers)).Methods("GET")
	router.HandleFunc("/users", handlers.HandleError(handler.HandleCreateUser)).Methods("POST")
	router.HandleFunc("/users/{id}/profile", handlers.HandleError(handler.HandleUpdateProfile)).Methods("PUT")
	router.HandleFunc("/users/{id}/account", handlers.HandleError(handler.HandleUpdateAccount)).Methods("PUT")

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
