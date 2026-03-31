package service

import (
	"API_Server/internal/model"
	"API_Server/internal/repository"
	"context"
	"time"

	"github.com/google/uuid"
)

type CandidateService struct {
	candidateRepo *repository.CandidateRepository
}

func NewCandidateService(candidateRepo *repository.CandidateRepository) *CandidateService {
	return &CandidateService{candidateRepo: candidateRepo}
}

func (s *CandidateService) RegisterCandidate(ctx context.Context, sessionID, userID string, req model.RegisterCandidateRequest) (*model.Candidate, error) {
	candidate := &model.Candidate{
		CandidateID: uuid.New().String(),
		SessionID:   sessionID,
		UserID:      userID,
		Type:        req.Type,
		IP:          req.IP,
		Port:        req.Port,
		Protocol:    req.Protocol,
		CreatedAt:   time.Now().UTC(),
	}

	if err := s.candidateRepo.Create(ctx, candidate); err != nil {
		return nil, err
	}
	return candidate, nil
}

func (s *CandidateService) ListCandidates(ctx context.Context, sessionID string) ([]model.Candidate, error) {
	return s.candidateRepo.FindBySessionID(ctx, sessionID)
}
