package http

import (
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
)

func (o *Endpoints) handleListConversations(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	conversations, err := o.conversationDBSvc.ListConversations(ctx)
	if err != nil {
		slog.ErrorContext(ctx, fmt.Sprintf("error listing conversations: %v", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = sendJsonResponse(w, conversations)
	if err != nil {
		slog.ErrorContext(ctx, fmt.Sprintf("error sending json response: %v", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (o *Endpoints) handleListMessages(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := r.PathValue("id")
	conversationId, err := strconv.Atoi(id)
	if err != nil {
		slog.Warn("bad conversation id", slog.String("conversation_id", id))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	messages, err := o.conversationDBSvc.GetMessagesByConversationID(ctx, conversationId)
	if err != nil {
		slog.ErrorContext(ctx, fmt.Sprintf("error getting messages for conversation %d: %v", conversationId, err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = sendJsonResponse(w, messages)
	if err != nil {
		slog.ErrorContext(ctx, fmt.Sprintf("error sending json response: %v", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
