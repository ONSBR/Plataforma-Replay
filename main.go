package main

import (
	"fmt"

	"github.com/ONSBR/Plataforma-Replay/api"
)

func main() {
	fmt.Println(logo())

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
