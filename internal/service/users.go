package service

import (
	"context"

	"github.com/All-Things-Muchiri/server/internal/domain"
	"github.com/All-Things-Muchiri/server/internal/repository"
)

type UserService struct {
	userRepo *repository.UserRepository
}

func NewUserService(userRepo *repository.UserRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

func (s *UserService) CreateUser(ctx context.Context, userRequest *domain.UserRequest) error {
	_, err := s.userRepo.CreateUser(ctx, userRequest)
	if err != nil {
		return err
	}

	return nil
}

func (s *UserService) UpdateUser(ctx context.Context, userRequest *domain.UserRequest) error {
	_, err := s.userRepo.UpdateUser(ctx, userRequest)
	if err != nil {
		return err
	}

	return nil
}
