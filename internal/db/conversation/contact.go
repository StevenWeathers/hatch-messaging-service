package conversation

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/StevenWeathers/hatch-messaging-service/domain"
	"github.com/jackc/pgx/v5"
)

func (s *Service) UpsertContact(ctx context.Context, contact domain.Contact) (int, error) {
	var id int

	tx, err := s.DB.BeginTx(ctx, pgx.TxOptions{
		IsoLevel: pgx.ReadCommitted,
	})
	if err != nil {
		slog.Error(fmt.Sprintf("error starting transaction: %v", err))
		return 0, fmt.Errorf("error starting transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	// this isn't realistic, in a real application you would ideally have better a defined contact model that includes the users multiple contact methods
	// for now, we will just use phone number or email to identify the contact
	query := "SELECT id FROM hatch.contact WHERE email = $1"
	var queryValue string = contact.Email
	if contact.Email == "" {
		query = "SELECT id FROM hatch.contact WHERE phone_number = $1"
		queryValue = contact.PhoneNumber
	}
	err = tx.QueryRow(ctx, query, queryValue).Scan(&id)
	if err != nil {
		if err == pgx.ErrNoRows {
			// Contact does not exist, insert new contact
			err = tx.QueryRow(ctx, `
				INSERT INTO hatch.contact (first_name, last_name, phone_number, email)
				VALUES ($1, $2, NULLIF($3, ''), NULLIF($4, ''))
				RETURNING id
			`, contact.FirstName, contact.LastName, contact.PhoneNumber, contact.Email).Scan(&id)
			if err != nil {
				slog.ErrorContext(ctx, fmt.Sprintf("error inserting new contact: %v", err))
				return 0, fmt.Errorf("error inserting new contact: %w", err)
			}
		} else {
			slog.ErrorContext(ctx, fmt.Sprintf("error querying contact: %v", err))
			return 0, fmt.Errorf("error querying contact: %w", err)
		}
	}

	if commitErr := tx.Commit(ctx); commitErr != nil {
		slog.Error(fmt.Sprintf("error committing transaction: %v", commitErr))
		return 0, fmt.Errorf("error committing transaction: %w", commitErr)
	}

	return id, nil
}
