package handler

import (
	"API_Server/internal/model"
	"API_Server/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AgentHandler struct {
	agentSvc *service.AgentService
}

func NewAgentHandler(agentSvc *service.AgentService) *AgentHandler {
	return &AgentHandler{agentSvc: agentSvc}
}

// POST /users/me/agents
func (h *AgentHandler) RegisterAgent(c *gin.Context) {
	var req model.RegisterAgentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{Code: "BAD_REQUEST", Message: err.Error()})
		return
	}

	userID, _ := c.Get("userID")
	agent, err := h.agentSvc.RegisterAgent(c.Request.Context(), userID.(string), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{Code: "INTERNAL", Message: "서버 오류"})
		return
	}

	c.JSON(http.StatusCreated, agent)
}

// GET /users/me/agents
func (h *AgentHandler) ListAgents(c *gin.Context) {
	userID, _ := c.Get("userID")

	agents, err := h.agentSvc.ListAgents(c.Request.Context(), userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{Code: "INTERNAL", Message: "서버 오류"})
		return
	}

	if agents == nil {
		agents = []model.Agent{}
	}
	c.JSON(http.StatusOK, agents)
}
