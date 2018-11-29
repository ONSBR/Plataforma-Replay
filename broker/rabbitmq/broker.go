package rabbitmq

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/ONSBR/Plataforma-Deployer/sdk/apicore"
	"github.com/ONSBR/Plataforma-EventManager/domain"
	"github.com/PMoneda/carrot"
	"github.com/labstack/gommon/log"
)

type RabbitBroker struct {
	Subscriber *carrot.Subscriber
	Publisher  *carrot.Publisher
}

//Publish an event to executor queue
func (broker *RabbitBroker) Publish(event *domain.Event) error {
	body, err := json.Marshal(event.ToCeleryMessage())
	if err != nil {
		return err
	}
	return broker.Publisher.Publish("events.publish", "executor", carrot.Message{
		ContentType: "application/json",
		Encoding:    "utf-8",
		Data:        body,
	})
}

func (broker *RabbitBroker) Subscribe(action func(*domain.Event) error) error {
	systems := getInstalledSystems()
	for _, system := range systems {
		err := broker.Subscriber.Subscribe(carrot.SubscribeWorker{
			AutoAck: false,
			Queue:   fmt.Sprintf("event.replay.%s.queue", system),
			Scale:   1,
			Handler: func(context *carrot.MessageContext) error {
				celeryMessage := new(domain.CeleryMessage)
				if err := json.Unmarshal(context.Message.Data, celeryMessage); err != nil {
					return err
				} else {
					if len(celeryMessage.Args) > 0 {
						if err := action(&celeryMessage.Args[0]); err != nil {
							return context.Nack(true)
						} else {
							return context.Ack()
						}
					} else {
						return context.Nack(false)
					}
				}
			},
		})
		if err != nil {
			log.Error(err)
		}
	}
	return nil
}

func getInstalledSystems() []string {
	type system struct {
		ID string `json:"id"`
	}
	list := make([]system, 0)
	for {
		err := apicore.Query(apicore.Filter{
			Map:    "core",
			Entity: "system",
			Name:   "",
		}, &list)
		if err == nil {
			break
		} else {
			log.Error("cannot connect to apicore retry in 10 seconds...")
			time.Sleep(10 * time.Second)
		}
	}

	ids := make([]string, len(list))
	for i, v := range list {
		ids[i] = v.ID
	}
	return ids
}

func parseMessage(message interface{}) (body []byte, err error) {
	switch t := message.(type) {
	case []byte:
		body = t
	default:
		body, err = json.Marshal(t)
		if err != nil {
			return nil, err
		}
		return
	}
	return
}

func NewRabbitBroker() *RabbitBroker {
	return new(RabbitBroker)
}
