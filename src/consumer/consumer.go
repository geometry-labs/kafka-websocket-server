package consumer

import (
	"log"
	"time"

	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
)

type KafkaTopicConsumer struct {
	TopicName string
	TopicChan chan *kafka.Message

	BrokerURL string
}

func (k *KafkaTopicConsumer) ConsumeAndBroadcastTopics() {

	var consumer *kafka.Consumer

	retries := 0
	max_retries := 60
	for {
		var err error
		consumer, err = kafka.NewConsumer(&kafka.ConfigMap{
			"bootstrap.servers": k.BrokerURL,
			"group.id":          "websocket-api-group",
			"auto.offset.reset": "latest",
		})

		if err != nil {
			retries++

			// Failed to connect
			if retries >= max_retries {
				log.Printf("Consumer: ERROR - Failed to connect to kafka broker after %d retries", max_retries)
				panic(err)
			}

			log.Printf("Consumer: Could not connect to kafka @ %s, %d/%d retries...", k.BrokerURL, retries, max_retries)

			time.Sleep(1 * time.Second)
			continue
		}

		defer consumer.Close()
		break
	}

	consumer.SubscribeTopics([]string{k.TopicName}, nil)
	log.Printf("Consumer: Reading from topic %s...", k.TopicName)

	for {
		msg, err := consumer.ReadMessage(-1)
		if err == nil {

			// Select for non-blocking channels
			select {
			case k.TopicChan <- msg:
			default:
			}
		}
	}
}
