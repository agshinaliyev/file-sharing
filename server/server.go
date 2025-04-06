package server

import (
	"file-sharing/config"
	"file-sharing/handler"
	"file-sharing/middleware"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
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

	// Services

	// Handlers
	authHandler := handler.NewAuthHandler()

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
		authGroup.GET("/view", handler.ViewFileHandler(minioClient))
		authGroup.GET("/share", handler.ShareFileHandler())
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
