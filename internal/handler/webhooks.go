package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/All-Things-Muchiri/server/internal/config"
	"github.com/All-Things-Muchiri/server/internal/domain"
	"github.com/All-Things-Muchiri/server/internal/service"
)

type WebhookHandler struct {
	userService *service.UserService
	whSecret string
}

func NewWebhookHandler(whSecret string, service *service.UserService) *WebhookHandler {
	return &WebhookHandler{
		userService: service,
		whSecret: whSecret,
	}
}

func (wh *WebhookHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("Failed to read body: %v", err)
		http.Error(w, "Failed to read body", http.StatusBadRequest)
		return
	}
	
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
		var usrData config.UserCreatedEvent
		
		if err := json.Unmarshal(whEvent.Data, &usrData); err != nil {
			log.Printf("Failed to decode event: %v", err)
			http.Error(w, "Failed to decode event", http.StatusInternalServerError)
			return
		}
		
		usrRequest := &domain.UserRequest{
			ID: usrData.UserID,
			Name: fmt.Sprintf("%s %s", usrData.Firstname, usrData.Lastname),
			Email: usrData.EmailAddresses[0].EmailAdresses,
			Image: usrData.ImageURL,
			EmailVerified: true,
		}
		
		processingError = wh.userService.CreateUser(ctx, usrRequest)
	}
	
	if processingError != nil {
		log.Printf("Failed to save user: %v", processingError)
		http.Error(w, "Failed to save user", http.StatusBadRequest)
		return
	}
	
	w.WriteHeader(http.StatusOK)
}