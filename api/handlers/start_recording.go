package handlers

import (
	"github.com/ONSBR/Plataforma-Replay/recorder"
	"github.com/labstack/echo"
)

//StartRecording starts to record all request events for an application
func StartRecording(ctx echo.Context) error {
	systemID := ctx.Get("systemID").(string)
	rec := recorder.GetRecorder(systemID)
	tape, err := rec.GetOrCreateTape(systemID)
	if err != nil {
		return err
	}
	return ctx.JSON(201, tape)
}
