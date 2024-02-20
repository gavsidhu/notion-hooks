package models

import (
	"time"

	"github.com/gavsidhu/notion-hooks/internal/notion"
)

type Webhook struct {
	ID               string     `json:"id"`
	Name             string     `json:"name"`
	Description      string     `json:"description"`
	UserID           string     `json:"user_id"`
	URL              string     `json:"url"`
	Secret           string     `json:"secret"`
	Events           []string   `json:"events"`
	IsActive         bool       `json:"is_active"`
	Status           string     `json:"status"`
	PollingInterval  int        `json:"polling_interval"`
	LastPolled       *time.Time `json:"last_polled"`
	NotionObjectID   string     `json:"notion_object_id"`
	NotionObjectType string     `json:"notion_object_type"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
}

type WebhookResponse struct {
	Webhook Webhook `json:"webhook"`
}

type WebhooksResponse struct {
	Webhooks []Webhook `json:"webhooks"`
}

type WebhookCreateRequest struct {
	URL             string `json:"url"`
	ContentType     string `json:"content_type"`
	Secret          string `json:"secret"`
	Events          string `json:"events"`
	IsActive        bool   `json:"is_active"`
	PollingInterval int    `json:"polling_interval"`
	NotionDataID    string `json:"notion_data_id"`
}

type WebhookUpdateRequest struct {
	URL             string `json:"url"`
	ContentType     string `json:"content_type"`
	Secret          string `json:"secret"`
	Events          string `json:"events"`
	IsActive        bool   `json:"is_active"`
	PollingInterval int    `json:"polling_interval"`
	NotionDataID    string `json:"notion_data_id"`
}

type NotionDatabasePageIDRow struct {
	ID        int       `json:"id"`
	WebhookID string    `json:"webhook_id"`
	UserID    string    `json:"user_id"`
	PageIDs   []string  `json:"page_ids"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type NotionDatabaseDetailRow struct {
	ID                  int                          `json:"id"`
	WebhooksID          string                       `json:"webhooks_id"`
	UserID              string                       `json:"user_id"`
	DatabasePageDetails notion.DatabaseQueryResponse `json:"database_page_details"`
	CreatedAt           time.Time                    `json:"created_at"`
	UpdatedAt           time.Time                    `json:"updated_at"`
}

type WebhookLog struct {
	ID           string `json:"id"`
	WebhookID    string `json:"webhook_id"`
	Status       string `json:"status"`
	ResponseCode int    `json:"response_code"`
	Payload      string `json:"payload,omitempty"`
	ErrorMessage string `json:"error_message,omitempty"`
	AttemptedAt  string `json:"attempted_at"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
}

type EventData struct {
	ObjectID   string `json:"object_id"`
	ObjectType string `json:"object_type"`
	CreatedAt  int64  `json:"created_at"`
}

type EventsToSend struct {
	Type      string    `json:"type"`
	UserID    string    `json:"user_id"`
	WebhookID string    `json:"webhook_id"`
	Data      EventData `json:"data"`
}

type Event struct {
	ID        string    `json:"id"`
	WebhookID string    `json:"webhook_id"`
	Type      string    `json:"type"`
	Data      EventData `json:"data"`
	CreatedAt int64     `json:"created_at"`
}

type InitialPollMessage struct {
	WebhookID        string `json:"webhook_id"`
	UserID           string `json:"user_id"`
	NotionObjectID   string `json:"notion_object_id"`
	NotionObjectType string `json:"notion_object_type"`
}
