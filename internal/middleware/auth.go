package middleware

import (
	"bookstore-api/internal/models"
	"bookstore-api/internal/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// @Summary JWT Authentication
// @Description Protects endpoints through JWT token
// @Security ApiKeyAuth
// @Param Authorization header string true "JWT Token" default(Bearer <token>)
// @Failure 400 {object} models.ErrorResponse "Invalid token"
// @Failure 401 {object} models.ErrorResponse "Required authorization"
func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {

		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, models.ErrorResponse{
				Error: "Authorization field empty",
			})
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, models.ErrorResponse{
				Error: "Authorization header format must be 'Bearer <token>'",
			})
			return
		}

		token, err := utils.ParseToken(parts[1])
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, models.ErrorResponse{
				Error: "Invalid token",
			})
			return
		}

		if !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, models.ErrorResponse{
				Error: "Token is expired or invalid",
			})
			return
		}

		claims := token.Claims.(jwt.MapClaims)

		c.Set("userID", claims["user_id"])
		c.Next()
	}
}

// @Summary Basic Authentication for Admin
// @Description Protects endpoints through Basic Auth
// @Security BasicAuth
// @Failure 401 {object} models.ErrorResponse "Invalid credentials"
func AdminAuth() gin.HandlerFunc {
	auth := gin.BasicAuth(gin.Accounts{
		"SuperUser": "qwerty12345",
	})
	return func(c *gin.Context) {
		auth(c)
		if c.IsAborted() {
			c.AbortWithStatusJSON(http.StatusUnauthorized, models.ErrorResponse{
				Error: "check username or password",
			})
			return
		}

		c.Next()
	}
}
