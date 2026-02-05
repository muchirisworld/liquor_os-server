package repository

import (
	"github.com/All-Things-Muchiri/server/internal/domain"
	"github.com/jmoiron/sqlx"
)

type OrganizationsRepository struct {
	db *sqlx.DB
}

func NewOrganizationsRepository(db *sqlx.DB) *OrganizationsRepository {
	return &OrganizationsRepository{
		db: db,
	}
}

func (r *OrganizationsRepository) CreateOrganization(organization *domain.OrganizationRequest) (*domain.Organization, error) {
	query := `
		INSERT INTO organizations (
			id, name, slug, created_by, logo_url, image_url, metadata, max_allowed_memberships
		)
		VALUES (
			:id, :name, :slug, :created_by, :logo_url, :image_url, :metadata, :max_allowed_memberships
		)
		RETURNING id, name, slug, created_by, logo_url, image_url, metadata, max_allowed_memberships, created_at, updated_at
	`

	rows, err := r.db.NamedQuery(query, organization)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var createdOrg domain.Organization
	if rows.Next() {
		err = rows.StructScan(&createdOrg)
		if err != nil {
			return nil, err
		}
	}

	return &createdOrg, nil
}
