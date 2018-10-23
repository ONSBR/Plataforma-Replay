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
	group.POST("/startRecording/:system_id", handlers.StartRecording)
	group.POST("/stopRecording/:system_id", handlers.StopRecording)
	group.POST("/play", handlers.Play)

	// Start server
	port := infra.GetEnv("PORT", ":6081")
	e.Logger.Fatal(e.Start(port))
}
