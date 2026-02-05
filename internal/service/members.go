package service

import (
    "context"
    "github.com/All-Things-Muchiri/server/internal/domain"
    "github.com/All-Things-Muchiri/server/internal/repository"
)

type MembershipService struct {
    membersRepo *repository.MembershipRepository
}

func NewMembershipService(membersRepo *repository.MembershipRepository) *MembershipService {
    return &MembershipService{
        membersRepo: membersRepo,
    }
}

func (r *MembershipService) CreateMember(ctx context.Context, req *domain.MembershipRequest) error {
    _, err := r.membersRepo.CreateMember(ctx, req)
    if err != nil {
        return err
    }
    
    return nil
}