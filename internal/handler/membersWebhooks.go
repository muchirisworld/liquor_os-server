package handler

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"

	"github.com/All-Things-Muchiri/server/internal/config"
	"github.com/All-Things-Muchiri/server/internal/domain"
	"github.com/All-Things-Muchiri/server/internal/service"
)

type OrganizationMembershipCreatedOrUpdatedEvent struct {
	MembershipID string              `json:"id"`
	Organization domain.Organization `json:"organization"`
	Permissions  json.RawMessage     `json:"permissions"`
	User         domain.Member       `json:"public_user_data"`
	Role         domain.Role         `json:"role"`
	RoleName     domain.RoleName     `json:"role_name"`
}

type MembershipWebhookHandler struct {
	membershipService *service.MembershipService
	whSecret          string
}

func NewMembershipWebhookHandler(whSecret string, membershipService *service.MembershipService) *MembershipWebhookHandler {
	return &MembershipWebhookHandler{
		membershipService: membershipService,
		whSecret:          whSecret,
	}
}

func (wh *MembershipWebhookHandler) HandleMembershipWebhook(w http.ResponseWriter, r *http.Request) {
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

	whEvent := config.NewWebhookEvent()
	if err := json.Unmarshal(body, whEvent); err != nil {
		log.Printf("Failed to decode event: %v", err)
		http.Error(w, "Failed to decode event", http.StatusInternalServerError)
		return
	}

	var processingError error
	switch whEvent.Type {
	case "organizationMembership.created":
		processingError = wh.processOrganizationMembershipCreated(ctx, whEvent.Data)
	default:
		log.Printf("Unhandled event type: %s", whEvent.Type)
		w.WriteHeader(http.StatusOK)
		return
	}

	if processingError != nil {
		log.Printf("Failed to process membership event: %v", processingError)
		http.Error(w, "Failed to process membership event", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (wh *MembershipWebhookHandler) processOrganizationMembershipCreated(ctx context.Context, data json.RawMessage) error {
	membershipReq := &domain.MembershipRequest{}
	if err := json.Unmarshal(data, membershipReq); err != nil {
		log.Printf("Failed to read membership request data: %v", err)
		return err
	}
	return wh.membershipService.CreateMember(ctx, membershipReq)
}
