package entities

import (
    "time"
    "github.com/google/uuid"
)

type Conversation struct {
    ID        uuid.UUID `db:"id"`
    BookingID uuid.UUID `db:"booking_id"`
    CreatedAt time.Time `db:"created_at"`
}

type ConvParticipant struct{
    ConversationID  uuid.UUID `db:"conversation_id"`
	Role            string    `db:"role"`
	UserId          uuid.UUID `db:"user_id"`
}