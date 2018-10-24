package api

import (
	"github.com/ONSBR/Plataforma-EventManager/infra"
	"github.com/ONSBR/Plataforma-Replay/api/handlers"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

//RunAPI build, config and run API
func RunAPI() {
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	group := e.Group("v1")

	// Routes
	group.POST("/startRecording/:systemId", handlers.StartRecording)
	group.POST("/stopRecording/:systemId", handlers.StopRecording)
	group.GET("/tape/:systemId/available", handlers.AvailableTapes)
	group.GET("/tape/:systemId/recording", handlers.GetRecording)
	group.GET("/tape/:systemId/download/:id", handlers.DownloadTape)
	group.POST("/play/:systemId", handlers.Play)

	// Start server
	port := infra.GetEnv("PORT", ":6081")
	e.Logger.Fatal(e.Start(port))
}
