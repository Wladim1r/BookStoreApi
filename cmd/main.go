package main

import (
	"bookstore-api/internal/api"
	"bookstore-api/internal/database"
	"bookstore-api/internal/service"
	"bookstore-api/internal/utils"

	"github.com/gin-gonic/gin"
)

func main() {
	gin.SetMode(gin.ReleaseMode)

	db := database.InitDB()

	bookRepo := service.NewBookRepository(db)
	bookServ := service.NewBookService(bookRepo)
	bookHandler := api.NewBookHandler(bookServ)

	r := gin.Default()
	r.Use(gin.LoggerWithFormatter(utils.Log))

	admin := r.Group("/admin")
	{
		admin.PATCH("/books/:id", bookHandler.UpdateBook)
		admin.DELETE("/books/:id", bookHandler.DeleteBook)

	}

	r.GET("/books", bookHandler.GetBooks)
	r.GET("/books/:id", bookHandler.GetBook)
	r.POST("/books", bookHandler.PostBook)

	r.Run(":8080")
}
