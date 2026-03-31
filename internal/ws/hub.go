package ws

import (
	"API_Server/internal/model"
	"encoding/json"
	"log"
	"sync"

	"github.com/gorilla/websocket"
)

// Client 는 하나의 WebSocket 연결을 나타냄
type Client struct {
	UserID string
	Conn   *websocket.Conn
	Send   chan []byte
}

// Hub 은 모든 WebSocket 클라이언트를 관리
type Hub struct {
	mu         sync.RWMutex
	clients    map[string]*Client // userID → Client
	register   chan *Client
	unregister chan *Client
}

func NewHub() *Hub {
	return &Hub{
		clients:    make(map[string]*Client),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client.UserID] = client
			h.mu.Unlock()
			log.Printf("WS 클라이언트 등록: %s", client.UserID)

		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client.UserID]; ok {
				delete(h.clients, client.UserID)
				close(client.Send)
			}
			h.mu.Unlock()
			log.Printf("WS 클라이언트 해제: %s", client.UserID)
		}
	}
}

func (h *Hub) Register(client *Client) {
	h.register <- client
}

func (h *Hub) Unregister(client *Client) {
	h.unregister <- client
}

// SendToUser 는 특정 사용자에게 메시지를 전송
func (h *Hub) SendToUser(userID string, msg model.WSMessage) {
	h.mu.RLock()
	client, ok := h.clients[userID]
	h.mu.RUnlock()

	if !ok {
		return
	}

	data, err := json.Marshal(msg)
	if err != nil {
		return
	}

	select {
	case client.Send <- data:
	default:
		// 버퍼가 꽉 차면 연결 종료
		h.Unregister(client)
	}
}

// BroadcastToSession 은 특정 세션의 참가자들에게 메시지를 전송
// 간단 구현: targetUserID 로 한 명에게 전송
func (h *Hub) BroadcastToSession(targetUserID string, msg model.WSMessage) {
	h.SendToUser(targetUserID, msg)
}
