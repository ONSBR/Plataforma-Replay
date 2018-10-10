package broker

import (
	"github.com/ONSBR/Plataforma-EventManager/domain"
	"github.com/ONSBR/Plataforma-Replay/broker/rabbitmq"
)

type Broker interface {
	Publish(event *domain.Event) error
	Subscribe(action func(*domain.Event) error) error
}

func GetBroker() Broker {
	return rabbitmq.GetBroker()
}

func Init(callback func(broker Broker)) {
	rabbitmq.Init()
	callback(GetBroker())
}
