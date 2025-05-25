package main

import (
	hand "bookstore-api/api/handlers"
	repo "bookstore-api/api/repository"
	serv "bookstore-api/api/service"
	_ "bookstore-api/docs"
	db "bookstore-api/internal/database"
	"bookstore-api/internal/middleware"
	"bookstore-api/internal/utils"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title BookStore API
// @version 1.1.3
// @description RESP API for managing books and personal user list of books
// @host localhost:8080
// @BasePath /
//
// @securityDefinitions.basic BasicAuth
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
// @description JWT token with 'Bearer ' prefix. Example: `Bearer eyJhbGci...`
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

	url := ginSwagger.URL("/swagger/doc.json")
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))

	// #######################___PUBLIC___######################
	public := r.Group("/auth")
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
