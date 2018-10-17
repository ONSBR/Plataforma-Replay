package recorder

import (
	"github.com/ONSBR/Plataforma-EventManager/domain"
	"github.com/ONSBR/Plataforma-Replay/tape"
)

//Recorder controls the flow of saving information about database dumping and executed events
type Recorder interface {
	Insert(tape *tape.Tape) error
	Rec(event *domain.Event) error
	Eject(systemID string) (*tape.Tape, error)
	IsClosed() (bool, error)
}

//GetRecorder returns default implementation
func GetRecorder() Recorder {
	return new(DefaultRecorder)
}
