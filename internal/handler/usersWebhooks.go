package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/All-Things-Muchiri/server/internal/config"
	"github.com/All-Things-Muchiri/server/internal/domain"
	"github.com/All-Things-Muchiri/server/internal/service"
)

// Limit request body size to 1MB
const MAX_BYTES = 1 << 20 // 1MB

type UserCreatedOrUpdatedEvent struct {
	Firstname      string         `json:"first_name"`
	Lastname       string         `json:"last_name"`
	UserID         string         `json:"id"`
	EmailAddresses []EmailAddress `json:"email_addresses"`
	ImageURL       string         `json:"image_url"`
}

type EmailAddress struct {
	EmailAdresses string `json:"email_address"`
}

type UsersWebhookHandler struct {
	userService *service.UserService
	whSecret    string
}

func NewUsersWebhookHandler(whSecret string, service *service.UserService) *UsersWebhookHandler {
	return &UsersWebhookHandler{
		userService: service,
		whSecret:    whSecret,
	}
}

func (wh *UsersWebhookHandler) HandleUsersWebhooks(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	r.Body = http.MaxBytesReader(w, r.Body, MAX_BYTES)

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("Failed to read body: %v", err)
		var maxBytesErr *http.MaxBytesError
		if errors.As(err, &maxBytesErr) {
			http.Error(w, "Request body too large", http.StatusRequestEntityTooLarge)
		} else {
			http.Error(w, "Failed to read body", http.StatusBadRequest)
		}
		return
	}

	log.Printf("Event payload: %v", string(body))

	// TODO: Verify webhook
	whEvent := config.NewWebhookEvent()
	if err := json.Unmarshal(body, whEvent); err != nil {
		log.Printf("Failed to decode event: %v", err)
		http.Error(w, "Failed to decode event", http.StatusInternalServerError)
		return
	}

	var processingError error
	switch whEvent.Type {
	case "user.created":
		processingError = wh.handleUserCreated(ctx, whEvent.Data)
	case "user.updated":
		processingError = wh.handleUserUpdated(ctx, whEvent.Data)
	default:
		log.Printf("Unhandled event type: %s", whEvent.Type)
		w.WriteHeader(http.StatusOK)
		return
	}

	if processingError != nil {
		log.Printf("Failed to process user event: %v", processingError)
		http.Error(w, "Failed to process user event", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (wh *UsersWebhookHandler) handleUserCreated(ctx context.Context, data json.RawMessage) error {
	var usrData UserCreatedOrUpdatedEvent

	if err := json.Unmarshal(data, &usrData); err != nil {
		log.Printf("Failed to decode user.created event: %v", err)
		return err
	}

	email := ""
	if len(usrData.EmailAddresses) > 0 {
		email = usrData.EmailAddresses[0].EmailAdresses
	}

	usrRequest := &domain.UserRequest{
		ID:            usrData.UserID,
		Name:          fmt.Sprintf("%s %s", usrData.Firstname, usrData.Lastname),
		Email:         email,
		Image:         usrData.ImageURL,
		EmailVerified: true,
	}

	return wh.userService.CreateUser(ctx, usrRequest)
}

func (wh *UsersWebhookHandler) handleUserUpdated(ctx context.Context, data json.RawMessage) error {
	var usrData UserCreatedOrUpdatedEvent

	if err := json.Unmarshal(data, &usrData); err != nil {
		log.Printf("Failed to decode user.updated event: %v", err)
		return err
	}

	email := ""
	if len(usrData.EmailAddresses) > 0 {
		email = usrData.EmailAddresses[0].EmailAdresses
	}

	usrRequest := &domain.UserRequest{
		ID:            usrData.UserID,
		Name:          fmt.Sprintf("%s %s", usrData.Firstname, usrData.Lastname),
		Email:         email,
		Image:         usrData.ImageURL,
		EmailVerified: true,
	}

	return wh.userService.UpdateUser(ctx, usrRequest)
}
