package scheduled

import (
	"context"

	"github.com/StevenWeathers/hatch-messaging-service/domain"
)

type Service struct {
	dbSvc ConversationDBSvc
}

type ConversationDBSvc interface {
	GetScheduledMessages(context.Context) ([]domain.Message, error)
}
