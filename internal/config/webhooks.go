package config

import (
	"database/sql"
	"encoding/json"
	"time"
)

type WebhookEvent struct {
	ID        string          `json:"instance_id"`
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

func NewWebhookEvent() *WebhookEvent {
	return &WebhookEvent{
		Timestamp: time.Now(),
	}
}
