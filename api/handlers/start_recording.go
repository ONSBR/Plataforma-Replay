package handlers

import (
	"fmt"

	"github.com/ONSBR/Plataforma-Replay/recorder"
	"github.com/labstack/echo"
	"github.com/labstack/gommon/log"
)

//StartRecording starts to record all request events for an application
func StartRecording(ctx echo.Context) error {
	systemID := ctx.Param("systemID")
	rec := recorder.GetRecorder(systemID)
	log.Info(fmt.Sprintf("creating tape for system %s", systemID))
	tape, err := rec.GetOrCreateTape(systemID)
	if err != nil {
		return err
	}
	if !tape.IsRecording() {
		ctx.JSON(400, H{"message": "cannot start recording"})
	}
	return ctx.JSON(201, tape)
}
