package entities

import (
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
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
        return errors.New("only the sender can edit a message")
    }
    if m.Status != 0 {
        return errors.New("cannot edit a deleted message")
    }
    return nil
}

func (m Message) ValidateContent(content string) error {
    if strings.TrimSpace(content) == "" {
        return errors.New("message content cannot be empty")
    }
    return nil
}

func (m Message) CanTransitionTo(target int) error {
	if target == m.Status {
		return errors.New("already in this state")
	}
	if target < m.Status {
		return errors.New("cannot reverse a deletion")
	}
	if target < 1 || target > 2 {
		return errors.New("invalid status")
	}
	return nil
}


type MessageUpdate struct {
    Status *int
    Content *string
}