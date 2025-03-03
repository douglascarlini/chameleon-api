package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

func wsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket upgrade error:", err)
		return
	}

	clientID := uuid.New().String()
	mutex.Lock()
	connections[clientID] = conn
	mutex.Unlock()

	log.Printf("new ws client: %s\n", clientID)

	conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("Your ID: %s", clientID)))

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Printf("ws read error (Client %s): %v\n", clientID, err)
			mutex.Lock()
			delete(connections, clientID)
			mutex.Unlock()
			conn.Close()
			break
		}

		log.Printf("new ws message from %s: %s\n", clientID, msg)
	}
}

func sendMessageToClient(clientID string, msg any) error {
	mutex.Lock()
	conn, exists := connections[clientID]
	mutex.Unlock()

	message, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	if exists {
		err := conn.WriteMessage(websocket.TextMessage, []byte(message))
		if err != nil {
			conn.Close()
			mutex.Lock()
			delete(connections, clientID)
			mutex.Unlock()
			return err
		}
	} else {
		return fmt.Errorf("client %s not found", clientID)
	}

	return nil

}

func sendBroadcastMessage(message string) {
	mutex.Lock()
	defer mutex.Unlock()

	for id, conn := range connections {
		err := conn.WriteMessage(websocket.TextMessage, []byte(message))
		if err != nil {
			conn.Close()
			delete(connections, id)
			log.Printf("ws write error [%s]: %v\n", id, err)
		}
	}
}
