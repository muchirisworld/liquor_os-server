package service

import (
	"context"

	"github.com/All-Things-Muchiri/server/internal/domain"
	"github.com/All-Things-Muchiri/server/internal/repository"
)

type OrganizationsService struct {
	orgsRepo *repository.OrganizationsRepository
}

func NewOrganizationsService(repo *repository.OrganizationsRepository) *OrganizationsService {
	return &OrganizationsService{
		orgsRepo: repo,
	}
}

func (s *OrganizationsService) CreateOrganization(ctx context.Context, organization *domain.OrganizationRequest) error {
	_, err := s.orgsRepo.CreateOrganization(organization)
	if err != nil {
		return err
	}

	return nil
}
