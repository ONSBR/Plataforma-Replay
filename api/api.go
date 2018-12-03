package api

import (
	"github.com/ONSBR/Plataforma-EventManager/infra"
	"github.com/ONSBR/Plataforma-Replay/api/handlers"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/labstack/gommon/log"
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

	tapeGroup.POST("/upload", handlers.UploadTape)
	tapeGroup.DELETE("/:tapeID", handlers.DeleteTape)
	tapeGroup.POST("/:systemID/rec", handlers.StartRecording)
	tapeGroup.POST("/:systemID/stop", handlers.StopRecording)
	tapeGroup.GET("/:systemID/availables", handlers.AvailableTapes)
	tapeGroup.GET("/:systemID/recording", handlers.GetRecording)
	tapeGroup.GET("/:systemID/download/:id", handlers.DownloadTape)
	tapeGroup.POST("/:systemID/play", handlers.Play)

	// Start server
	port := infra.GetEnv("PORT", "6081")
	log.Info(port)
	e.Logger.Fatal(e.Start(":" + port))
}
