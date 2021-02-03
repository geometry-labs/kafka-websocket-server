package consumer

import (
  "gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
)


type KafkaTopicConsumer struct {
  Topic_name    string
  Topic_chan    chan *kafka.Message
}

func (k *KafkaTopicConsumer)ConsumeAndBroadcastTopics() {
  c, err := kafka.NewConsumer(&kafka.ConfigMap{
    "bootstrap.servers": "kafka:9092",
		"group.id":          "websocket-api-group",
		"auto.offset.reset": "latest",
	})

	if err != nil {
		panic(err)
	}
  defer c.Close()

	c.SubscribeTopics([]string{k.Topic_name}, nil)

  for {
    msg, err := c.ReadMessage(-1)
		if err == nil {

      // NOTE: use select statement for non-blocking channels
      select {
      case k.Topic_chan <- msg:
      default:
      }
    }
  }
}
