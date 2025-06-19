package conversation

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/StevenWeathers/hatch-messaging-service/domain"
	"github.com/jackc/pgx/v5"
)

type Service struct {
	DB *pgx.Conn
}

func (s *Service) ListConversations(ctx context.Context) ([]domain.Conversation, error) {
	conversations := make([]domain.Conversation, 0)

	rows, err := s.DB.Query(ctx, `
		SELECT id FROM hatch.conversation;
	`)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		r := domain.Conversation{}
		if err := rows.Scan(&r.ID); err != nil {
			slog.Error(fmt.Sprintf("error scanning Conversation: %v", err))
		} else {
			conversations = append(conversations, r)
		}
	}

	return conversations, nil
}

func (s *Service) GetMessagesByConversationID(ctx context.Context, conversationId int) ([]domain.Message, error) {
	conversations := make([]domain.Message, 0)

	rows, err := s.DB.Query(ctx, `
		SELECT id, conversation_id, sender_id, type, body, timestamp FROM hatch.conversation_message
		WHERE conversation_id = $1;
	`, conversationId)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		r := domain.Message{}
		if err := rows.Scan(&r.ID, r.ConversationID, &r.SenderID, &r.Type, &r.Body, &r.Timestamp); err != nil {
			slog.Error(fmt.Sprintf("error scanning Conversation: %v", err))
		} else {
			conversations = append(conversations, r)
		}
	}

	return conversations, nil
}

func (s *Service) UpsertConversation(ctx context.Context, senderId int, recipientId int) (int, error) {
	var id int

	tx, err := s.DB.BeginTx(ctx, pgx.TxOptions{
		IsoLevel: pgx.ReadCommitted,
	})
	if err != nil {
		slog.Error(fmt.Sprintf("error starting transaction: %v", err))
		return 0, fmt.Errorf("error starting transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	// Check if a conversation already exists between the sender and recipient
	err = tx.QueryRow(ctx, `
		SELECT id
		FROM hatch.conversation
		WHERE id IN (
			SELECT conversation_id
			FROM hatch.conversation_contact
			WHERE contact_id IN ($1, $2)
			GROUP BY conversation_id
			HAVING COUNT(DISTINCT contact_id) = 2
		)
		LIMIT 1;
	`, senderId, recipientId).Scan(&id)
	if err != nil {
		if err == pgx.ErrNoRows {
			// No existing conversation, create a new one
			err = tx.QueryRow(ctx, `
				INSERT INTO hatch.conversation DEFAULT VALUES
				RETURNING id;
			`).Scan(&id)
			if err != nil {
				slog.Error(fmt.Sprintf("error creating new conversation: %v", err))
				return 0, fmt.Errorf("error creating new conversation: %w", err)
			}

			// Insert the sender and recipient into the conversation_contact table
			_, err = s.DB.Exec(ctx, `
				INSERT INTO hatch.conversation_contact (conversation_id, contact_id)
				VALUES ($1, $2), ($1, $3);
			`, id, senderId, recipientId)
			if err != nil {
				slog.Error(fmt.Sprintf("error inserting contacts into conversation: %v", err))
				return 0, fmt.Errorf("error inserting contacts into conversation: %w", err)
			}
		} else {
			slog.Error(fmt.Sprintf("error checking existing conversation: %v", err))
			return 0, fmt.Errorf("error checking existing conversation: %w", err)
		}
	}

	if commitErr := tx.Commit(ctx); commitErr != nil {
		slog.Error(fmt.Sprintf("error committing transaction: %v", commitErr))
		return 0, fmt.Errorf("error committing transaction: %w", commitErr)
	}

	return id, nil
}
