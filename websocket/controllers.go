package websocket

import (
	"encoding/json"
	"fmt"
	"go-backend/database"
	"go-backend/models"
	"log"

	"github.com/gorilla/websocket"
)

func (h *Hub) Register(userID string, conn *websocket.Conn) {
	h.mu.Lock()
	h.clients[userID] = conn
	h.mu.Unlock()
	fmt.Println("🔗 User connected:", userID)
}

func (h *Hub) Unregister(userID string) {
	h.mu.Lock()
	delete(h.clients, userID)
	h.mu.Unlock()
	fmt.Println("❌ User disconnected:", userID)
}

// Send a message to a user using the websocket connection and fallback storage
func (h *Hub) Send(senderId string, msg []byte) (err error) {
	var wsMsg models.Message
	if err := json.Unmarshal(msg, &wsMsg); err != nil {
		log.Println("Invalid message:", err)
		return err
	}
	wsMsg.SenderID = senderId

	// Broadcast au receiver
	outgoing, _ := json.Marshal(wsMsg)
	h.mu.RLock()
	defer h.mu.RUnlock()
	database.DB.Create(&wsMsg)

	if conn, ok := h.clients[wsMsg.ReceiverID]; ok {
		conn.WriteMessage(websocket.TextMessage, outgoing)
	}
	fmt.Printf("💬 Message from %s to %s: %s\n", wsMsg.SenderID, wsMsg.ReceiverID, wsMsg.Content)
	return nil
}
