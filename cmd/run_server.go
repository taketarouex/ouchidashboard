package main

import (
	"fmt"
	"os"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func main() {
	e := echo.New()

	e.Use(middleware.Logger())
	e.POST("/", collectorHandler)
	e.GET("/api/rooms/:roomName/logs/:logType", getLogsHandler)
	e.GET("/api/rooms", getRoomNamesHandler)
	e.Static("/ui", "ui")
	e.Static("/_next", "ui/_next")
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", port)))
}
