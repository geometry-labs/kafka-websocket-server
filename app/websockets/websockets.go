package websockets

import (
  "log"
  "net/http"

  "gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
  "github.com/gorilla/websocket"
)


type KafkaWebsocketServer struct {
  Topic_chans map[string]chan *kafka.Message
}

func (ws *KafkaWebsocketServer)ListenAndServe() {

  for t, c := range ws.Topic_chans {

    http.HandleFunc("/" + t, readAndBroadcastKafkaTopic(c))
  }

  log.Fatal(http.ListenAndServe(":8080", nil))
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
