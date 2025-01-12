package main

import (
	"brain/auth"
	"brain/card"
	"brain/db"
	"brain/user"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	// Initialize database
	err := db.InitDB()
	if err != nil {
		panic(err)
	}

	// Create Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Login route
	e.POST("/login", auth.Login)

	// Restricted routes
	restricted := e.Group("/cards")
	restricted.Use(auth.AuthMiddleware)

	// Card routes
	restricted.GET("", card.GetCards)
	restricted.GET("/:id", card.GetCard)
	restricted.POST("", card.CreateCard)
	restricted.PUT("/:id", card.UpdateCard)
	restricted.DELETE("/:id", card.DeleteCard)

	// User routes
	e.GET("/users", user.GetUsers)
	e.GET("/users/:id", user.GetUser)
	e.POST("/users", user.CreateUser)
	e.PUT("/users/:id", user.UpdateUser)
	e.DELETE("/users/:id", user.DeleteUser)

	// Start server
	e.Logger.Fatal(e.Start(":1323"))
}
