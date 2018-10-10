package tape

import (
	"github.com/ONSBR/Plataforma-EventManager/domain"
)

//Tape is datalog struct to keep all events, dump file and metadata
type Tape struct {
	//Dump
	//Events
	//Metadata
}

func (t *Tape) Record(event *domain.Event) error {
	return nil
}

func GetTape(systemID string) *Tape {
	return nil
}
