package main

import (
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"

	"kafka-websocket-server/consumer"
	"kafka-websocket-server/health"
	"kafka-websocket-server/websockets"
)

func main() {

	env := getEnvironment()

	topic_names := strings.Split(env.Topics, ",")
	broadcasters := make(map[string]*websockets.TopicBroadcaster)

	// One consumer per topic
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
			env.BrokerURL,
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
		env.Port,
		env.Prefix,
	}

	// Start server
	go websocket_server.ListenAndServe()
	log.Printf("Websocket server listening on :%s%s/...", env.Port, env.Prefix)

	// Start health server
	go health.ListenAndServe(env.HealthPort, env.Port, env.Prefix, topic_names)

	// Listen for close sig
	sigCh := make(chan os.Signal)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM, syscall.SIGKILL)

	// Keep main thread alive
	for sig := range sigCh {
		log.Printf("Stopping websocket server...%s", sig.String())
	}
}
