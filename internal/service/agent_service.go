package service

import (
	"API_Server/internal/model"
	"API_Server/internal/repository"
	"context"
	"time"

	"github.com/google/uuid"
)

type AgentService struct {
	agentRepo *repository.AgentRepository
}

func NewAgentService(agentRepo *repository.AgentRepository) *AgentService {
	return &AgentService{agentRepo: agentRepo}
}

func (s *AgentService) RegisterAgent(ctx context.Context, userID string, req model.RegisterAgentRequest) (*model.Agent, error) {
	now := time.Now().UTC()
	agent := &model.Agent{
		AgentID:      uuid.New().String(),
		DeviceName:   req.DeviceName,
		Platform:     req.Platform,
		AgentVersion: req.AgentVersion,
		MultiAddress: req.MultiAddress,
		RegisteredAt: now,
		LastSeenAt:   now,
		Status:       "offline",
	}

	if err := s.agentRepo.Create(ctx, userID, agent); err != nil {
		return nil, err
	}
	return agent, nil
}

func (s *AgentService) ListAgents(ctx context.Context, userID string) ([]model.Agent, error) {
	return s.agentRepo.FindByUserID(ctx, userID)
}
