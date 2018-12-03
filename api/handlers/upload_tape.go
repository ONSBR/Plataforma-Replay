package handlers

import (
	"fmt"
	"io"
	"os"

	"github.com/ONSBR/Plataforma-Replay/tape"
	"github.com/labstack/echo"
	"github.com/labstack/gommon/log"
)

//UploadTape upload tapes
func UploadTape(c echo.Context) error {
	log.Info("uploading tape")
	file, err := c.FormFile("tape")
	if err != nil {
		log.Error(err)
		return err
	}
	src, err := file.Open()
	if err != nil {
		log.Error(err)
		return err
	}
	defer src.Close()
	log.Info(fmt.Sprintf("saving tape %s", file.Filename))
	// Destination
	dst, err := os.Create(fmt.Sprintf("%s/%s", tape.GetTapesPath(), file.Filename))
	if err != nil {
		log.Error(err)
		return err
	}
	defer dst.Close()

	// Copy
	if _, err = io.Copy(dst, src); err != nil {
		log.Error(err)
		return err
	}
	return c.JSON(200, H{"message": fmt.Sprintf("tape %s was uploaded", file.Filename)})
}
