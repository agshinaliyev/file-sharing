package handler

import (
	"file-sharing/repository"
	"github.com/gin-gonic/gin"
	"net/http"
	"path/filepath"
)

type StorageHandler struct {
	repo *repository.MinioRepository
}

func NewStorageHandler(repo *repository.MinioRepository) *StorageHandler {
	return &StorageHandler{repo: repo}
}

// Only HTTP handling (no MinIO logic!)
func (h *StorageHandler) UploadFile(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file uploaded"})
		return
	}

	tempPath := filepath.Join("/tmp", file.Filename)
	if err := c.SaveUploadedFile(file, tempPath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
		return
	}

	if err := h.repo.UploadFile(c.Request.Context(), file.Filename, tempPath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Upload failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "File uploaded"})
}
