package main

import (
	"log"
	"net/http"

	"github.com/CustomCloudStorage/config"
	"github.com/CustomCloudStorage/databases"
	"github.com/CustomCloudStorage/handlers"
	"github.com/CustomCloudStorage/repositories"
	"github.com/CustomCloudStorage/services"
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

	userRepo := repositories.NewUserRepository(postgresDB)
	fileRepo := repositories.NewFileRepository(postgresDB)
	folderRepo := repositories.NewFolderRepository(postgresDB)
	uploadSessionRepo := repositories.NewUploadSessionRepository(postgresDB)
	uploadPartRepo := repositories.NewUploadPartRepository(postgresDB)

	fileService := services.NewFileService(userRepo, fileRepo, folderRepo, cfg.Service)
	folderService := services.NewFolderService(fileRepo, folderRepo, cfg.Service)
	uploadService := services.NewUploadService(userRepo, fileRepo, uploadSessionRepo, uploadPartRepo, cfg.Service)

	userHandler := handlers.NewUserHandler(userRepo, fileRepo, fileService)
	fileHandler := handlers.NewFileHandler(fileRepo, fileService)
	folderHandler := handlers.NewFolderHandler(folderRepo, folderService)
	uploadHandler := handlers.NewUploadHandler(uploadService)

	router := mux.NewRouter()

	router.HandleFunc("/users/{id}", handlers.HandleError(userHandler.HandleGetUser)).Methods("GET")
	router.HandleFunc("/users", handlers.HandleError(userHandler.HandleListUsers)).Methods("GET")
	router.HandleFunc("/users", handlers.HandleError(userHandler.HandleCreateUser)).Methods("POST")
	router.HandleFunc("/users/{id}/profile", handlers.HandleError(userHandler.HandleUpdateProfile)).Methods("PUT")
	router.HandleFunc("/users/{id}/account", handlers.HandleError(userHandler.HandleUpdateAccount)).Methods("PUT")
	router.HandleFunc("/users/{id}/credentials", handlers.HandleError(userHandler.HandleUpdateCredentials)).Methods("PUT")
	router.HandleFunc("/users/{id}", handlers.HandleError(userHandler.HandleDeleteUser)).Methods("DELETE")

	router.HandleFunc("/users/{id}/folders", handlers.HandleError(folderHandler.HandleCreateFolder)).Methods("POST")
	router.HandleFunc("/users/{id}/folders/{folderID}", handlers.HandleError(folderHandler.HandleGetFolder)).Methods("GET")
	router.HandleFunc("/users/{id}/folders/{folderID}", handlers.HandleError(folderHandler.HandleUpdateFolder)).Methods("PUT")
	router.HandleFunc("/users/{id}/folders", handlers.HandleError(folderHandler.HandleListFolders)).Methods("GET")
	router.HandleFunc("/users/{userID}/folders/{folderID}/download", handlers.HandleError(folderHandler.DownloadFolderHandler)).Methods("GET")

	router.HandleFunc("/uploads/init", handlers.HandleError(uploadHandler.InitSessionHandler)).Methods("POST")
	router.HandleFunc("/uploads/{sessionID}/{partNumber}", handlers.HandleError(uploadHandler.UploadPartHandler)).Methods("PUT")
	router.HandleFunc("/uploads/{sessionID}/progress", handlers.HandleError(uploadHandler.ProgressHandler)).Methods("GET")
	router.HandleFunc("/uploads/{sessionID}/complete", handlers.HandleError(uploadHandler.CompleteHandler)).Methods("POST")
	router.HandleFunc("/uploads/{sessionID}", handlers.HandleError(uploadHandler.AbortHandler)).Methods("DELETE")

	router.HandleFunc("/users/{id}/files/{fileID}", handlers.HandleError(fileHandler.HandleGetFile)).Methods("GET")
	router.HandleFunc("/users/{id}/files", handlers.HandleError(fileHandler.HandleListFiles)).Methods("GET")
	router.HandleFunc("/users/{id}/files/{fileID}/name", handlers.HandleError(fileHandler.HandleUpdateName)).Methods("PUT")
	router.HandleFunc("/users/{id}/files/{fileID}/folderID", handlers.HandleError(fileHandler.HandleUpdateFolderID)).Methods("PUT")
	router.HandleFunc("/users/{userID}/files/{fileID}/download-url", handlers.HandleError(fileHandler.DownloadURLHandler)).Methods("GET")
	router.HandleFunc("/files/download", handlers.HandleError(fileHandler.DownloadByTokenHandler)).Methods("GET")

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
