package main

import (
	"github.com/caarlos0/env"
	log "github.com/sirupsen/logrus"

	"github.com/lacazethomas/goTodo/app"
	"github.com/lacazethomas/goTodo/config"
)

func main() {

	log.SetFormatter(&log.JSONFormatter{})

	config := config.DB{}
	err := env.Parse(&config)
	if err != nil {
		log.Println("Can not get env variable")
	}
	log.Printf("%+v\n", config)
	app := &app.App{}
	app.Initialize(config)
	app.Run(":8000")
}
