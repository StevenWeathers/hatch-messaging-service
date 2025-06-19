package http

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/StevenWeathers/hatch-messaging-service/domain"
)

func (o *Endpoints) handleSMSWebhook(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	body, err := io.ReadAll(r.Body)
	if err != nil {
		slog.ErrorContext(ctx, fmt.Sprintf("error reading body: %v", err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var u = SMSWebhookRequest{}
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

	err = o.conversationDBSvc.SaveMessage(ctx, domain.Message{
		ConversationID: convoId,
		SenderID:       senderID,
		Type:           u.Type,
		Body:           u.Body,
		Timestamp:      u.Timestamp,
	})
	if err != nil {
		slog.ErrorContext(ctx, fmt.Sprintf("error saving message: %v", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Here I would process the SMS webhook request and save to the database, if an error occurs, log it and return an error status code
	// Depending on the upstream service, I might also implement a queue to handle retries if the upsteram service doesn't handle retries itself.

	err = sendJsonResponse(w, `{"alive": true}`)
	if err != nil {
		slog.ErrorContext(ctx, fmt.Sprintf("error sending json response: %v", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (o *Endpoints) handleEmailWebhook(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	body, err := io.ReadAll(r.Body)
	if err != nil {
		slog.ErrorContext(ctx, fmt.Sprintf("error reading body: %v", err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var u = EmailWebhookRequest{}
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

	// Here I would process the SMS webhook request and save to the database, if an error occurs, log it and return an error status code
	// Depending on the upstream service, I might also implement a queue to handle retries if the upsteram service doesn't handle retries itself.

	err = sendJsonResponse(w, `{"alive": true}`)
	if err != nil {
		slog.ErrorContext(ctx, fmt.Sprintf("error sending json response: %v", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
