package websockets

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
)

type KafkaWebsocketServer struct {
	Broadcasters map[string]*TopicBroadcaster
	Port         string
	Prefix       string
}

func (ws *KafkaWebsocketServer) ListenAndServe() {

	for t, b := range ws.Broadcasters {

		endpoint_path := ws.Prefix + "/" + t

		http.HandleFunc(endpoint_path, readAndBroadcastKafkaTopic(b))
	}

	log.Fatal(http.ListenAndServe(":"+ws.Port, nil))
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func readAndBroadcastKafkaTopic(broadcaster *TopicBroadcaster) func(w http.ResponseWriter, r *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {
		c, err := upgrader.Upgrade(w, r, nil)

		if err != nil {
			log.Print("upgrade:", err)
			return
		}
		defer c.Close()

		// Add broadcaster
		topic_chan := make(chan *kafka.Message)
		id := broadcaster.AddWebsocketChannel(topic_chan)
		defer func() {
			// Remove broadcaster
			broadcaster.RemoveWebsocketChannel(id)
		}()

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
