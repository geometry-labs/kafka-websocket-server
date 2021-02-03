package main

import (
  "os"
  "strings"

  "gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"

  "websocket-api/consumer"
  "websocket-api/websockets"
)

func main() {

  // comma seperated list of topic names to serve
  topics_env_raw := os.Getenv("WEBSOCKET_API_TOPICS")

  topic_names := strings.Split(topics_env_raw, ",")
  topic_chans := make(map[string]chan *kafka.Message)

  for _, topic := range topic_names {
    // Create channel
    topic_chans[topic] = make(chan *kafka.Message)

    // Create consumer
    kafka_consumer := consumer.KafkaTopicConsumer{
      topic,
      topic_chans[topic],
    }

    // Start consumer
    go kafka_consumer.ConsumeAndBroadcastTopics()
  }

  // Create server
  websocket_server := websockets.KafkaWebsocketServer{
    topic_chans,
  }

  // Start server
  go websocket_server.ListenAndServe()

  // Keep main thread alive
  for {
  }
}
