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

	InitSuperuser(postgresDB, cfg.Superuser)

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
	trashService := services.NewTrashService(trashRepo, fileService)
	authService := services.NewAuthService(authRepo, redis, cfg.Auth)
	emailService := services.NewEmailService(redis, email, templates)
	registerService := services.NewRegistrationService(registerRepo, userRepo, emailService, cfg.Service)
	userService := services.NewUserService(userRepo, cfg.Service)

	authMiddleware := middleware.NewAuthMiddleware(authRepo, authService, cfg.Auth)

	userHandler := handlers.NewUserHandler(userRepo, fileRepo, fileService, userService)
	fileHandler := handlers.NewFileHandler(fileRepo, fileService)
	folderHandler := handlers.NewFolderHandler(folderRepo, folderService)
	uploadHandler := handlers.NewUploadHandler(uploadService)
	trashHandler := handlers.NewTrashHandler(trashRepo, trashService)
	authHandler := handlers.NewAuthHandler(authRepo, authService)
	registerhandler := handlers.NewRegistrationHandler(registerRepo, registerService)

	router := mux.NewRouter()

	router.Use(authMiddleware.AuthMiddleWare())
	router.HandleFunc("/auth/register", middleware.HandleError(registerhandler.Register)).Methods("POST")
	router.HandleFunc("/auth/register/confirm", middleware.HandleError(registerhandler.Confirm)).Methods("POST")
	router.HandleFunc("/auth/register/resend", middleware.HandleError(registerhandler.ResendCode)).Methods("POST")
	router.HandleFunc("/auth/login", middleware.HandleError(authHandler.HandleLogIn)).Methods("POST")
	router.HandleFunc("/auth/logout", middleware.HandleError(authHandler.HandleLogOut)).Methods("POST")
	router.HandleFunc("/auth/me", middleware.HandleError(authHandler.HandleAuthMe)).Methods("GET")

	adminRouter := router.PathPrefix("/admin").Subrouter()
	adminRouter.Use(authMiddleware.RequireRole("admin", "superuser"))

	adminRouter.HandleFunc("/users/{id}", middleware.HandleError(userHandler.HandleGetUser)).Methods("GET")
	adminRouter.HandleFunc("/users", middleware.HandleError(userHandler.HandleListUsers)).Methods("GET")
	router.HandleFunc("/me/profile", middleware.HandleError(userHandler.HandleUpdateProfile)).Methods("PUT")
	adminRouter.HandleFunc("/users/{id}/account", middleware.HandleError(userHandler.HandleUpdateAccount)).Methods("PUT")
	router.HandleFunc("/me/credentials", middleware.HandleError(userHandler.HandleUpdateCredentials)).Methods("PUT")
	adminRouter.HandleFunc("/users/{id}", middleware.HandleError(userHandler.HandleDeleteUser)).Methods("DELETE")
	adminRouter.HandleFunc("/users/storage", middleware.HandleError(userHandler.HandleStorageStats)).Methods("GET")

	router.HandleFunc("/folders", middleware.HandleError(folderHandler.HandleCreateFolder)).Methods("POST")
	router.HandleFunc("/folders/{folderID}", middleware.HandleError(folderHandler.HandleGetFolder)).Methods("GET")
	router.HandleFunc("/folders/{folderID}", middleware.HandleError(folderHandler.HandleUpdateFolder)).Methods("PUT")
	router.HandleFunc("/folders", middleware.HandleError(folderHandler.HandleListFolders)).Methods("GET")
	router.HandleFunc("/folders/{folderID}/download", middleware.HandleError(folderHandler.DownloadFolderHandler)).Methods("GET")

	router.HandleFunc("/uploads/init", middleware.HandleError(uploadHandler.InitSessionHandler)).Methods("POST")
	router.HandleFunc("/uploads/{sessionID}/{partNumber}", middleware.HandleError(uploadHandler.UploadPartHandler)).Methods("PUT")
	router.HandleFunc("/uploads/{sessionID}/progress", middleware.HandleError(uploadHandler.ProgressHandler)).Methods("GET")
	router.HandleFunc("/uploads/{sessionID}/complete", middleware.HandleError(uploadHandler.CompleteHandler)).Methods("POST")
	router.HandleFunc("/uploads/{sessionID}", middleware.HandleError(uploadHandler.AbortHandler)).Methods("DELETE")

	router.HandleFunc("/files/{fileID}", middleware.HandleError(fileHandler.HandleGetFile)).Methods("GET")
	router.HandleFunc("/files", middleware.HandleError(fileHandler.HandleListFiles)).Methods("GET")
	router.HandleFunc("/files/{fileID}/name", middleware.HandleError(fileHandler.HandleUpdateName)).Methods("PUT")
	router.HandleFunc("/files/{fileID}/folderID", middleware.HandleError(fileHandler.HandleUpdateFolderID)).Methods("PUT")
	router.HandleFunc("/files/{fileID}/download-url", middleware.HandleError(fileHandler.DownloadURLHandler)).Methods("GET")
	router.HandleFunc("/files/download", middleware.HandleError(fileHandler.DownloadByTokenHandler)).Methods("GET")

	router.HandleFunc("/trash/files", middleware.HandleError(trashHandler.ListFilesHandler)).Methods("GET")
	router.HandleFunc("/trash/files/{fileID}", middleware.HandleError(trashHandler.DeleteFileHandler)).Methods("DELETE")
	router.HandleFunc("/trash/files/{fileID}/restore", middleware.HandleError(trashHandler.RestoreFileHandler)).Methods("POST")
	router.HandleFunc("/trash/files/{fileID}/permanent", middleware.HandleError(trashHandler.PermanentDeleteFileHandler)).Methods("DELETE")
	router.HandleFunc("/trash/folders", middleware.HandleError(trashHandler.ListFoldersHandler)).Methods("GET")
	router.HandleFunc("/trash/folders/{folderID}", middleware.HandleError(trashHandler.DeleteFolderHandler)).Methods("DELETE")
	router.HandleFunc("/trash/folders/{folderID}/restore", middleware.HandleError(trashHandler.RestoreFolderHandler)).Methods("POST")
	router.HandleFunc("/trash/files/{fileID}/permanent", middleware.HandleError(trashHandler.PermanentDeleteFileHandler)).Methods("DELETE")

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{cfg.Cors.AllowedOrigin},
		AllowCredentials: true,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders:   []string{"Authorization", "Content-Type", "Token"},
		Debug:            false,
	})

	h := c.Handler(router)

	log.Println("Server is up and running!")
	log.Fatal(http.ListenAndServe(":"+cfg.Port, h))
}
