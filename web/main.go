package main

import (
	"github.com/labstack/echo/v4"
	"github.com/remshams/jira-control/web/handler"
)

func main() {
	e := echo.New()
	e.Static("/assets", "assets")
	e.GET("/log", handler.HandleLogWork)

	e.Logger.Fatal(e.Start(":8080"))
}
