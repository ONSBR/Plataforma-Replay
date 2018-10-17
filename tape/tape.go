package tape

import (
	"github.com/ONSBR/Plataforma-EventManager/domain"
)

//Tape is datalog struct to keep all events, dump file and metadata
type Tape struct {
	path     string
	state    string
	segments []*Segment
}

type Segment struct {
	fileName    string
	segmentType string
}

func (t *Tape) Record(event *domain.Event) error {
	return nil
}

func (t *Tape) RecordDump(systemID string) error {
	return nil
}

func (t *Tape) Dest() string {
	return t.path
}

func GetTape(systemID string) *Tape {
	return nil
}
