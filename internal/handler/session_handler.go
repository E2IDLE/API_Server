package handler

//Session_handler.go 는 세션 생성, 참여, 조회, 삭제 및 세션 기록 조회 기능을 담당하는 핸들러입니다.

import (
	"API_Server/internal/model"
	"API_Server/internal/service"
	"API_Server/internal/ws"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type SessionHandler struct {
	sessionSvc *service.SessionService
	hub        *ws.Hub
	userSvc    *service.UserService
}

func NewSessionHandler(sessionSvc *service.SessionService, hub *ws.Hub, userSvc *service.UserService) *SessionHandler {
	return &SessionHandler{sessionSvc: sessionSvc, hub: hub, userSvc: userSvc}
}

// POST /sessions + 토큰
func (h *SessionHandler) CreateSession(c *gin.Context) {
	userID, _ := c.Get("userID")
	token, _ := c.Get("token")

	var req model.CreateSessionRequest
	_ = c.ShouldBindJSON(&req) // body 없어도 허용

	session, err := h.sessionSvc.CreateSession(c.Request.Context(), userID.(string), token.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{Code: "INTERNAL", Message: "서버 오류"})
		return
	}

	// 대상 유저가 지정된 경우 WebSocket으로 연결 신청 알림 전송
	if req.UserID != "" {
		senderProfile, err2 := h.userSvc.GetProfile(c.Request.Context(), userID.(string))
		senderNickname := ""
		if err2 == nil && senderProfile != nil {
			senderNickname = senderProfile.Nickname
		}
		h.hub.BroadcastToSession(req.UserID, model.WSMessage{
			Event: "session:connection_request",
			Data: map[string]interface{}{
				"sessionId":      session.SessionID,
				"inviteCode":     session.InviteCode,
				"senderId":       userID.(string),
				"senderNickname": senderNickname,
			},
		})
	}

	c.JSON(http.StatusCreated, session)
}

// POST /sessions/:sessionId/join +토큰
func (h *SessionHandler) JoinSession(c *gin.Context) {
	sessionID := c.Param("sessionId")

	var req model.JoinSessionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{Code: "BAD_REQUEST", Message: err.Error()})
		return
	}

	userID, _ := c.Get("userID")
	token, _ := c.Get("token")

	session, err := h.sessionSvc.JoinSession(c.Request.Context(), sessionID, userID.(string), token.(string), req)
	if err != nil {
		if errors.Is(err, service.ErrSessionNotFound) {
			c.JSON(http.StatusNotFound, model.ErrorResponse{Code: "NOT_FOUND", Message: err.Error()})
			return
		}
		if errors.Is(err, service.ErrInvalidInvite) {
			c.JSON(http.StatusBadRequest, model.ErrorResponse{Code: "BAD_REQUEST", Message: err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{Code: "INTERNAL", Message: "서버 오류"})
		return
	}

	// WebSocket 알림
	profile, err2 := h.userSvc.GetProfile(c.Request.Context(), userID.(string))
	nickname := ""
	if err2 == nil && profile != nil {
		nickname = profile.Nickname
	}
	h.hub.BroadcastToSession(session.SenderID, model.WSMessage{
		Event: "session:peer_joined",
		Data: map[string]interface{}{
			"sessionId":    session.SessionID,
			"inviteCode":   session.InviteCode,
			"peerId":       userID.(string),
			"peerNickname": nickname,
		},
	})

	c.JSON(http.StatusOK, session)
}

// GET /sessions/:sessionId
func (h *SessionHandler) GetSession(c *gin.Context) {
	sessionID := c.Param("sessionId")
	userID, _ := c.Get("userID")

	session, err := h.sessionSvc.GetSession(c.Request.Context(), sessionID, userID.(string))
	if err != nil {
		if errors.Is(err, service.ErrSessionNotFound) {
			c.JSON(http.StatusNotFound, model.ErrorResponse{Code: "NOT_FOUND", Message: err.Error()})
			return
		}
		if errors.Is(err, service.ErrForbidden) {
			c.JSON(http.StatusForbidden, model.ErrorResponse{Code: "FORBIDDEN", Message: err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{Code: "INTERNAL", Message: "서버 오류"})
		return
	}

	c.JSON(http.StatusOK, session)
}

// DELETE /sessions/:sessionId
func (h *SessionHandler) DeleteSession(c *gin.Context) {
	sessionID := c.Param("sessionId")
	userID, _ := c.Get("userID")

	err := h.sessionSvc.DeleteSession(c.Request.Context(), sessionID, userID.(string))
	if err != nil {
		if errors.Is(err, service.ErrSessionNotFound) {
			c.JSON(http.StatusNotFound, model.ErrorResponse{Code: "NOT_FOUND", Message: err.Error()})
			return
		}
		if errors.Is(err, service.ErrForbidden) {
			c.JSON(http.StatusForbidden, model.ErrorResponse{Code: "FORBIDDEN", Message: err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{Code: "INTERNAL", Message: "서버 오류"})
		return
	}

	c.Status(http.StatusNoContent)
}

// GET /sessions/history
func (h *SessionHandler) GetHistory(c *gin.Context) {
	userID, _ := c.Get("userID")

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))

	resp, err := h.sessionSvc.GetHistory(c.Request.Context(), userID.(string), page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{Code: "INTERNAL", Message: "서버 오류"})
		return
	}

	c.JSON(http.StatusOK, resp)
}
