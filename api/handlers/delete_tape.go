package handlers

import (
	"fmt"

	"github.com/ONSBR/Plataforma-Replay/tape"
	"github.com/labstack/echo"
)

//DeleteTape by id
func DeleteTape(ctx echo.Context) error {
	tapeZipID := ctx.Param("tapeID")
	if err := tape.Delete(tapeZipID); err == nil {
		return ctx.JSON(200, H{"message": fmt.Sprintf("Tape %s was deleted", tapeZipID)})
	} else {
		return ctx.JSON(400, H{"message": err})
	}

}
