package handler

import (
	"API_Server/internal/model"
	"API_Server/internal/service"
	"net/http"

	"API_Server/internal/ws"

	"github.com/gin-gonic/gin"
)

type CandidateHandler struct {
	candidateSvc *service.CandidateService
	hub          *ws.Hub
}

func NewCandidateHandler(candidateSvc *service.CandidateService, hub *ws.Hub) *CandidateHandler {
	return &CandidateHandler{candidateSvc: candidateSvc, hub: hub}
}

// POST /sessions/:sessionId/candidates
func (h *CandidateHandler) RegisterCandidate(c *gin.Context) {
	sessionID := c.Param("sessionId")

	var req model.RegisterCandidateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{Code: "BAD_REQUEST", Message: err.Error()})
		return
	}

	userID, _ := c.Get("userID")
	candidate, err := h.candidateSvc.RegisterCandidate(c.Request.Context(), sessionID, userID.(string), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{Code: "INTERNAL", Message: "서버 오류"})
		return
	}

	// WebSocket 알림: session:candidate_added
	h.hub.BroadcastToSession(sessionID, model.WSMessage{
		Event: "session:candidate_added",
		Data: map[string]interface{}{
			"sessionId":   sessionID,
			"candidateId": candidate.CandidateID,
			"type":        candidate.Type,
			"ip":          candidate.IP,
			"port":        candidate.Port,
			"protocol":    candidate.Protocol,
		},
	})

	c.JSON(http.StatusCreated, candidate)
}

// GET /sessions/:sessionId/candidates
func (h *CandidateHandler) ListCandidates(c *gin.Context) {
	sessionID := c.Param("sessionId")

	candidates, err := h.candidateSvc.ListCandidates(c.Request.Context(), sessionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{Code: "INTERNAL", Message: "서버 오류"})
		return
	}

	if candidates == nil {
		candidates = []model.Candidate{}
	}
	c.JSON(http.StatusOK, candidates)
}
