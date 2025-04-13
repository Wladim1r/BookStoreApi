package middleware

import "github.com/gin-gonic/gin"

func Authentification() gin.HandlerFunc {
	auth := gin.BasicAuth(gin.Accounts{
		"admin": "secret",
	})

	return auth
}
