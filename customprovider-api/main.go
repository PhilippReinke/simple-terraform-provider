package main

import (
	"api/items"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	items := items.New()
	items.AddItem("item1")

	e := echo.New()
	e.Use(middleware.Logger())

	// create
	e.POST("/items/:name", func(c echo.Context) error {
		name := c.Param("name")
		item := items.AddItem(name)
		return c.JSON(http.StatusOK, item)
	})

	// read
	e.GET("/items", func(c echo.Context) error {
		return c.JSON(http.StatusOK, items.ReadItems())
	})
	e.GET("/items/:id", func(c echo.Context) error {
		id := c.Param("id")
		item, err := items.ReadItemByID(id)
		if err != nil {
			return c.NoContent(http.StatusNotFound)
		}
		return c.JSON(http.StatusOK, item)
	})

	// update
	e.PUT("/items/:id/:name", func(c echo.Context) error {
		id := c.Param("id")
		name := c.Param("name")
		item, err := items.Update(id, name)
		if err != nil {
			return c.NoContent(http.StatusNotFound)
		}
		return c.JSON(http.StatusOK, item)
	})

	// remove
	e.DELETE("/items/:id", func(c echo.Context) error {
		id := c.Param("id")
		items.Remove(id)
		return c.NoContent(http.StatusOK)
	})

	e.Logger.Fatal(e.Start(":8080"))
}
