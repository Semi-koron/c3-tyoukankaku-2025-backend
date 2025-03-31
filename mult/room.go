package mult

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

type SendData struct {
	ID       string  `json:"id"`
	Color    string  `json:"color"`
	Position struct {
		X float64 `json:"x"`
		Y float64 `json:"y"`
		Z float64 `json:"z"`
	} `json:"position"`
	Rotation struct {
		X float64 `json:"x"`
		Y float64 `json:"y"`
		Z float64 `json:"z"`
		W float64 `json:"w"`
	} `json:"rotation"`
	Chat string `json:"chat"`
}

type ReceiveData []SendData

var (
	clients   = make(map[*websocket.Conn]bool)
	dataStore = make(ReceiveData, 0)
	lock      sync.Mutex
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func handleConnection(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Failed to upgrade connection:", err)
		return
	}
	defer conn.Close()
	clients[conn] = true

	for {
		var msg SendData
		err := conn.ReadJSON(&msg)
		if err != nil {
			fmt.Println("Error reading message:", err)
			delete(clients, conn)
			break
		}

		lock.Lock()
		dataStore = append(dataStore, msg)
		lock.Unlock()

		sendToAllClients()
	}
}

func sendToAllClients() {
	lock.Lock()
	data, err := json.Marshal(dataStore)
	lock.Unlock()
	if err != nil {
		fmt.Println("Error marshalling data:", err)
		return
	}

	for client := range clients {
		err := client.WriteMessage(websocket.TextMessage, data)
		if err != nil {
			fmt.Println("Error sending message:", err)
			client.Close()
			delete(clients, client)
		}
	}
}

func main() {
	http.HandleFunc("/ws", handleConnection)
	fmt.Println("WebSocket server started on ws://localhost:8080/ws")
	http.ListenAndServe(":8080", nil)
}
