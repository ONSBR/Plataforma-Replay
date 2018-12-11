package handlers

import (
	"github.com/ONSBR/Plataforma-EventManager/domain"

	"github.com/ONSBR/Plataforma-Replay/recorder"

	"github.com/labstack/echo"
)

func RecordEvent(ctx echo.Context) error {
	systemID := ctx.Param("systemID")
	rec := recorder.GetRecorder(systemID)
	event := new(domain.Event)
	ctx.Bind(event)
	if err := rec.Rec(event); err != nil {
		return ctx.JSON(400, H{"message": err.Error()})
	} else {
		return ctx.JSON(200, H{"message": "ok"})
	}
}
