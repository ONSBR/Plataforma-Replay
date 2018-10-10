package rabbitmq

import (
	"os"

	"github.com/PMoneda/carrot"
	"github.com/labstack/gommon/log"
)

var broker *RabbitBroker

//Init broker connection
func Init() {

	config := carrot.ConnectionConfig{
		Host:     os.Getenv("RABBITMQ_HOST"),
		Username: os.Getenv("RABBITMQ_USERNAME"),
		Password: os.Getenv("RABBITMQ_PASSWORD"),
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
