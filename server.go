package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	e.GET("/cards", GetCards)
	e.GET("/cards/:id", GetCard)
	e.POST("/cards", CreateCard)
	e.PUT("/cards/:id", UpdateCard)
	e.DELETE("/cards/:id", DeleteCard)

	e.Logger.Fatal(e.Start(":1323"))
}
