package player

import "github.com/ONSBR/Plataforma-Replay/tape"

type Player interface {
	Play(tapeID string) error
}

func GetPlayer(systemID string) Player {
	p := new(defaultPlayer)
	p.path = tape.GetTapesPath()
	p.systemID = systemID
	return p
}
