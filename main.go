package main

import (
	"fmt"

	"github.com/ONSBR/Plataforma-Replay/actions"
	"github.com/ONSBR/Plataforma-Replay/api"
	"github.com/ONSBR/Plataforma-Replay/broker"
	"github.com/labstack/gommon/log"
)

func main() {
	log.SetLevel(log.INFO)
	fmt.Println(logo())
	go broker.Init(func(broker broker.Broker) {
		broker.Subscribe(actions.ReceiveEvent)
	})

	api.RunAPI()
}

func logo() string {
	return `

                 _
                | |
  _ __ ___ _ __ | | __ _ _   _
 | '__/ _ \ '_ \| |/ _' | | | |
 | | |  __/ |_) | | (_| | |_| |
 |_|  \___| .__/|_|\__,_|\__, |
          | |             __/ |
          |_|            |___/

	`
}
