package domain

import "time"

type Contact struct {
	ID          int    `json:"id"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Email       string `json:"email"`        // e.g., "contact@usechatapp.com"
	PhoneNumber string `json:"phone_number"` // e.g., "+12016661234"
}

type Conversation struct {
	ID int `json:"id"`
}

type Message struct {
	ID             int        `json:"id"`
	ConversationID int        `json:"conversation_id"`
	SenderID       int        `json:"sender_id"`
	Type           string     `json:"type"` // e.g., "sms", "mms", "email"
	Body           string     `json:"body"`
	Attachments    *[]string  `json:"attachments,omitempty"` // e.g., ["https://example.com/image.jpg"]
	Timestamp      time.Time  `json:"timestamp"`             // e.g., "2024-11-01T14:00:00Z"
	ScheduledTime  *time.Time `json:"scheduled_time"`        // e.g., "2025-02-07T14:00:00Z"
}
