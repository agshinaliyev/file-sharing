package handler

import (
	jwt2 "file-sharing/jwt"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/minio/minio-go/v7"
	"io"
	"net/http"
	"net/url"
	"strings"
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
			"path":    userPrefix + file.Filename, // Return file path
			"view":    fmt.Sprintf("http://localhost:8080/view?object=%s", url.QueryEscape(userPrefix+file.Filename)),
		})
	}
}
func ViewFileHandler(minioClient *minio.Client) gin.HandlerFunc {
	return func(c *gin.Context) {

		objectName := c.Query("object")
		if objectName == "" {
			c.String(http.StatusBadRequest, "Missing 'object' query param")
			return
		}
		shareToken := c.Query("token")
		if shareToken != "" {
			token, err := jwt.Parse(shareToken, func(token *jwt.Token) (interface{}, error) {
				return jwt2.JwtShared, nil
			})
			if err != nil || !token.Valid {
				c.String(http.StatusUnauthorized, "Invalid or expired share token")
				return
			}

			claims := token.Claims.(jwt.MapClaims)
			if claims["object"] != objectName {
				c.String(http.StatusForbidden, "Token does not match requested object")
				return
			}
		} else {
			// Fallback to requiring Authorization token
			authHeader := c.GetHeader("Authorization")
			if !strings.HasPrefix(authHeader, "Bearer ") {
				c.String(http.StatusUnauthorized, "Missing or invalid token")
				return
			}

			// You can verify login token here if needed
			// or skip if already handled by middleware
		}

		object, err := minioClient.GetObject(
			c.Request.Context(),
			"uploads", // your bucket name
			objectName,
			minio.GetObjectOptions{},
		)
		if err != nil {
			c.String(http.StatusInternalServerError, "Failed to get object: %v", err)
			return
		}

		// Get metadata to set the proper content type
		stat, err := object.Stat()
		if err != nil {
			c.String(http.StatusInternalServerError, "Failed to stat object: %v", err)
			return
		}

		// Set headers to match the file's content
		c.Header("Content-Type", stat.ContentType)
		c.Header("Content-Length", fmt.Sprintf("%d", stat.Size))
		c.Header("Content-Disposition", "inline") // Show in browser if possible

		// Stream it directly to the client
		_, err = io.Copy(c.Writer, object)
		if err != nil {
			c.String(http.StatusInternalServerError, "Failed to stream object: %v", err)
			return
		}
	}
}
func ShareFileHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		object := c.Query("object")
		if object == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Missing object param"})
			return
		}

		token, err := jwt2.SharedToken(object, 1*time.Minute) // valid for 10 min
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate share token"})
			return
		}

		link := fmt.Sprintf("http://localhost:8080/view?object=%s&token=%s", object, token)
		c.JSON(http.StatusOK, gin.H{"share_link": link})
	}
}

//func GenerateTempURL(minioClient *minio.Client) gin.HandlerFunc {
//	return func(c *gin.Context) {
//		// Get filename from URL params
//		//username := c.MustGet("username").(string)
//
//		filename := c.Param("filename")
//
//		// Set expiry time (e.g., 24 hours)
//		//expiryHours := 10
//		expiry := 5 * time.Second
//
//		// Generate presigned URL
//		url, err := minioClient.PresignedGetObject(
//			c.Request.Context(),
//			"uploads", // Bucket name
//			filename,  // File name
//			expiry,    // Link validity duration
//			nil,       // Optional request parameters
//		)
//		if err != nil {
//			c.JSON(http.StatusInternalServerError, gin.H{
//				"error": "Failed to generate temporary URL",
//			})
//			return
//		}
//
//		// Return the temporary URL
//		c.JSON(http.StatusOK, gin.H{
//			"url":    url.String(),
//			"expiry": expiry,
//		})
//	}
//}
