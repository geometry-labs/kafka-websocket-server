package consumer

import (
	"testing"
	"time"

	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
)

func TestKafkaTopicConsumer(t *testing.T) {

	// create test producer
	go func() {
		p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": "kafka:9092"})
		if err != nil {
			t.Logf("Failed to connect to kafka broker")
			t.Fail()
		}

		defer p.Close()

		topic := "test_topic"
		for {
			p.Produce(&kafka.Message{
				TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
				Value:          []byte("test_message"),
			}, nil)

			time.Sleep(1 * time.Second)
		}
	}()

	topic_chan := make(chan *kafka.Message)

	topic_consumer := KafkaTopicConsumer{
		"test_topic",
		topic_chan,
		"kafka:9092",
	}

	go topic_consumer.ConsumeAndBroadcastTopics()

	select {
	case res := <-topic_chan:
		msg := string(res.Value)
		if msg != "test_message" {
			t.Logf("Failed to assert topic message value")
			t.Fail()
		}
	case <-time.After(10 * time.Second):
		t.Logf("Failed to receive message from kafka")
		t.Fail()
	}

}
