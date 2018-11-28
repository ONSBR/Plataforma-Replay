package rabbitmq

import (
	"github.com/ONSBR/Plataforma-Deployer/env"
	"github.com/PMoneda/carrot"
	"github.com/labstack/gommon/log"
)

var broker *RabbitBroker

//Init broker connection
func Init() {

	config := carrot.ConnectionConfig{
		Host:     env.Get("RABBITMQ_HOST", "localhost"),
		Username: env.Get("RABBITMQ_USERNAME", "guest"),
		Password: env.Get("RABBITMQ_PASSWORD", "guest"),
		VHost:    "/",
	}

	conn, err := carrot.NewBrokerClient(&config)
	if err != nil {
		log.Error(err)
		panic(err)
	}
	builder := carrot.NewBuilder(conn)
	builder.UseVHost("plataforma_v1.0")
	subConn, _ := carrot.NewBrokerClient(&config)
	subscriber := carrot.NewSubscriber(subConn)
	pubConn, _ := carrot.NewBrokerClient(&config)
	publisher := carrot.NewPublisher(pubConn)

	broker = NewRabbitBroker()
	broker.Publisher = publisher
	broker.Subscriber = subscriber
}

func GetBroker() *RabbitBroker {
	return broker
}
