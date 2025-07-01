package conversation

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/StevenWeathers/hatch-messaging-service/domain"
	"github.com/jackc/pgx/v5"
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

// GetMessageCountBySenderID gets a count of messages by sender ID
func (s *Service) GetMessageCountBySenderID(ctx context.Context, senderID int) (int, error) {
	count := 0

	err := s.DB.QueryRow(ctx, `
		SELECT count(timestamp) FROM hatch.conversation_message WHERE sender_id = $1 AND timestamp < current_timestamp - interval '60 seconds';
	`, senderID).Scan(&count)
	if err != nil {
		if err == pgx.ErrNoRows {
			return 0, nil
		}
		return count, err
	}

	return count, nil
}

func (s *Service) GetScheduledMessages(ctx context.Context) ([]domain.Message, error) {
	messages := make([]domain.Message, 0)

	rows, err := s.DB.Query(ctx, `
		SELECT id, conversation_id, sender_id, type, body, attachments, timestamp FROM hatch.conversation_message WHERE scheduled_time < current_timestamp AND message_sent IS NOT TRUE;
	`)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		r := domain.Message{}
		if err := rows.Scan(&r.ID, &r.ConversationID, &r.SenderID, &r.Type, &r.Body, &r.Attachments, &r.Timestamp); err != nil {
			slog.Error(fmt.Sprintf("error scanning Message: %v", err))
		} else {
			messages = append(messages, r)
		}
	}

	return messages, nil
}
