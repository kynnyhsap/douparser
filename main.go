package main

import (
	"dou-parser/handlers"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {

	r := echo.New()

	r.Pre(middleware.RemoveTrailingSlash())

	r.GET("/tags", handlers.GetTags)
	r.GET("/calendar/events", handlers.GetEventsList)
	r.GET("/calendar/events/:id", handlers.GetEvent)

	r.Logger.Fatal(r.Start(":8080"))
}
