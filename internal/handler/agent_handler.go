package handler

//agent_handler.go 는 사용자 에이전트 등록 및 조회 기능을 담당하는 핸들러입니다.

import (
	"API_Server/internal/model"
	"API_Server/internal/service"
	"log"
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
	userID, exists := c.Get("userID")
	log.Printf("ListAgents 호출 - userID: %v, exists: %v", userID, exists)

	agents, err := h.agentSvc.ListAgents(c.Request.Context(), userID.(string))
	if err != nil {
		log.Printf("ListAgents error: %v", err)
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{Code: "INTERNAL", Message: "서버 오류"})
		return
	}

	log.Printf("ListAgents 결과: %d건", len(agents))
	if agents == nil {
		agents = []model.Agent{}
	}
	c.JSON(http.StatusOK, agents)
}
