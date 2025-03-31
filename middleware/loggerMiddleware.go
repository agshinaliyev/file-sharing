package middleware

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func Log() gin.HandlerFunc {
	return func(c *gin.Context) {
		logMsg := "Request: {Method: \"%s\", Endpoint: \"%s\"}"
		log.Infof(logMsg, c.Request.Method, c.Request.URL.Path)

		c.Next()
	}
}
