package service

import (
	"API_Server/internal/model"
	"API_Server/internal/repository"
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrEmailExists   = errors.New("이미 존재하는 이메일입니다")
	ErrInvalidLogin  = errors.New("이메일 또는 비밀번호가 일치하지 않습니다")
	ErrWrongPassword = errors.New("현재 비밀번호가 일치하지 않습니다")
)

type AuthService struct {
	userRepo  *repository.UserRepository
	tokenRepo *repository.TokenRepository
}

func NewAuthService(userRepo *repository.UserRepository, tokenRepo *repository.TokenRepository) *AuthService {
	return &AuthService{userRepo: userRepo, tokenRepo: tokenRepo}
}

func (s *AuthService) Register(ctx context.Context, req model.RegisterRequest) (*model.RegisterResponse, error) {
	// 이메일 중복 체크
	existing, err := s.userRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, ErrEmailExists
	}

	// 비밀번호 해시
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &model.User{
		ID:           uuid.New().String(),
		Email:        req.Email,
		PasswordHash: string(hash),
		Nickname:     req.Nickname,
		CreatedAt:    time.Now().UTC(),
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	return &model.RegisterResponse{
		UserID:   user.ID,
		Email:    user.Email,
		Nickname: user.Nickname,
	}, nil
}

func (s *AuthService) Login(ctx context.Context, req model.LoginRequest) (*model.LoginResponse, error) {
	user, err := s.userRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrInvalidLogin
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, ErrInvalidLogin
	}

	// opaque 토큰 생성
	tokenBytes := make([]byte, 32)
	if _, err := rand.Read(tokenBytes); err != nil {
		return nil, err
	}
	token := hex.EncodeToString(tokenBytes)
	expiresAt := time.Now().UTC().Add(24 * time.Hour)

	authToken := &model.AuthToken{
		Token:     token,
		UserID:    user.ID,
		ExpiresAt: expiresAt,
	}
	if err := s.tokenRepo.Create(ctx, authToken); err != nil {
		return nil, err
	}

	return &model.LoginResponse{
		Token:     token,
		ExpiresAt: expiresAt,
	}, nil
}

func (s *AuthService) Logout(ctx context.Context, token string) error {
	return s.tokenRepo.DeleteByToken(ctx, token)
}

func (s *AuthService) DeleteAccount(ctx context.Context, userID string) error {
	// 모든 토큰 삭제 → 사용자 삭제
	if err := s.tokenRepo.DeleteByUserID(ctx, userID); err != nil {
		return err
	}
	return s.userRepo.Delete(ctx, userID)
}

func (s *AuthService) ChangePassword(ctx context.Context, userID string, req model.ChangePasswordRequest) error {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil || user == nil {
		return ErrInvalidLogin
	}

	// 현재 비밀번호 확인
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.CurrentPassword)); err != nil {
		return ErrWrongPassword
	}

	// 새 비밀번호 해시
	hash, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	return s.userRepo.UpdatePassword(ctx, userID, string(hash))
}
