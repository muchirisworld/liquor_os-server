package service

import (
	"context"

	"github.com/All-Things-Muchiri/server/internal/domain"
	"github.com/All-Things-Muchiri/server/internal/repository"
)

type UserService struct {
	userRepo repository.UserRepository
}

func (s *UserService) CreateUser(ctx context.Context, userRequest *domain.UserRequest) (*domain.User, error) {
	usr, err := s.userRepo.CreateUser(ctx, userRequest)
	if err != nil {
		return nil, err
	}
	
	return usr, nil
}