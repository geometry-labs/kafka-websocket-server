package websockets

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
)

type KafkaWebsocketServer struct {
	TopicChans map[string]chan *kafka.Message
	Port       string
}

func (ws *KafkaWebsocketServer) ListenAndServe() {

	for t, c := range ws.TopicChans {

		http.HandleFunc("/"+t, readAndBroadcastKafkaTopic(c))
	}

	log.Fatal(http.ListenAndServe(":"+ws.Port, nil))
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func readAndBroadcastKafkaTopic(topic_chan chan *kafka.Message) func(w http.ResponseWriter, r *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {
		c, err := upgrader.Upgrade(w, r, nil)

		if err != nil {
			log.Print("upgrade:", err)
			return
		}
		defer c.Close()

		// Read for close
		client_close_sig := make(chan bool)
		go func() {
			for {
				_, _, err := c.ReadMessage()
				if err != nil {
					client_close_sig <- true
					break
				}
			}
		}()

		for {
			// Read
			msg := <-topic_chan

			// Broadcast
			err = c.WriteMessage(websocket.TextMessage, msg.Value)
			if err != nil {
				break
			}

			// check for client close
			select {
			case _ = <-client_close_sig:
				break
			default:
				continue
			}
		}
	}
}

// Icon ETL spacific websocket handler
func filterAndBroadcastKafkaTopic(topic_chan chan *kafka.Message) func(w http.ResponseWriter, r *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {
		c, err := upgrader.Upgrade(w, r, nil)

		if err != nil {
			log.Print("upgrade:", err)
			return
		}
		defer c.Close()

		// Registration id given by the resigtration service
		// Use for filtering topic messages
		registration_id := ""

		// Read registration
		client_close_sig := make(chan bool)
		go func() {
			for {
				_, json_data, err := c.ReadMessage()

				if err != nil {
					_ = c.WriteMessage(websocket.TextMessage, []byte("Closing websocket"))

					client_close_sig <- true
					break
				}

				// Make registration POST request
				resp, err := http.Post(
					"registration/register/broadcaster", // NOTE: should be "/broadcaster/register"
					"application/json",
					bytes.NewBuffer(json_data),
				)
				if err != nil || resp.StatusCode != 200 {
					_ = c.WriteMessage(websocket.TextMessage, []byte("Error registering events...closing websocket"))

					client_close_sig <- true
					break
				}

				var registration_json map[string]interface{}
				json.NewDecoder(resp.Body).Decode(&registration_json)

				registration_id = registration_json["registration_id"].(string)

				// unregister on disconnect
				defer func() {
					unregister_existing_blob_format := `
						{
							registration_id: %s
						}
					`
					unregister_existing_blob := fmt.Sprintf(unregister_existing_blob_format, registration_id)

					_, _ = http.Post(
						"registration/unregister/broadcaster", // NOTE: should be "/broadcaster/unregister"
						"application/json",
						bytes.NewBuffer([]byte(unregister_existing_blob)),
					)
				}()
			}
		}()

		for {
			// Read
			msg := <-topic_chan

			// TODO filter registration_id

			// Broadcast
			err = c.WriteMessage(websocket.TextMessage, msg.Value)
			if err != nil {
				break
			}

			// check for client close
			select {
			case _ = <-client_close_sig:
				break
			default:
				continue
			}
		}
	}
}
