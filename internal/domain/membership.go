package domain

type Member struct {
	UserID          string `json:"user_id" db:"user_id"`
	Firstname       string `json:"first_name" db:"first_name"`
	Lastname        string `json:"last_name" db:"last_name"`
	Identifier      string `json:"identifier" db:"identifier"`
	ImageURL        string `json:"profile_image_url" db:"profile_image_url"`
	ProfileImageURL string `json:"image_url" db:"image_url"`
}

type Role struct {
	Role string `json:"role" db:"role"`
}

type RoleName struct {
	RoleName string `json:"role_name" db:"role_name"`
}

type MembershipRequest struct {
	UserID          string `json:"user_id" db:"user_id"`
	OrganizationID  string `json:"organization_id" db:"organization_id"`
	Identifier      string `json:"identifier" db:"identifier"`
	ImageURL        string `json:"image_url" db:"image_url"`
	ProfileImageURL string `json:"profile_image_url" db:"profile_image_url"`
	Role            string `json:"role" db:"role"`
	RoleName        string `json:"role_name" db:"role_name"`
}
