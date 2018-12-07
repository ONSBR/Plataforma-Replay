package handlers

import (
	"fmt"

	"github.com/ONSBR/Plataforma-Replay/player"
	"github.com/labstack/echo"
	"github.com/labstack/gommon/log"
)

//Play a tape on current installed platform
func Play(ctx echo.Context) error {
	tapeID := ctx.Param("id")
	p := player.GetPlayer(ctx.Param("systemID"))
	log.Info(fmt.Sprintf("Tape %s on system %s", tapeID, ctx.Param("systemID")))
	err := p.Play(tapeID)
	if err != nil {
		log.Error(err)
		return ctx.JSON(400, H{"message": err.Error()})
	}
	return ctx.JSON(200, H{"message": "playing"})
}
