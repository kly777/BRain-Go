package main

import (


	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"yourproject/db"
	"yourproject/card"
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

	// Routes
	e.GET("/cards", card.GetCards)
	e.GET("/cards/:id", card.GetCard)
	e.POST("/cards", card.CreateCard)
	e.PUT("/cards/:id", card.UpdateCard)
	e.DELETE("/cards/:id", card.DeleteCard)

	// Start server
	e.Logger.Fatal(e.Start(":1323"))
}
