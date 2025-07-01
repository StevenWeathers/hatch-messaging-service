package scheduled

import (
	"context"
	"log/slog"
	"time"
)

func New(dbSvc ConversationDBSvc) *Service {
	return &Service{
		dbSvc: dbSvc,
	}
}

func (s *Service) SendScheduledMessages() {
	ticker := time.NewTicker(1 * time.Second)
	quit := make(chan struct{})

	go func() {
		for {
			select {
			case <-ticker.C:
				messages, err := s.dbSvc.GetScheduledMessages(context.Background())
				if err != nil {
					slog.Error("Error getting scheduled messages")
				}

				for _, m := range messages {
					// send the message
					_ = m
				}
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()
}
