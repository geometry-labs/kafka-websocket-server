package main

import (
	"log"

	"github.com/kelseyhightower/envconfig"
)

type Environment struct {
	Topics      string `envconfig:"KAFKA_WEBSOCKET_SERVER_TOPICS" required:"true"`
	BrokerURL   string `envconfig:"KAFKA_WEBSOCKET_SERVER_BROKER_URL" required:"true"`
	Port        string `envconfig:"KAFKA_WEBSOCKET_SERVER_PORT" required:"false" default:"8080"`
	Prefix      string `envconfig:"KAFKA_WEBSOCKET_SERVER_PREFIX" required:"false" default:""`
	HealthPort  string `envconfig:"KAFKA_WEBSOCKET_SERVER_HEALTH_PORT" required:"false" default:"5001"`
	MetricsPort string `envconfig:"KAFKA_WEBSOCKET_SERVER_METRICS_PORT" required:"false" default:"9400"`
	NetworkName string `envconfig:"KAFKA_WEBSOCKET_SERVER_NETWORK_NAME" required:"false" default:"mainnet"`
}

func getEnvironment() Environment {
	var env Environment
	err := envconfig.Process("", &env)
	if err != nil {
		log.Fatalf("ERROR: envconfig - %s\n", err.Error())
	}
	return env
}
