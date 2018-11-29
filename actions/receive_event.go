package actions

import (
	"github.com/ONSBR/Plataforma-EventManager/domain"
	"github.com/ONSBR/Plataforma-Replay/broker"
	"github.com/ONSBR/Plataforma-Replay/recorder"
	"github.com/labstack/gommon/log"
)

/*
ReceiveEvent Processa um evento, antes de ser enviado a fila de execução, para gravar toda a execução de um determinado sistema
Esse processo é feito da seguinte forma:
	* Pegamos uma fita e um gravador e solicitamos ao gravador para gravar o evento de chegada na fita exclusiva de um sistema
	  Cada sistema terá sua fila em disco exclusiva onde cada fita irá conter o dump da base de dados do início da gravação e
	  os eventos que foram executados a partir do início da gravação;
	* O momento da gravação ocorre antes do evento ir para a fila de execução da plataforma.

Pontos de falha:
	* Se houver uma falha na verificação da disponibilidade de gravação de uma fita então será apenas logado um erro
	  e o processo de execução continua
	* Se a fita já estiver fechada para gravação então será logada apenas uma informação de que a fita já foi fechada e
	  a execução continuará sem a gravação do evento
	* Se na tentativa de gravação na fita ocorrer um erro então a fita será marcada como corrompida, será logado o erro
	  de gravação e a execução seguirá mas a gravação desse sistema será interrompida, será necessário habilitar novamente
	  o modo de gravação para gerar uma outra fita

Observação: A ideia é que o serviço de replay não seja um gargalo de execução para a plataforma
mas a execução do processo irá continuar;
TODO: traduzir para o inglês de forma adequada!*/
func ReceiveEvent(event *domain.Event) error {
	log.Info("received event %s from system %s", event.Name, event.SystemID)
	recorder := recorder.GetRecorder(event.SystemID)
	log.Info("recording event")
	err := recorder.Rec(event)
	log.Info("event recorded")
	if err != nil {
		log.Error(err)
	}
	brk := broker.GetBroker()
	log.Info("publishing event to execution")
	return brk.Publish(event)
}
