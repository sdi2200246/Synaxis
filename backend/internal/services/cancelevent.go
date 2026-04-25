package services

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/sdi2200246/synaxis/internal/entities"
	apperr "github.com/sdi2200246/synaxis/internal/error"
	"github.com/sdi2200246/synaxis/internal/interfaces"
)

type CancelEventService struct {
	eventRepo   interfaces.EventRepository
	bookingRepo interfaces.BookingRepository
	messageRepo interfaces.MessagesRepository
	eventBus    interfaces.EventBus
}

func NewCancelEventService(er interfaces.EventRepository,br interfaces.BookingRepository,mr interfaces.MessagesRepository,
	eb interfaces.EventBus,
) *CancelEventService {
	return &CancelEventService{
		eventRepo:   er,
		bookingRepo: br,
		messageRepo: mr,
		eventBus:    eb,
	}
}

func (s *CancelEventService) Subscribe() {
	ch := s.eventBus.Subscribe("EventCancelled")
	go func() {
		for event := range ch {
			cancelled, ok := event.(EventCancelled)
			if !ok {
				slog.Error("CancelEventService: unexpected event type")
				continue
			}
			if err := s.handleCancellation(context.Background(), cancelled.EventID); err != nil {
				slog.Error("CancelEventService: handleCancellation failed",
					"error", err,
					"event_id", cancelled.EventID,
				)
			}
		}
	}()
}

func (s *CancelEventService) handleCancellation(ctx context.Context, eventID uuid.UUID) error {
	event, err := s.eventRepo.GetByID(ctx, eventID)
	if err != nil {
		slog.Error("CancelEventService: failed to get event",
			"error", err,
			"event_id", eventID,
		)
		return err
	}

	bookings, err := s.bookingRepo.GetByEventID(ctx, eventID)
	if err != nil {
		slog.Error("CancelEventService: failed to get bookings",
			"error", err,
			"event_id", eventID,
		)
		return err
	}

	if len(bookings) == 0 {
		return nil
	}

	for _, booking := range bookings {
		msg := entities.Message{
			ID:       uuid.New(),
			SenderID: event.OrganizerID,
			Content: fmt.Sprintf(
				"We're sorry to inform you that \"%s\" has been cancelled. If you have any questions please reach out to us.",
				event.Title,
			),
			IsRead: false,
			Status: 0,
			SentAt: time.Now(),
		}

		existing, err := s.messageRepo.GetByBookingID(ctx, booking.ID)
		if err != nil && !errors.Is(err, apperr.ErrNotFound) {
			slog.Error("CancelEventService: failed to check existing conversation",
				"error", err,
				"booking_id", booking.ID,
			)
			continue
		}

		if errors.Is(err, apperr.ErrNotFound) {
			conv := entities.Conversation{
				ID:        uuid.New(),
				BookingID: booking.ID,
				CreatedAt: time.Now(),
			}
			if err := s.messageRepo.CreateConversationWithMessage(
				ctx,
				conv,
				event.OrganizerID,
				booking.UserID,
				msg,
			); err != nil {
				slog.Error("CancelEventService: failed to notify attendee",
					"error", err,
					"booking_id", booking.ID,
					"user_id", booking.UserID,
				)
				continue
			}
		} else {
			msg.ConversationID = existing.ID
			if err := s.messageRepo.Create(ctx, msg); err != nil {
				slog.Error("CancelEventService: failed to send cancellation message",
					"error", err,
					"conversation_id", existing.ID,
					"user_id", booking.UserID,
				)
				continue
			}
		}

		slog.Info("CancelEventService: attendee notified",
			"booking_id", booking.ID,
			"user_id", booking.UserID,
			"event_id", eventID,
		)
	}

	return nil
}