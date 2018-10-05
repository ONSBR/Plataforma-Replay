package recorder

import (
	"github.com/ONSBR/Plataforma-Replay/tape"
)

type Recorder interface {
	Rec() error
	Stop() error
	Eject() (*tape.Tape, error)
	Rewind() error
}
