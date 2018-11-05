package handlers

import (
	"github.com/ONSBR/Plataforma-Replay/recorder"
	"github.com/labstack/echo"
)

//DownloadTape return available tapes for system
func DownloadTape(ctx echo.Context) error {
	systemID := ctx.Param("systemID")
	tapeZipID := ctx.Param("id")
	rec := recorder.GetRecorder(systemID)
	return ctx.File(rec.TapePath(tapeZipID))
}
