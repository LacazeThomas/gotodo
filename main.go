package main

import (
	"github.com/caarlos0/env"
	log "github.com/sirupsen/logrus"

	"github.com/lacazethomas/goTodo/app"
	"github.com/lacazethomas/goTodo/config"
)

//TODO:
//Better file organization
func main() {
	log.SetFormatter(&log.JSONFormatter{})

	cfg := config.DB{}
	err := env.Parse(&cfg)
	if err != nil {
		log.Println("Can not get env variable")
	}

	router := &app.App{}
	router.Initialize(cfg)
	router.Run(":8000")
}
