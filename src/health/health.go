package health

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
)

type HealthChecker struct {
	ws_port   string
	ws_prefix string
	topics    []string
}

var healthChecker HealthChecker

func ListenAndServe(port string, ws_port string, ws_prefix string, topics []string) {
	healthChecker = HealthChecker{
		ws_port,
		ws_prefix,
		topics,
	}

	http.HandleFunc("/", healthChecker.generalCheckHealth)

	http.ListenAndServe(":"+port, nil)
}

func (h *HealthChecker) generalCheckHealth(w http.ResponseWriter, r *http.Request) {

	// Check websocket server
	for _, t := range h.topics {
		url := fmt.Sprintf("ws://localhost:%s%s/%s", h.ws_port, h.ws_prefix, t)

		websocket_client, _, err := websocket.DefaultDialer.Dial(url, nil)
		if err != nil {
			err_msg := fmt.Sprintf(`{"err":"Failed to connect to websocket server on /%s"}`, t)

			w.WriteHeader(500)
			w.Write([]byte(err_msg))
			return
		}
		defer websocket_client.Close()
	}

	w.WriteHeader(200)
	w.Write([]byte("{}"))
}
