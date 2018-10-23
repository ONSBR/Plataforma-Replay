package handlers

import (
	"github.com/ONSBR/Plataforma-Replay/recorder"
	"github.com/labstack/echo"
)

//StopRecording for an application
func StopRecording(ctx echo.Context) error {
	systemID := ctx.Get("systemID").(string)
	rec := recorder.GetRecorder(systemID)
	tape, err := rec.GetTape(systemID)
	if err != nil {
		return err
	}
	err = tape.Close()
	if err != nil {
		return err
	}
	return ctx.JSON(201, tape)
}
