package player

import "github.com/ONSBR/Plataforma-Replay/tape"

type Player interface {
	Play(tape.Tape)
}
