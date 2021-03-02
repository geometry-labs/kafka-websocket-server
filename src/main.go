package main

import (
	"log"
	"os"
	"strings"

	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"

	"kafka-websocket-server/consumer"
	"kafka-websocket-server/health"
	"kafka-websocket-server/websockets"
)

func main() {

	topics_env := os.Getenv("KAFKA_WEBSOCKET_SERVER_TOPICS")
	broker_url_env := os.Getenv("KAFKA_WEBSOCKET_SERVER_BROKER_URL")
	port_env := os.Getenv("KAFKA_WEBSOCKET_SERVER_PORT")
	prefix_env := os.Getenv("KAFKA_WEBSOCKET_SERVER_PREFIX")
	health_port_env := os.Getenv("KAFKA_WEBSOCKET_SERVER_HEALTH_PORT")

	if topics_env == "" {
		log.Println("ERROR: required enviroment variable missing: WEBSOCKET_API_TOPICS")
		os.Exit(1)
	}
	if broker_url_env == "" {
		log.Println("ERROR: required enviroment variable missing: WEBSOCKET_API_BROKER_URL")
		os.Exit(1)
	}
	if port_env == "" {
		port_env = "8080"
	}
	if prefix_env == "" {
		prefix_env = ""
	}
	if health_port_env == "" {
		health_port_env = "5001"
	}

	topic_names := strings.Split(topics_env, ",")
	broadcasters := make(map[string]*websockets.TopicBroadcaster)

	for _, topic := range topic_names {
		// Create channel
		topic_chan := make(chan *kafka.Message)

		// Create broadcaster
		broadcasters[topic] = &websockets.TopicBroadcaster{
			topic_chan,
			make(map[websockets.BroadcasterID]chan *kafka.Message),
		}

		// Create consumer
		kafka_consumer := consumer.KafkaTopicConsumer{
			topic,
			topic_chan,
			broker_url_env,
		}

		// Start consumer
		go kafka_consumer.ConsumeAndBroadcastTopics()
		log.Printf("Kafka consumer created for %s", topic)

		// Start broadcaster
		go broadcasters[topic].Broadcast()
		log.Printf("Topic broadcaster created for %s", topic)
	}

	// Create server
	websocket_server := websockets.KafkaWebsocketServer{
		broadcasters,
		port_env,
		prefix_env,
	}

	// Start server
	go websocket_server.ListenAndServe()
	log.Printf("Websocket server listening on :%s%s/...", port_env, prefix_env)

	// Start health server
	go health.ListenAndServe(health_port_env, port_env, prefix_env, topic_names)

	// Keep main thread alive
	for {
	}
}
