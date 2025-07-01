package http

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/StevenWeathers/hatch-messaging-service/domain"
)

func (o *Endpoints) handleSMSMessage(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	body, err := io.ReadAll(r.Body)
	if err != nil {
		slog.ErrorContext(ctx, fmt.Sprintf("error reading body: %v", err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var u = SMSMessageRequest{}
	jsonErr := json.Unmarshal(body, &u)
	if jsonErr != nil {
		slog.ErrorContext(ctx, fmt.Sprintf("error unmarshalling body: %v", jsonErr))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	senderID, err := o.conversationDBSvc.UpsertContact(ctx, domain.Contact{
		FirstName:   "TestFirstname",
		LastName:    "TestLastname",
		PhoneNumber: u.From,
	})
	if err != nil {
		slog.ErrorContext(ctx, fmt.Sprintf("error upserting sender contact: %v", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// rate limit 5 per 60 seconds
	// call the db to get a count of messages the senderID has sent in last 60 seconds
	// if 5 >= then would respond with an error message and status code
	messageCount, err := o.conversationDBSvc.GetMessageCountBySenderID(ctx, senderID)
	if err != nil {
		slog.ErrorContext(ctx, fmt.Sprintf("error upserting sender contact: %v", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if messageCount >= 5 {
		slog.ErrorContext(ctx, "max message count of 5 per 60 seconds reached")
		w.WriteHeader(http.StatusTooManyRequests)
		return
	}

	receiverID, err := o.conversationDBSvc.UpsertContact(ctx, domain.Contact{
		FirstName:   "TestFirstname",
		LastName:    "TestLastname",
		PhoneNumber: u.To,
	})
	if err != nil {
		slog.ErrorContext(ctx, fmt.Sprintf("error upserting receiver contact: %v", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	convoId, err := o.conversationDBSvc.UpsertConversation(ctx, senderID, receiverID)
	if err != nil {
		slog.ErrorContext(ctx, fmt.Sprintf("error upserting conversation: %v", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	message := domain.Message{
		ConversationID: convoId,
		SenderID:       senderID,
		Type:           u.Type,
		Body:           u.Body,
		Timestamp:      u.Timestamp,
		ScheduledTime:  u.ScheduledTime,
	}

	err = o.conversationDBSvc.SaveMessage(ctx, message)
	if err != nil {
		slog.ErrorContext(ctx, fmt.Sprintf("error saving message: %v", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Here I would send to appropriate service (e.g., Twilio, etc.)
	// If the service returns an error I would log it and return an error response.
	// Additionally, depending on the service design I might store the send status in the database and handle retries.
	// or use a queuing system like Kafka to handle incoming messages.

	if message.ScheduledTime != nil {
		// do nothing else with message and handle it in a cron like job
		err = sendJsonResponse(w, `{"success": true}`)
		if err != nil {
			slog.ErrorContext(ctx, fmt.Sprintf("error sending json response: %v", err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	// update message with SENT column boolean to TRUE

	err = sendJsonResponse(w, `{"success": true}`)
	if err != nil {
		slog.ErrorContext(ctx, fmt.Sprintf("error sending json response: %v", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (o *Endpoints) handleEmailMessage(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	body, err := io.ReadAll(r.Body)
	if err != nil {
		slog.ErrorContext(ctx, fmt.Sprintf("error reading body: %v", err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var u = EmailMessageRequest{}
	jsonErr := json.Unmarshal(body, &u)
	if jsonErr != nil {
		slog.ErrorContext(ctx, fmt.Sprintf("error unmarshalling body: %v", jsonErr))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	senderID, err := o.conversationDBSvc.UpsertContact(ctx, domain.Contact{
		FirstName: "TestFirstname",
		LastName:  "TestLastname",
		Email:     u.From,
	})
	if err != nil {
		slog.ErrorContext(ctx, fmt.Sprintf("error upserting sender contact: %v", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	receiverID, err := o.conversationDBSvc.UpsertContact(ctx, domain.Contact{
		FirstName: "TestFirstname",
		LastName:  "TestLastname",
		Email:     u.To,
	})
	if err != nil {
		slog.ErrorContext(ctx, fmt.Sprintf("error upserting receiver contact: %v", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	convoId, err := o.conversationDBSvc.UpsertConversation(ctx, senderID, receiverID)
	if err != nil {
		slog.ErrorContext(ctx, fmt.Sprintf("error upserting conversation: %v", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = o.conversationDBSvc.SaveMessage(ctx, domain.Message{
		ConversationID: convoId,
		SenderID:       senderID,
		Type:           "email",
		Body:           u.Body,
		Timestamp:      u.Timestamp,
	})
	if err != nil {
		slog.ErrorContext(ctx, fmt.Sprintf("error saving message: %v", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Here I would send to appropriate service (e.g., SendGrid, etc.)
	// If the service returns an error I would log it and return an error response.
	// Additionally, depending on the service design I might store the send status in the database and handle retries.
	// or use a queuing system like Kafka to handle incoming messages.

	err = sendJsonResponse(w, `{"success": true}`)
	if err != nil {
		slog.ErrorContext(ctx, fmt.Sprintf("error sending json response: %v", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
