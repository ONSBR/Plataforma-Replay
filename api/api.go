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
	tapeGroup := group.Group("/tape")
	// Routes
	tapeGroup.POST("/:systemId/rec", handlers.StartRecording)
	tapeGroup.POST("/:systemId/stop", handlers.StopRecording)
	tapeGroup.GET("/:systemId/available", handlers.AvailableTapes)
	tapeGroup.GET("/:systemId/recording", handlers.GetRecording)
	tapeGroup.GET("/:systemId/download/:id", handlers.DownloadTape)
	tapeGroup.POST("/:systemId/play", handlers.Play)

	// Start server
	port := infra.GetEnv("PORT", ":6081")
	e.Logger.Fatal(e.Start(port))
}
