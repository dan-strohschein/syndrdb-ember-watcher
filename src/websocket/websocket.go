/*
This package provides WebSocket functionality for real-time communication between this middle layer, and a web app
This allows a web app to connect to this middle layer via WebSocket protocol to receive real-time updates and visualizations of SyndrDB metrics.
*/
package websocket

import (
	"fmt"
	"net/http"
	"syndrdb-ember-watcher/src/models"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all connections for simplicity in this example
	},
}

var ws *websocket.Conn

func handleConnections(w http.ResponseWriter, r *http.Request) {
	// Upgrade initial GET request to a WebSocket
	var err error
	ws, err = upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer ws.Close() // Close the connection when the function returns

	for {
		// Read message from browser
		messageType, message, err := ws.ReadMessage()
		if err != nil {
			fmt.Println("Error reading message:", err)
			break
		}
		fmt.Printf("Received: %s\n", message)

		// Echo the message back to the client
		err = ws.WriteMessage(messageType, message)
		if err != nil {
			fmt.Println("Error writing message:", err)
			break
		}
	}
}

func StartWebSocketServer(port int) {
	http.HandleFunc("/ws", handleConnections)
	fmt.Printf("WebSocket server started on :%d\n", port)
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}

func StopWebSocketServer() {
	// Implement server shutdown logic if needed
	ws.Close()
}

func BroadcastMetrics(metrics *models.MetricsBlock) {
	if ws != nil {
		err := ws.WriteJSON(metrics)
		if err != nil {
			fmt.Println("Error broadcasting metrics:", err)
		}
	}
}
