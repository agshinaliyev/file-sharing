package server

import (
	"file-sharing/config"
	"file-sharing/handler"
	"file-sharing/middleware"
	"file-sharing/repository"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func Init() error {

	r := gin.New()
	err := r.SetTrustedProxies([]string{"127.0.0.1"})
	if err != nil {
		log.Fatalf("Failed to set trusted proxies %v", err)
	}
	minioCfg := config.LoadMinioConfig()
	minioRepo, err := repository.NewMinioRepository(minioCfg)
	if err != nil {
		log.Fatal("MinIO initialization failed:", err)
	}

	// Add this where you register routes
	storageHandler := handler.NewStorageHandler(minioRepo)
	a := r.Group("/api/v1/auth")
	a.POST("/upload", storageHandler.UploadFile)
	authHandler := handler.NewAuthHandler()

	// Auth resource
	ar := r.Group("/api/v1/auth")
	{
		ar.POST("/register", authHandler.Register)
		ar.POST("/login", authHandler.Login)

		// Routes

		ar.Use(middleware.Log(), middleware.Recover())
	}
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
