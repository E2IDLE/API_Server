package handler

//Turn_handler.go 는 TURN 서버 자격 증명 발급 기능을 담당하는 핸들러

import (
	"API_Server/internal/model"
	"API_Server/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type TurnHandler struct {
	turnSvc *service.TurnService
}

func NewTurnHandler(turnSvc *service.TurnService) *TurnHandler {
	return &TurnHandler{turnSvc: turnSvc}
}

// POST /sessions/:sessionId/turn-credentials
func (h *TurnHandler) IssueTurnCredentials(c *gin.Context) {
	sessionID := c.Param("sessionId")

	creds, err := h.turnSvc.IssueTurnCredentials(sessionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{Code: "INTERNAL", Message: "서버 오류"})
		return
	}

	c.JSON(http.StatusOK, creds)
}
