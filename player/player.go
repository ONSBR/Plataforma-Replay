package player

import (
	"github.com/ONSBR/Plataforma-EventManager/domain"
	"github.com/ONSBR/Plataforma-Replay/tape"
)

type Player interface {
	Play(tapeID string) error
	StopDomainContaners(name string) error
	StartDomainContaners(name string) error
	RestoreDataBase() error
	GetEventsFromTape() ([]*domain.Event, error)
	EmitEvents(events []*domain.Event) error
	Emit(event *domain.Event) error
}

func GetPlayer(systemID string) Player {
	p := new(defaultPlayer)
	p.path = tape.GetTapesPath()
	p.systemID = systemID
	return p
}
