package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v7"
	"net/http"
	"time"
)

func UploadHandler(minioClient *minio.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		username := c.MustGet("username").(string)
		file, err := c.FormFile("file")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "No file uploaded"})
			return
		}

		// Open the file
		fileReader, err := file.Open()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read file"})
			return
		}
		defer fileReader.Close()

		// Upload to MinIO bucket "uploads"
		userPrefix := "user:" + username + "/"
		_, err = minioClient.PutObject(
			c.Request.Context(),
			"uploads",                // Bucket name
			userPrefix+file.Filename, // Use base name for security
			fileReader,
			file.Size,
			minio.PutObjectOptions{
				ContentType: file.Header.Get("Content-Type"),
			},
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload to MinIO"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": fmt.Sprintf("'%s' uploaded successfully!", file.Filename),
			"path":    fmt.Sprintf("/uploads/%s", file.Filename), // Return file path
		})
	}
}
func GenerateTempURL(minioClient *minio.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get filename from URL params
		//username := c.MustGet("username").(string)

		filename := c.Param("filename")

		// Set expiry time (e.g., 24 hours)
		//expiryHours := 10
		expiry := 5 * time.Second

		// Generate presigned URL
		url, err := minioClient.PresignedGetObject(
			c.Request.Context(),
			"uploads", // Bucket name
			filename,  // File name
			expiry,    // Link validity duration
			nil,       // Optional request parameters
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to generate temporary URL",
			})
			return
		}

		// Return the temporary URL
		c.JSON(http.StatusOK, gin.H{
			"url":    url.String(),
			"expiry": expiry,
		})
	}
}
