package domain

import (
	"encoding/json"
	"time"
)

type Organization struct {
	OrganizationID        string          `json:"id" db:"id"`
	Name                  string          `json:"name" db:"name"`
	Slug                  string          `json:"slug" db:"slug"`
	CreatedBy             string          `json:"created_by" db:"created_by"`
	LogoURL               string          `json:"logo_url" db:"logo_url"`
	ImageURL              string          `json:"image_url" db:"image_url"`
	Metadata              json.RawMessage `json:"metadata" db:"metadata"`
	MaxAllowedMemberships int32           `json:"max_allowed_memberships" db:"max_allowed_memberships"`
	CreatedAt             time.Time       `json:"created_at" db:"created_at"`
	UpdatedAt             time.Time       `json:"updated_at" db:"updated_at"`
}

type OrganizationRequest struct {
	OrganizationID        string          `json:"id" db:"id"`
	Name                  string          `json:"name" db:"name"`
	Slug                  string          `json:"slug" db:"slug"`
	CreatedBy             string          `json:"created_by" db:"created_by"`
	LogoURL               string          `json:"logo_url" db:"logo_url"`
	ImageURL              string          `json:"image_url" db:"image_url"`
	Metadata              json.RawMessage `json:"metadata" db:"metadata"`
	MaxAllowedMemberships int32           `json:"max_allowed_memberships" db:"max_allowed_memberships"`
}