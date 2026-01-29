package config

import (
	"database/sql"
	"encoding/json"
	"time"
)

type WebhookEvent struct {
	ID        string          `json:"id"`
	Type      string          `json:"type"`
	Timestamp time.Time       `json:"-"`
	Data      json.RawMessage `json:"data"`
}

type WebhookLog struct {
	ID              int64
	EventID         string
	EventType       string
	Payload         string
	ProcessedAt     time.Time
	ProcessingError sql.NullString
	Success         bool
}

type UserCreatedEvent struct {
	Firstname      string         `json:"first_name"`
	Lastname       string         `json:"last_name"`
	UserID         string         `json:"id"`
	EmailAddresses []EmailAddress `json:"email_addresses"`
	ImageURL       string         `json:"image_url"`
}


type EmailAddress struct {
	EmailAdresses string `json:"email_address"`
}

func NewWebhookEvent() *WebhookEvent {
	return &WebhookEvent{
		Timestamp: time.Now(),
	}
}
