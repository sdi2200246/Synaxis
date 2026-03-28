package entities

import (
    "time"
    "github.com/google/uuid"
)

type Message struct {
    ID             uuid.UUID `json:"id" db:"id"`
    ConversationID uuid.UUID `json:"conversation_id" db:"conversation_id"`
    SenderID       uuid.UUID `json:"sender_id" db:"sender_id"`
    Content        string    `json:"content" db:"content"`
    IsRead         bool      `json:"is_read" db:"is_read"`
    SentAt         time.Time `json:"sent_at" db:"sent_at"`
}
