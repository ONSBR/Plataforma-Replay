package handlers

import (
	"fmt"

	"github.com/ONSBR/Plataforma-Replay/recorder"
	"github.com/labstack/echo"
)

//GetRecording return available tapes for system
func GetRecording(ctx echo.Context) error {
	systemID := ctx.Param("systemID")

	rec := recorder.GetRecorder(systemID)
	tape, err := rec.GetTape(systemID)
	if err != nil {
		return ctx.JSON(404, map[string]interface{}{"message": fmt.Sprintf("system %s is not in recording mode", systemID)})
	}
	return ctx.JSON(200, tape)
}
