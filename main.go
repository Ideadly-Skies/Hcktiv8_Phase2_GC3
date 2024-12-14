package main

import (
	config "w3/gc3/config/database"
	user_handler "w3/gc3/internal/userHandler"
	post_handler "w3/gc3/internal/postHandler"
	comment_handler "w3/gc3/internal/commentHandler"
	activity_handler "w3/gc3/internal/activityHandler"
	cust_middleware "w3/gc3/internal/middleware"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/swaggo/echo-swagger"
	_ "w3/gc3/docs"
)

func main(){
	// migrate data to supabase
	// config.MigrateData()

	// connect to db
	config.InitDB()
	defer config.CloseDB()

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// public routes
	e.POST("users/register", user_handler.Register)
	e.POST("users/login", user_handler.Login)

	// protected routes //
	// post
	e.POST("posts", post_handler.CreatePost, cust_middleware.JWTMiddleware)
	e.GET("posts", post_handler.GetAllPosts, cust_middleware.JWTMiddleware)
	e.GET("posts/:id", post_handler.GetPostByID, cust_middleware.JWTMiddleware)
	e.DELETE("posts/:id", post_handler.DeletePost, cust_middleware.JWTMiddleware)	

	// comments
	e.POST("/comments", comment_handler.CreateComment, cust_middleware.JWTMiddleware)
	e.GET("/comments/:id", comment_handler.GetCommentByID, cust_middleware.JWTMiddleware)
	e.DELETE("/comments/:id", comment_handler.DeleteCommentByID, cust_middleware.JWTMiddleware)

	// activity
	e.GET("activities", activity_handler.GetActivities, cust_middleware.JWTMiddleware)

	// swagger
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	// start the server at 8080
	e.Logger.Fatal(e.Start(":8080"))
}