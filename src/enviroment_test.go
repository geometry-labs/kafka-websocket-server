package main

import (
	"os"

	"testing"
)

func TestEnvironment(t *testing.T) {

	// Set env
	env_map := map[string]string{
		"KAFKA_WEBSOCKET_SERVER_TOPICS":      "topics",
		"KAFKA_WEBSOCKET_SERVER_BROKER_URL":  "broker_url_env",
		"KAFKA_WEBSOCKET_SERVER_PORT":        "port_env",
		"KAFKA_WEBSOCKET_SERVER_PREFIX":      "prefix_env",
		"KAFKA_WEBSOCKET_SERVER_HEALTH_PORT": "health_port_env",
	}

	for k, v := range env_map {
		os.Setenv(k, v)
	}

	// Check env
	env := getEnvironment()

	if env.Topics != env_map["KAFKA_WEBSOCKET_SERVER_TOPICS"] {
		t.Errorf("Invalid value for env variable: KAFKA_WEBSOCKET_SERVER_TOPICS")
	}
	if env.BrokerURL != env_map["KAFKA_WEBSOCKET_SERVER_BROKER_URL"] {
		t.Errorf("Invalid value for env variable: KAFKA_WEBSOCKET_SERVER_BROKER_URL")
	}
	if env.Port != env_map["KAFKA_WEBSOCKET_SERVER_PORT"] {
		t.Errorf("Invalid value for env variable: KAFKA_WEBSOCKET_SERVER_PORT")
	}
	if env.Prefix != env_map["KAFKA_WEBSOCKET_SERVER_PREFIX"] {
		t.Errorf("Invalid value for env variable: KAFKA_WEBSOCKET_SERVER_PREFIX")
	}
	if env.HealthPort != env_map["KAFKA_WEBSOCKET_SERVER_HEALTH_PORT"] {
		t.Errorf("Invalid value for env variable: KAFKA_WEBSOCKET_SERVER_HEALTH_PORT")
	}
}
