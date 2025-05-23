package main

import (
	"log"
	"net/http"
	"os"
	"text/template"

	"github.com/CustomCloudStorage/config"
	"github.com/CustomCloudStorage/databases"
	"github.com/CustomCloudStorage/handlers"
	"github.com/CustomCloudStorage/infrastructure/email"
	"github.com/CustomCloudStorage/middleware"
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
	redisDB, err := databases.GetRedis(cfg.Redis)
	if err != nil {
		log.Fatalf("Failed to connect to the redis: %v", err)
	}

	userRepo := repositories.NewUserRepository(postgresDB)
	fileRepo := repositories.NewFileRepository(postgresDB)
	folderRepo := repositories.NewFolderRepository(postgresDB)
	uploadSessionRepo := repositories.NewUploadSessionRepository(postgresDB)
	uploadPartRepo := repositories.NewUploadPartRepository(postgresDB)
	trashRepo := repositories.NewTrashRepository(postgresDB)
	authRepo := repositories.NewAuthRepository(postgresDB)
	registerRepo := repositories.NewRegistrationRepository(postgresDB)
	redis := repositories.NewRedisCache(redisDB)

	email := email.NewSMTPMailer(cfg.SMTP)
	var templates = template.New("")
	for _, name := range []string{"registration_confirmation"} {
		content, err := os.ReadFile("templates/" + name + ".tmpl")
		if err != nil {
			log.Fatal(err)
		}
		templates = template.Must(templates.New(name).Parse(string(content)))
	}

	fileService := services.NewFileService(userRepo, fileRepo, folderRepo, cfg.Service)
	folderService := services.NewFolderService(fileRepo, folderRepo, cfg.Service)
	uploadService := services.NewUploadService(userRepo, fileRepo, uploadSessionRepo, uploadPartRepo, cfg.Service)
	trashService := services.NewTrashService(trashRepo, cfg.Service)
	authService := services.NewAuthService(authRepo, redis, cfg.Auth)
	emailService := services.NewEmailService(redis, email, templates)
	registerService := services.NewRegistrationService(registerRepo, userRepo, emailService)

	authMiddleware := middleware.NewAuthMiddleware(authRepo, authService, cfg.Auth)

	userHandler := handlers.NewUserHandler(userRepo, fileRepo, fileService)
	fileHandler := handlers.NewFileHandler(fileRepo, fileService)
	folderHandler := handlers.NewFolderHandler(folderRepo, folderService)
	uploadHandler := handlers.NewUploadHandler(uploadService)
	trashHandler := handlers.NewTrashHandler(trashRepo, trashService)
	authHandler := handlers.NewAuthHandler(authRepo, authService)
	registerhandler := handlers.NewRegistrationHandler(registerService)

	router := mux.NewRouter()

	router.Use(authMiddleware.AuthMiddleWare())
	router.HandleFunc("/auth/register", handlers.HandleError(registerhandler.Register)).Methods("POST")
	router.HandleFunc("/auth/register/confirm", handlers.HandleError(registerhandler.Confirm)).Methods("POST")
	router.HandleFunc("/auth/register/resend", handlers.HandleError(registerhandler.ResendCode)).Methods("POST")
	router.HandleFunc("/auth/login", handlers.HandleError(authHandler.HandleLogIn)).Methods("POST")
	router.HandleFunc("/auth/logout", handlers.HandleError(authHandler.HandleLogOut)).Methods("POST")
	router.HandleFunc("/auth/me", handlers.HandleError(authHandler.HandleAuthMe)).Methods("GET")

	adminRouter := router.PathPrefix("/admin").Subrouter()
	adminRouter.Use(authMiddleware.RequireRole("admin", "superuser"))

	adminRouter.HandleFunc("/users/{id}", handlers.HandleError(userHandler.HandleGetUser)).Methods("GET")
	adminRouter.HandleFunc("/users", handlers.HandleError(userHandler.HandleListUsers)).Methods("GET")
	router.HandleFunc("/me/profile", handlers.HandleError(userHandler.HandleUpdateProfile)).Methods("PUT")
	adminRouter.HandleFunc("/users/{id}/account", handlers.HandleError(userHandler.HandleUpdateAccount)).Methods("PUT")
	router.HandleFunc("/me/credentials", handlers.HandleError(userHandler.HandleUpdateCredentials)).Methods("PUT")
	adminRouter.HandleFunc("/users/{id}", handlers.HandleError(userHandler.HandleDeleteUser)).Methods("DELETE")

	router.HandleFunc("/folders", handlers.HandleError(folderHandler.HandleCreateFolder)).Methods("POST")
	router.HandleFunc("/folders/{folderID}", handlers.HandleError(folderHandler.HandleGetFolder)).Methods("GET")
	router.HandleFunc("/folders/{folderID}", handlers.HandleError(folderHandler.HandleUpdateFolder)).Methods("PUT")
	router.HandleFunc("/folders", handlers.HandleError(folderHandler.HandleListFolders)).Methods("GET")
	router.HandleFunc("/folders/{folderID}/download", handlers.HandleError(folderHandler.DownloadFolderHandler)).Methods("GET")

	router.HandleFunc("/uploads/init", handlers.HandleError(uploadHandler.InitSessionHandler)).Methods("POST")
	router.HandleFunc("/uploads/{sessionID}/{partNumber}", handlers.HandleError(uploadHandler.UploadPartHandler)).Methods("PUT")
	router.HandleFunc("/uploads/{sessionID}/progress", handlers.HandleError(uploadHandler.ProgressHandler)).Methods("GET")
	router.HandleFunc("/uploads/{sessionID}/complete", handlers.HandleError(uploadHandler.CompleteHandler)).Methods("POST")
	router.HandleFunc("/uploads/{sessionID}", handlers.HandleError(uploadHandler.AbortHandler)).Methods("DELETE")

	router.HandleFunc("/files/{fileID}", handlers.HandleError(fileHandler.HandleGetFile)).Methods("GET")
	router.HandleFunc("/files", handlers.HandleError(fileHandler.HandleListFiles)).Methods("GET")
	router.HandleFunc("/files/{fileID}/name", handlers.HandleError(fileHandler.HandleUpdateName)).Methods("PUT")
	router.HandleFunc("/files/{fileID}/folderID", handlers.HandleError(fileHandler.HandleUpdateFolderID)).Methods("PUT")
	router.HandleFunc("/files/{fileID}/download-url", handlers.HandleError(fileHandler.DownloadURLHandler)).Methods("GET")
	router.HandleFunc("/files/download", handlers.HandleError(fileHandler.DownloadByTokenHandler)).Methods("GET")

	router.HandleFunc("/trash/files", handlers.HandleError(trashHandler.ListFilesHandler)).Methods("GET")
	router.HandleFunc("/trash/files/{fileID}", handlers.HandleError(trashHandler.DeleteFileHandler)).Methods("DELETE")
	router.HandleFunc("/trash/files/{fileID}/restore", handlers.HandleError(trashHandler.RestoreFileHandler)).Methods("POST")
	router.HandleFunc("/trash/files/{fileID}/permanent", handlers.HandleError(trashHandler.PermanentDeleteFileHandler)).Methods("DELETE")
	router.HandleFunc("/trash/folders", handlers.HandleError(trashHandler.ListFoldersHandler)).Methods("GET")
	router.HandleFunc("/trash/folders/{folderID}", handlers.HandleError(trashHandler.DeleteFolderHandler)).Methods("DELETE")
	router.HandleFunc("/trash/folders/{folderID}/restore", handlers.HandleError(trashHandler.RestoreFolderHandler)).Methods("POST")
	router.HandleFunc("/trash/files/{fileID}/permanent", handlers.HandleError(trashHandler.PermanentDeleteFileHandler)).Methods("DELETE")

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
