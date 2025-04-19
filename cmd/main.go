package main

import (
	hand "bookstore-api/api/handlers"
	repo "bookstore-api/api/repository"
	serv "bookstore-api/api/service"
	db "bookstore-api/internal/database"
	"bookstore-api/internal/middleware"
	"bookstore-api/internal/utils"

	"github.com/gin-gonic/gin"
)

func main() {
	gin.SetMode(gin.ReleaseMode)

	db := db.InitDB()

	bookRepo := repo.NewBookRepository(db)
	bookServ := serv.NewBookService(bookRepo)
	bookHandler := hand.NewBookHandler(bookServ)

	userRepo := repo.NewUserRepository(db)
	userServ := serv.NewUserService(userRepo)
	userHandler := hand.NewUserHandler(userServ)

	r := gin.Default()
	r.Use(gin.LoggerWithFormatter(utils.Log))

	// #######################___PUBLIC___######################
	public := r.Group("/api")
	{
		public.POST("/register", userHandler.Register)
		public.POST("/login", userHandler.Login)
	}
	// #########################################################
	//
	//
	//
	// #######################___PRIVATE___#####################
	private := r.Group("/api")
	private.Use(middleware.JWTAuth())
	{
		private.GET("/books/:id", bookHandler.GetUserBook)
		private.GET("/books", bookHandler.GetUserBooks)
		private.POST("/books", bookHandler.PostBook)
		private.PATCH("/books/:id", bookHandler.UpdateBook)
		private.DELETE("/books/:id", bookHandler.DeleteBook)
	}
	// #########################################################
	//
	//
	//
	// ########################___ADMIN___######################
	admin := r.Group("/admin")
	admin.Use(middleware.AdminAuth())
	{
		admin.GET("/books", bookHandler.GetAllBooks)
		admin.GET("/users", userHandler.GetAllUsers)
		admin.DELETE("/users/:username", userHandler.DeleteByUsername)
	}
	// #########################################################

	r.Run(":8080")
}
