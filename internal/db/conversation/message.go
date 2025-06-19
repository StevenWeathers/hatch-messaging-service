package conversation

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/StevenWeathers/hatch-messaging-service/domain"
)

func (s *Service) SaveMessage(ctx context.Context, message domain.Message) error {
	_, err := s.DB.Exec(ctx, `
		INSERT INTO hatch.conversation_message (conversation_id, sender_id, type, body, timestamp)
		VALUES ($1, $2, $3, $4, $5);
	`, message.ConversationID, message.SenderID, message.Type, message.Body, message.Timestamp)
	if err != nil {
		slog.Error(fmt.Sprintf("error saving message: %v", err))
	}
	return err
}
