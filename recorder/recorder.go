package recorder

import (
	"github.com/ONSBR/Plataforma-EventManager/domain"
	"github.com/ONSBR/Plataforma-Replay/tape"
)

//Recorder controls the flow of saving information about database dumping and executed events
type Recorder interface {
	Rec(event *domain.Event) error
	Eject() (*tape.Tape, error)
	IsRecording() bool
	GetOrCreateTape(systemID string) (*tape.Tape, error)
	GetTape(systemID string) (*tape.Tape, error)
	InsertTape(tape *tape.Tape)
	AvailableTapesToDownload(systemID string) ([]string, error)
	TapePath(name string) string
}

//GetRecorder returns default implementation
func GetRecorder(systemID string) Recorder {
	return newDefaultRecorder(systemID)
}

func newDefaultRecorder(systemID string) Recorder {
	rec := new(DefaultRecorder)
	rec.systemID = systemID
	rec.path = tape.GetTapesPath()
	return rec
}
