package service

import (
	"API_Server/internal/model"
	"API_Server/internal/repository"
	"context"
)

type UserService struct {
	userRepo *repository.UserRepository
}

func NewUserService(userRepo *repository.UserRepository) *UserService {
	return &UserService{userRepo: userRepo}
}

func (s *UserService) GetProfile(ctx context.Context, userID string) (*model.UserProfile, error) {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil || user == nil {
		return nil, err
	}

	return &model.UserProfile{
		ID:           user.ID,
		Email:        user.Email,
		Nickname:     user.Nickname,
		ProfileImage: user.ProfileImage,
		CreatedAt:    user.CreatedAt,
		AgentStatus:  "offline", // TODO: 실제 에이전트 상태 조회
	}, nil
}

func (s *UserService) UpdateProfile(ctx context.Context, userID string, req model.UpdateProfileRequest) (*model.UserProfile, error) {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil || user == nil {
		return nil, err
	}

	if req.Nickname != nil {
		user.Nickname = *req.Nickname
	}
	// profileImage: 요청에 포함되면 업데이트 (null이면 삭제)
	if req.ProfileImage != nil {
		user.ProfileImage = req.ProfileImage
	}

	if err := s.userRepo.Update(ctx, user); err != nil {
		return nil, err
	}

	return &model.UserProfile{
		ID:           user.ID,
		Email:        user.Email,
		Nickname:     user.Nickname,
		ProfileImage: user.ProfileImage,
		CreatedAt:    user.CreatedAt,
		AgentStatus:  "offline",
	}, nil
}
