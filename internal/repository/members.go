package repository

import (
	"context"

	"github.com/All-Things-Muchiri/server/internal/domain"
	"github.com/jmoiron/sqlx"
)

type MembershipRepository struct {
	db *sqlx.DB
}

func NewMembershipRepository(db *sqlx.DB) *MembershipRepository {
	return &MembershipRepository{
		db: db,
	}
}

func (r *MembershipRepository) CreateMember(ctx context.Context, req *domain.MembershipRequest) (*domain.Member, error) {
    query := `
        WITH inserted AS (
            INSERT INTO members (user_id, organization_id, identifier, image_url, profile_image_url, role, role_name)
            VALUES (:user_id, :organization_id, :identifier, :image_url, :profile_image_url, :role, :role_name)
            RETURNING *
        )
        SELECT 
            i.user_id,
            split_part(u.name, ' ', 1) as first_name,
            split_part(u.name, ' ', 2) as last_name,
            i.identifier,
            i.image_url,
            i.profile_image_url
        FROM inserted i
        JOIN users u ON i.user_id = u.id
    `

    var member domain.Member
    rows, err := r.db.NamedQueryContext(ctx, query, req)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    if rows.Next() {
        err = rows.StructScan(&member)
    }
    
    return &member, err
}
