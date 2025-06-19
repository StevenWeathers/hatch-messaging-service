package http

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
)

func (o *Endpoints) handle() {
	// k8s healthcheck /healthz as per convention
	o.router.HandleFunc("GET /healthz", o.handleHealthz)

	// API endpoints for sending messages
	o.router.HandleFunc("POST /api/messages/sms", o.handleSMSMessage)
	o.router.HandleFunc("POST /api/messages/email", o.handleEmailMessage)

	// Webhook endpoints for receiving messages from external services
	o.router.HandleFunc("POST /api/webhooks/sms", o.handleSMSWebhook)
	o.router.HandleFunc("POST /api/webhooks/email", o.handleEmailWebhook)

	// API endpoints for viewing saved conversations and messages
	o.router.HandleFunc("GET /api/conversations", o.handleListConversations)
	o.router.HandleFunc("GET /api/conversations/{id}/messages", o.handleListMessages)
}

func (o *Endpoints) ListenAndServe(ctx context.Context) error {
	slog.InfoContext(ctx, fmt.Sprintf("Listening on %s", o.config.ListenAddress))
	return http.ListenAndServe(o.config.ListenAddress, o.router)
}

func New(config Config, convDbService ConversationDBSvc) *Endpoints {
	router := http.NewServeMux()
	e := &Endpoints{
		config:            config,
		router:            router,
		conversationDBSvc: convDbService,
	}
	e.handle()
	return e
}
