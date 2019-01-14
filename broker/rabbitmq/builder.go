package rabbitmq

import (
	"fmt"
	"strconv"
	"time"

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
		count, e := strconv.Atoi(env.Get("RABBITMQ_MAX_RETRIES", "30"))
		if e != nil {
			count = 30
		}
		log.Error(err)
		for err != nil && count > 0 {
			log.Info(fmt.Sprintf("trying again after 10 seconds: remaining %d", count))
			time.Sleep(10 * time.Second)
			conn, err = carrot.NewBrokerClient(&config)
			count--
		}
		log.Info("connected to rabbitmq")
		if count == 0 && err != nil {
			panic(err)
		}
	}
	builder := carrot.NewBuilder(conn)
	builder.UseVHost("plataforma_v1.0")
	subConn, _ := carrot.NewBrokerClient(&config)
	subscriber := carrot.NewSubscriber(subConn)
	subscriber.SetMaxRetries(30)
	pubConn, _ := carrot.NewBrokerClient(&config)
	publisher := carrot.NewPublisher(pubConn)
	broker = NewRabbitBroker()
	broker.Publisher = publisher
	broker.Subscriber = subscriber
}

func GetBroker() *RabbitBroker {
	return broker
}
