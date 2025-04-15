package main

import (
	"log"
	"net/http"

	"github.com/CustomCloudStorage/config"
	"github.com/CustomCloudStorage/databases"
	"github.com/CustomCloudStorage/handlers"
	"github.com/CustomCloudStorage/repositories"
	"github.com/CustomCloudStorage/service"
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
	service := service.NewService(repository, cfg.StorageDir)
	handler := handlers.NewHandler(repository, service)

	router := mux.NewRouter()

	router.HandleFunc("/users/{id}", handlers.HandleError(handler.HandleGetUser)).Methods("GET")
	router.HandleFunc("/users", handlers.HandleError(handler.HandleListUsers)).Methods("GET")
	router.HandleFunc("/users", handlers.HandleError(handler.HandleCreateUser)).Methods("POST")
	router.HandleFunc("/users/{id}/profile", handlers.HandleError(handler.HandleUpdateProfile)).Methods("PUT")
	router.HandleFunc("/users/{id}/account", handlers.HandleError(handler.HandleUpdateAccount)).Methods("PUT")
	router.HandleFunc("/users/{id}/credentials", handlers.HandleError(handler.HandleUpdateCredentials)).Methods("PUT")
	router.HandleFunc("/users/{id}", handlers.HandleError(handler.HandleDeleteUser)).Methods("DELETE")

	router.HandleFunc("/users/{id}/folders", handlers.HandleError(handler.HandleCreateFolder)).Methods("POST")
	router.HandleFunc("/users/{id}/folders/{folderID}", handlers.HandleError(handler.HandleGetFolder)).Methods("GET")
	router.HandleFunc("/users/{id}/folders/{folderID}", handlers.HandleError(handler.HandleUpdateFolder)).Methods("PUT")
	router.HandleFunc("/users/{id}/folders/{folderID}", handlers.HandleError(handler.HandleDeleteFolder)).Methods("DELETE")
	router.HandleFunc("/users/{id}/folders", handlers.HandleError(handler.HandleDeleteFolder)).Methods("GET")

	router.HandleFunc("/users/{id}/files", handlers.HandleError(handler.HandleUploadFile)).Methods("POST")
	router.HandleFunc("/users/{id}/files/{fileID}", handlers.HandleError(handler.HandleGetFile)).Methods("POST")
	router.HandleFunc("/users/{id}/files/{fileID}", handlers.HandleError(handler.HandleDeleteFile)).Methods("DELETE")
	router.HandleFunc("/users/{id}/files/{fileID}", handlers.HandleError(handler.HandleListFiles))
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
