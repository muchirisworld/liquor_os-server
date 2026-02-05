package handler

import (
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

type OrganizationCreatedOrUpdatedEvent struct {
	OrganizationID string          `json:"id"`
	Name           string          `json:"name"`
	Slug           string          `json:"slug"`
	CreatedBy      string          `json:"created_by"`
	LogoURL        string          `json:"logo_url"`
	ImageURL       string          `json:"image_url"`
	Metadata       json.RawMessage `json:"-"`
}

type OrgsWebhookHandler struct {
	orgService *service.OrganizationsService
	whSecret   string
}

func NewOrganizationsWebhookHandler(whSecret string, service *service.OrganizationsService) *OrgsWebhookHandler {
	return &OrgsWebhookHandler{
		orgService: service,
		whSecret:   whSecret,
	}
}

func (wh *OrgsWebhookHandler) HandleOrganizationsWebhook(w http.ResponseWriter, r *http.Request) {
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
	case "organization.created":
		orgReq := &domain.OrganizationRequest{}
		if err := json.Unmarshal(whEvent.Data, orgReq); err != nil {
			log.Printf("Failed to read organization request data: %v", err)
			http.Error(w, fmt.Sprintf("Failed to read organization request data: %v", err), http.StatusInternalServerError)
			return
		}
		processingError = wh.orgService.CreateOrganization(ctx, orgReq)

	default:
		log.Printf("Unhandled event type: %s", whEvent.Type)
		w.WriteHeader(http.StatusOK)
		return
	}

	if processingError != nil {
		log.Printf("Failed to process organization event: %v", processingError)
		http.Error(w, "Failed to process organization event", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}
