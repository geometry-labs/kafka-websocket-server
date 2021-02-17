package main

import (
	"log"
	"os"
	"strings"

	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"

	"kafka-websocket-server/consumer"
	"kafka-websocket-server/websockets"
)

func main() {

	topics_env := os.Getenv("KAFKA_WEBSOCKET_SERVER_TOPICS")
	broker_url_env := os.Getenv("KAFKA_WEBSOCKET_SERVER_BROKER_URL")
	port_env := os.Getenv("KAFKA_WEBSOCKET_SERVER_PORT")
	prefix_env := os.Getenv("KAFKA_WEBSOCKET_SERVER_PREFIX")

	if topics_env == "" {
		log.Println("ERROR: required enviroment variable missing: WEBSOCKET_API_TOPICS")
		return
	}
	if broker_url_env == "" {
		log.Println("ERROR: required enviroment variable missing: WEBSOCKET_API_BROKER_URL")
		return
	}
	if port_env == "" {
		port_env = "8080"
	}
	if prefix_env == "" {
		prefix_env = ""
	}

	topic_names := strings.Split(topics_env, ",")
	topic_chans := make(map[string]chan *kafka.Message)

	for _, topic := range topic_names {
		// Create channel
		topic_chans[topic] = make(chan *kafka.Message)

		// Create consumer
		kafka_consumer := consumer.KafkaTopicConsumer{
			topic,
			topic_chans[topic],
			broker_url_env,
		}

		// Start consumer
		go kafka_consumer.ConsumeAndBroadcastTopics()
		log.Printf("Kafka consumer create for %s", topic)
	}

	// Create server
	websocket_server := websockets.KafkaWebsocketServer{
		topic_chans,
		port_env,
		prefix_env,
	}

	// Start server
	go websocket_server.ListenAndServe()
	log.Printf("Websocket server listening on :%s%s/...", port_env, prefix_env)

	// Keep main thread alive
	for {
	}
}
