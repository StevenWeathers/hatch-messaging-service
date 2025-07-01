package http

import (
	"context"
	"net/http"
	"time"

	"github.com/StevenWeathers/hatch-messaging-service/domain"
)

type Config struct {
	ListenAddress string
}

// Endpoints represents the http service and its endpoints
type Endpoints struct {
	config            Config
	router            *http.ServeMux
	conversationDBSvc ConversationDBSvc
}

type SMSMessageRequest struct {
	From          string     `json:"from"`                  // e.g., "+12016661234"
	To            string     `json:"to"`                    // e.g., "+12016661234"
	Type          string     `json:"type"`                  // e.g., "sms", "mms"
	Body          string     `json:"body"`                  // e.g., "Hello! This is a test SMS message."
	Attachments   *[]string  `json:"attachments,omitempty"` // e.g., "["https://example.com/image.jpg"]"
	Timestamp     time.Time  `json:"timestamp"`             // e.g., "2024-11-01T14:00:00Z"
	ScheduledTime *time.Time `json:"scheduled_time"`        // e.g., "2025-02-07T14:00:00Z"
}
type EmailMessageRequest struct {
	From        string    `json:"from"`                  // e.g., "user@usehatchapp.com"
	To          string    `json:"to"`                    // e.g., "contact@gmail.com"
	Body        string    `json:"body"`                  // e.g., "Hello! This is a test email message with <b>HTML</b> formatting."
	Attachments *[]string `json:"attachments,omitempty"` // e.g., "["https://example.com/document.pdf"]"
	Timestamp   time.Time `json:"timestamp"`             // e.g., "2024-11-01T14:00:00Z"
}

type SMSWebhookRequest struct {
	MessagingProviderID string    `json:"messaging_provider_id"` // e.g., "1, 2"
	From                string    `json:"from"`                  // e.g., "+12016661234"
	To                  string    `json:"to"`                    // e.g., "+12016661234"
	Type                string    `json:"type"`                  // e.g., "sms", "mms"
	Body                string    `json:"body"`                  // e.g., "Hello! This is a test SMS message."
	Attachments         *[]string `json:"attachments,omitempty"` // e.g., "["https://example.com/image.jpg"]"
	Timestamp           time.Time `json:"timestamp"`             // e.g., "2024-11-01T14:00:00Z"
}
type EmailWebhookRequest struct {
	MessagingProviderID string    `json:"messaging_provider_id"` // e.g., "1, 2"
	From                string    `json:"from"`                  // e.g., "user@usehatchapp.com"
	To                  string    `json:"to"`                    // e.g., "contact@gmail.com"
	Body                string    `json:"body"`                  // e.g., "Hello! This is a test email message with <b>HTML</b> formatting."
	Attachments         *[]string `json:"attachments,omitempty"` // e.g., "["https://example.com/document.pdf"]"
	Timestamp           time.Time `json:"timestamp"`             // e.g., "2024-11-01T14:00:00Z"
}

type ConversationDBSvc interface {
	ListConversations(context.Context) ([]domain.Conversation, error)
	GetMessagesByConversationID(context.Context, int) ([]domain.Message, error)
	SaveMessage(ctx context.Context, message domain.Message) error
	UpsertConversation(ctx context.Context, senderId int, recipientId int) (int, error)
	UpsertContact(ctx context.Context, contact domain.Contact) (int, error)
	GetMessageCountBySenderID(context.Context, int) (int, error)
}
