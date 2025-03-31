package server

import (
	"file-sharing/config"
	"file-sharing/db"
	"file-sharing/handler"
	"file-sharing/middleware"
	"file-sharing/repository"
	"file-sharing/service"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"time"
)

func Init() error {

	r := gin.New()
	err := r.SetTrustedProxies([]string{"127.0.0.1"})
	if err != nil {
		log.Fatalf("Failed to set trusted proxies %v", err)
	}
	minioClient, err := config.NewMinIOClient()
	if err != nil {
		log.Fatal("Failed to initialize MinIO client:", err)
	}
	linkRepo := repository.NewLinkRepository(db.GetDb())

	// Services
	minioSvc := service.NewMinioService(minioClient, "uploads")
	sharingSvc := service.NewSharingService(linkRepo, minioSvc)

	// Handlers
	shareHandler := handler.NewShareHandler(sharingSvc)
	authHandler := handler.NewAuthHandler()
	go func() {
		ticker := time.NewTicker(1 * time.Hour)
		for range ticker.C {
			linkRepo.CleanExpired()
		}
	}()

	// Auth resource
	ar := r.Group("/api/v1/auth")
	{
		ar.POST("/register", authHandler.Register)
		ar.POST("/login", authHandler.Login)

		// Routes

		ar.Use(middleware.Log(), middleware.Recover())
	}
	// 3. Register Routes (with logging)
	authGroup := r.Group("/")
	authGroup.Use(middleware.Log(),
		middleware.Recover(),
		middleware.Auth())
	{
		authGroup.POST("/upload", handler.UploadHandler(minioClient))
		authGroup.POST("/share/:path", shareHandler.CreateLink)
	}
	r.GET("/share/:path", shareHandler.DownloadShared)

	// Profile resource
	pr := r.Group("/api/v1/profile")
	pr.Use(
		middleware.Log(),
		middleware.Recover(),
		middleware.Auth(),
	)

	pr.GET("/", authHandler.GetProfile)

	if err := r.Run(":" + config.Conf.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
	log.Infof("%s started on port:%s", config.Conf.AppName, config.Conf.DBPort)
	return nil

}
