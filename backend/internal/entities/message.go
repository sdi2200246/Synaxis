package entities

import (
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	apperr "github.com/sdi2200246/synaxis/internal/error"
)

type Message struct {
    ID             uuid.UUID `db:"id"`
    ConversationID uuid.UUID `db:"conversation_id"`
    SenderID       uuid.UUID `db:"sender_id"`
    Content        string    `db:"content"`
    Status         int       `db:"status"`
    IsRead         bool      `db:"is_read"`
    SentAt         time.Time `db:"sent_at"`
    UpdatedAt      *time.Time `db:"updated_at"`
}

func (m Message) CanEditContent(callerID uuid.UUID) error {
    if m.SenderID != callerID {
        return fmt.Errorf("only sender can modify this message: %w", apperr.ErrForbidden)
    }
    if m.Status != 0 {
        return fmt.Errorf("cannot edit a deleted message: %w" , apperr.ErrConflict)
    }
    return nil
}

func (m Message) ValidateContent(content string) error {
    if strings.TrimSpace(content) == "" {
        return fmt.Errorf("message content cannot be empty: %w", apperr.ErrBadInput)
    }
    return nil
}

func (m Message) CanTransitionTo(target int) error {
    if target < 1 || target > 2 {
        return fmt.Errorf("invalid deletion status %d, must be 1 or 2: %w", target, apperr.ErrBadInput)
    }
    if target == m.Status {
        return fmt.Errorf("message is already in status %d: %w", target, apperr.ErrConflict)
    }
    if target < m.Status {
        return fmt.Errorf("cannot reverse a deletion from status %d to %d: %w", m.Status, target, apperr.ErrConflict)
    }
    return nil
}


type MessageUpdate struct {
    Status *int
    Content *string
}