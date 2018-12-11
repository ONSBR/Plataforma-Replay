package player

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/ONSBR/Plataforma-Deployer/sdk/eventmanager"

	"github.com/PMoneda/whaler"

	"github.com/ONSBR/Plataforma-Replay/sdk"

	"github.com/ONSBR/Plataforma-EventManager/domain"
	"github.com/labstack/gommon/log"

	"github.com/ONSBR/Plataforma-Replay/db"

	"github.com/ONSBR/Plataforma-Replay/recorder"
	"github.com/ONSBR/Plataforma-Replay/tape"
)

type defaultPlayer struct {
	path     string
	systemID string
}

/*
O processo de play é feito da seguinte forma:

Pela tela, o usuário clica no botão play de uma respectiva fita, já disponível dentro da plataforma.
Quando clicar em Play o serviço verifica se a plataforma está em modo gravação e caso esteja retorna uma mensagem de erro
dizendo que não pode iniciar o processo de play porque a plataforma está em processo de gravação

Caso a plataforma não esteja em processo de gravação, ela irá primeiramente
1 - descompactar a fita
2 - se existir um banco de dados de mesmo nome a plataforma irá fazer uma fita de backup
3 - fazer o processo de restore da base de dados
4 - Para cada evento, o serviço irá dispará-lo como se fosse um execução
*/
func (p *defaultPlayer) Play(tapeID string) error {
	rec := recorder.GetRecorder(p.systemID)
	if rec.IsRecording() {
		return fmt.Errorf("Cannot start playing events while platform is recording")
	}
	dbName, err := sdk.GetDBName(p.systemID)
	if err != nil {
		return err
	}
	if err := p.StopDomainContaners(dbName); err != nil {
		return err
	}
	if err := tape.Restore(p.systemID, tapeID); err != nil {
		return err
	}
	if err := p.RestoreDataBase(); err != nil {
		return err
	}
	if err := p.StartDomainContaners(dbName); err != nil {
		return err
	}
	if events, err := p.GetEventsFromTape(); err != nil {
		return err
	} else if err := p.EmitEvents(events); err != nil {
		return err
	}
	//in the end remove opened tape after send all events
	t := tape.Tape{SystemID: p.systemID, Path: tape.GetTapesPath()}
	log.Info("clear tape folder ", t.Dest())
	if err := os.RemoveAll(t.Dest()); err != nil {
		log.Error(err)
		return err
	}
	return nil
}

func (p *defaultPlayer) StopDomainContaners(name string) error {
	defer func() {
		if r := recover(); r != nil {
			log.Error("Recovered in f", r)
		}
	}()
	containers, err := whaler.GetContainers(false)
	if err != nil {
		return err
	}
	var timeout time.Duration = 30 * time.Second
	for _, container := range containers {
		if strings.Contains(container.Name, name) {
			if err := whaler.StopContainer(container.ID, &timeout); err != nil {
				return err
			}
		}
	}
	return nil
}

func (p *defaultPlayer) StartDomainContaners(name string) error {
	containers, err := whaler.GetContainers(true)
	if err != nil {
		return err
	}
	for _, container := range containers {
		if strings.Contains(container.Name, name) {
			if err := whaler.StartContainer(container.ID); err != nil {
				return err
			}
		}
	}
	return nil
}

func (p *defaultPlayer) RestoreDataBase() error {
	dbName, err := sdk.GetDBName(p.systemID)
	if err != nil {
		return err
	}
	return db.GetDB().Restore(dbName, fmt.Sprintf("%s/%s/dump.sql", tape.GetTapesPath(), p.systemID))
}

func (p *defaultPlayer) GetEventsFromTape() ([]*domain.Event, error) {
	t, err := tape.GetOrCreateTape(p.systemID, tape.GetTapesPath())
	if err != nil {
		return nil, err
	}
	events := make([]*domain.Event, 0)
	for _, seg := range t.Segments {
		if seg.IsEvent() {
			if event, err := seg.Event(t); err != nil {
				return nil, err
			} else {
				events = append(events, event)
			}
		}
	}
	return events, nil
}

func (p *defaultPlayer) EmitEvents(events []*domain.Event) error {
	for _, evt := range events {
		if err := p.Emit(evt); err != nil {
			return err
		}
	}
	return nil
}

func (p *defaultPlayer) Emit(event *domain.Event) error {
	log.Info("emiting event ", event.Name)
	return eventmanager.Push(event)
}
