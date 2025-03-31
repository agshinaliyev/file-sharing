package handler

import (
	"file-sharing/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type ShareHandler struct {
	sharingSvc service.SharingService
}

func NewShareHandler(sharingSvc service.SharingService) *ShareHandler {
	return &ShareHandler{sharingSvc: sharingSvc}
}

func (h *ShareHandler) CreateLink(c *gin.Context) {
	filePath := c.Param("path")
	url, err := h.sharingSvc.CreateShareLink(filePath, 24*time.Hour)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"url": url})
}

func (h *ShareHandler) DownloadShared(c *gin.Context) {
	filePath := c.Param("path")
	token := c.Query("token")

	url, err := h.sharingSvc.ValidateAndGetURL(filePath, token)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	c.Redirect(http.StatusFound, url)
}
