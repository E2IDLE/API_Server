package service

import (
	"API_Server/internal/model"
	"API_Server/internal/repository"
	"context"
	"crypto/rand"
	"errors"
	"log"
	"math/big"
	"time"

	"github.com/google/uuid"
)

var (
	ErrSessionNotFound = errors.New("세션을 찾을 수 없습니다")
	ErrForbidden       = errors.New("접근 권한이 없습니다")
	ErrInvalidInvite   = errors.New("잘못된 초대 코드입니다")
)

type SessionService struct {
	sessionRepo *repository.SessionRepository
}

func NewSessionService(sessionRepo *repository.SessionRepository) *SessionService {
	return &SessionService{sessionRepo: sessionRepo}
}

func (s *SessionService) CreateSession(ctx context.Context, userID, token string) (*model.Session, error) {
	now := time.Now().UTC()
	session := &model.Session{
		SessionID:   uuid.New().String(),
		InviteCode:  generateInviteCode(6),
		Status:      "waiting",
		SenderID:    userID,
		ReceiverID:  nil,
		SenderToken: token,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if err := s.sessionRepo.Create(ctx, session); err != nil {
		return nil, err
	}
	return session, nil
}

func (s *SessionService) JoinSession(ctx context.Context, sessionID, userID, token string, req model.JoinSessionRequest) (*model.Session, error) {
	session, err := s.sessionRepo.FindByID(ctx, sessionID)
	if err != nil {
		return nil, err
	}
	if session == nil {
		return nil, ErrSessionNotFound
	}
	if session.InviteCode != req.InviteCode {
		return nil, ErrInvalidInvite
	}

	if err := s.sessionRepo.UpdateReceiver(ctx, sessionID, userID, token, "connecting"); err != nil {
		return nil, err
	}

	// 업데이트된 세션 반환
	return s.sessionRepo.FindByID(ctx, sessionID)
}

func (s *SessionService) GetSession(ctx context.Context, sessionID, userID string) (*model.Session, error) {
	session, err := s.sessionRepo.FindByID(ctx, sessionID)
	if err != nil {
		return nil, err
	}
	if session == nil {
		return nil, ErrSessionNotFound
	}

	// 접근 권한 확인
	if session.SenderID != userID && (session.ReceiverID == nil || *session.ReceiverID != userID) {
		return nil, ErrForbidden
	}
	return session, nil
}

func (s *SessionService) DeleteSession(ctx context.Context, sessionID, userID string) error {
	session, err := s.sessionRepo.FindByID(ctx, sessionID)
	if err != nil {
		return err
	}
	if session == nil {
		return ErrSessionNotFound
	}
	if session.SenderID != userID {
		return ErrForbidden
	}
	return s.sessionRepo.Delete(ctx, sessionID)
}

func (s *SessionService) GetHistory(ctx context.Context, userID string, page, pageSize int) (*model.SessionHistoryResponse, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	sessions, totalCount, err := s.sessionRepo.FindByUserIDPaginated(ctx, userID, page, pageSize)
	if err != nil {
		return nil, err
	}

	totalPages := 0
	if totalCount > 0 {
		totalPages = (totalCount + pageSize - 1) / pageSize
	}

	return &model.SessionHistoryResponse{
		Items: sessions,
		Pagination: model.Pagination{
			Page:       page,
			PageSize:   pageSize,
			TotalCount: totalCount,
			TotalPages: totalPages,
		},
	}, nil
}

// generateInviteCode 는 6자리 영숫자 초대 코드를 생성
func generateInviteCode(length int) string {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	code := make([]byte, length)
	for i := range code {
		n, _ := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		code[i] = charset[n.Int64()]
	}
	return string(code)
}

func (s *SessionService) destroySession(sessionID string) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	log.Printf("세션 자동 소멸 시작: %s", sessionID)
	s.sessionRepo.Delete(ctx, sessionID)
	log.Printf("세션 자동 소멸 완료: %s", sessionID)
}

// ★ 세션 상태 변경 + 자동 소멸 ★
func (s *SessionService) UpdateStatus(ctx context.Context, sessionID, status string) error {
	if err := s.sessionRepo.UpdateStatus(ctx, sessionID, status); err != nil {
		return err
	}

	// "completed" 또는 "error" → 자동 소멸
	if status == "completed" || status == "error" {
		go s.destroySession(sessionID)
	}

	return nil
}
