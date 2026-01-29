package repository

import (
	"context"
	"database/sql"

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
		`INSERT INTO users (
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
	if !rows.Next() {
		if err := rows.Err(); err != nil {
			return nil, err
		}
		return nil, sql.ErrNoRows
	}
	if err := rows.StructScan(&user); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) UpdateUser(ctx context.Context, userRequest *domain.UserRequest) (*domain.User, error) {
	rows, err := r.db.NamedQueryContext(ctx,
		`UPDATE users SET
			name = :name,
			email = :email,
			email_verified = :email_verified,
			image = :image,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = :id
		RETURNING *`,
		userRequest,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var user domain.User
	if !rows.Next() {
		if err := rows.Err(); err != nil {
			return nil, err
		}
		return nil, sql.ErrNoRows
	}
	if err := rows.StructScan(&user); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &user, nil
}
