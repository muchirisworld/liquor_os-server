package repository

import (
	"context"

	"github.com/All-Things-Muchiri/server/internal/domain"
	"github.com/jmoiron/sqlx"
)

type UserRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (r *UserRepository) CreateUser(ctx context.Context, userRequest *domain.UserRequest) (*domain.User, error) {
	rows, err := r.db.NamedQueryContext(ctx,
		`INSERT INTO "user" (
			 id, name, email, email_verified, image
		) VALUES (
			:id, :name, :email, :email_verified, :image
		) RETURNING *`,
		userRequest,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var user domain.User
	if rows.Next() {
		if err := rows.StructScan(&user); err != nil {
			return nil, err
		}
	}
	
	return &user, nil
}