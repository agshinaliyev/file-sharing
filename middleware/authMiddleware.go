package middleware

import (
	"errors"
	"file-sharing/jwt"
	errModel "file-sharing/model/error"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {

		tokenHeader := c.GetHeader("Authorization")

		userId, err := jwt.JwtParse(tokenHeader)

		if err != nil {

			if errors.Is(err, &errModel.InvalidJWTToken) {
				c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
				c.Abort()
				return
			} else if errors.Is(err, &errModel.InvalidJWTToken) {
				c.JSON(http.StatusUnauthorized, gin.H{"message": "Token is invalid"})
				c.Abort()
				return
			} else if errors.Is(err, &errModel.InvalidJWTToken) {
				c.JSON(http.StatusUnauthorized, gin.H{"message": "Token is expired"})
				c.Abort()
				return
			} else {
				c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
				c.Abort()
				return
			}

		}

		c.Set("USER_ID", userId)
		c.Next()

	}
}
