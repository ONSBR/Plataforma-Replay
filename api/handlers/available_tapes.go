package handlers

import (
	"github.com/ONSBR/Plataforma-Replay/recorder"
	"github.com/labstack/echo"
)

//AvailableTapes return available tapes for system
func AvailableTapes(ctx echo.Context) error {
	systemID := ctx.Get("systemID").(string)
	rec := recorder.GetRecorder(systemID)
	list, err := rec.AvailableTapesToDownload(systemID)
	if err != nil {
		return err
	}
	return ctx.JSON(200, list)
}
