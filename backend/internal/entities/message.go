package entities

import (
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

type MessageUpdate struct {
    Status *int
    Content *string
    IsRead *bool
}