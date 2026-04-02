package handler

//WebSocket 핸들러는 클라이언트와의 WebSocket 연결을 관리하고, 메시지를 주고받는 역할을 합니다.

import (
	"API_Server/internal/model"
	"API_Server/internal/repository"
	"API_Server/internal/service"
	"API_Server/internal/ws"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

type WSHandler struct {
	hub        *ws.Hub
	tokenRepo  *repository.TokenRepository
	sessionSvc *service.SessionService
}

func NewWSHandler(hub *ws.Hub, tokenRepo *repository.TokenRepository, sessionSvc *service.SessionService) *WSHandler {
	return &WSHandler{hub: hub, tokenRepo: tokenRepo, sessionSvc: sessionSvc}
}

func (h *WSHandler) HandleWebSocket(c *gin.Context) {
	token := c.Query("token")
	if token == "" {
		c.JSON(http.StatusUnauthorized, model.ErrorResponse{Code: "UNAUTHORIZED", Message: "토큰이 필요합니다."})
		return
	}

	authToken, err := h.tokenRepo.FindByToken(c.Request.Context(), token)
	if err != nil || authToken == nil {
		c.JSON(http.StatusUnauthorized, model.ErrorResponse{Code: "UNAUTHORIZED", Message: "유효하지 않은 토큰입니다."})
		return
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("WebSocket 업그레이드 실패: %v", err)
		return
	}

	client := &ws.Client{
		UserID: authToken.UserID,
		Conn:   conn,
		Send:   make(chan []byte, 256),
	}

	h.hub.Register(client)

	go h.writePump(client)
	go h.readPump(client)
}

func (h *WSHandler) writePump(client *ws.Client) {
	defer func() {
		client.Conn.Close()
	}()

	for {
		msg, ok := <-client.Send
		if !ok {
			client.Conn.WriteMessage(websocket.CloseMessage, []byte{})
			return
		}

		client.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
		if err := client.Conn.WriteMessage(websocket.TextMessage, msg); err != nil {
			return
		}
	}
}

func (h *WSHandler) readPump(client *ws.Client) {
	defer func() {
		h.hub.Unregister(client)
		client.Conn.Close()
	}()

	client.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	client.Conn.SetPongHandler(func(string) error {
		client.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	for {
		_, message, err := client.Conn.ReadMessage()
		if err != nil {
			break
		}

		var wsMsg model.WSMessage
		if err := json.Unmarshal(message, &wsMsg); err != nil {
			continue
		}

		switch wsMsg.Event {
		case "agent:ping":
			client.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))
			log.Printf("agent:ping from user %s", client.UserID)

		case "session:status_changed":
			if dataMap, ok := wsMsg.Data.(map[string]interface{}); ok {
				sessionID, _ := dataMap["sessionId"].(string)
				status, _ := dataMap["status"].(string)

				if status == "completed" || status == "error" {
					go h.sessionSvc.UpdateStatus(
						context.Background(), sessionID, status, // ← 쉼표 추가
					)
				}
			}
		}
	}
}
