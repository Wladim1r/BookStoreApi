package middleware

import (
	"bookstore-api/internal/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {

		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Authorization field empty",
			})
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Authorization header format must be 'Bearer <token>'",
			})
			return
		}

		token, err := utils.ParseToken(parts[1])
		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid token",
			})
			return
		}

		claims := token.Claims.(jwt.MapClaims)

		c.Set("userID", claims["user_id"])
		c.Next()
	}
}

func AdminAuth() gin.HandlerFunc {
	auth := gin.BasicAuth(gin.Accounts{
		"SuperUser": "qwerty12345",
	})
	return func(c *gin.Context) {
		auth(c)
		if c.IsAborted() {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "check username or password",
			})
			return
		}

		c.Next()
	}
}
