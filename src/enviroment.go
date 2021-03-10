package main

import (
	"log"

	"github.com/kelseyhightower/envconfig"
)

type Environment struct {
	Topics     string `envconfig:"KAFKA_WEBSOCKET_SERVER_TOPICS" required:"true"`
	BrokerURL  string `envconfig:"KAFKA_WEBSOCKET_SERVER_BROCKER_URL" required:"true"`
	Port       string `envconfig:"KAFKA_WEBSOCKET_SERVER_PORT" required:"false" default:"8080"`
	Prefix     string `envconfig:"KAFKA_WEBSOCKET_SERVER_PREFIX" required:"false" default:""`
	HealthPort string `envconfig:"KAFKA_WEBSOCKET_SERVER_HEALTH_PORT" required:"false" default:"5001"`
}

var env *Environment

func getEnvironment() *Environment {

	if env == nil {
		err := envconfig.Process("", env)
		if err != nil {
			log.Fatal(err.Error())
		}
	}

	return env
}
