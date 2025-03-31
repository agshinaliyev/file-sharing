package middleware

import (
	errModel "file-sharing/model/error"
	"github.com/gin-gonic/gin"

	log "github.com/sirupsen/logrus"
)

func Recover() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				log.Errorf("Panic recovered: %v\n", err)
				c.JSON(errModel.UnexpectedError.Code, gin.H{"error": errModel.UnexpectedError.Message})
				c.Abort()
			}
		}()

		c.Next()
	}
}
