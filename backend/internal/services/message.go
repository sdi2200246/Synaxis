package services

import (
	"context"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/sdi2200246/synaxis/internal/entities"
	apperr "github.com/sdi2200246/synaxis/internal/error"
	"github.com/sdi2200246/synaxis/internal/interfaces"
)

type CreateConversationInput struct {
	BookingID    uuid.UUID
	OrganizerID  uuid.UUID
	AttendeeID 	 uuid.UUID
}

type CreateMessageInput struct {
	ConversationID uuid.UUID
	SenderID       uuid.UUID
	Content        string
}

type Message struct {
	ID             uuid.UUID
	ConversationID uuid.UUID
	SenderID       uuid.UUID
	Content        string
	IsRead         bool
	Deleted        bool
	SentAt         time.Time
	UpdatedAt      *time.Time
}

type Conversation struct {
	ID        uuid.UUID
	BookingID uuid.UUID
	CreatedAt time.Time
	UnseenCount int
}
type ConvParticipant struct{
	Role string
	UserID uuid.UUID
}


type ConversationWithParticipants struct {
	Conversation Conversation
	Participants []ConvParticipant
}

type MessageService struct {
	messagesRepo interfaces.MessagesRepository
	bookingRepo  interfaces.BookingRepository
	eventRepo    interfaces.EventRepository
}


func NewMessageService(r interfaces.MessagesRepository , br interfaces.BookingRepository  , er interfaces.EventRepository) *MessageService {
	return &MessageService{messagesRepo: r , bookingRepo: br , eventRepo: er}
}

func (s *MessageService) CreateConversation(ctx context.Context, input CreateConversationInput) (uuid.UUID, error) {

	booking, err := s.bookingRepo.GetByID(ctx, input.BookingID)
	if err != nil {
		return uuid.Nil, err
	}

	event, err := s.eventRepo.GetByTicketTypeID(ctx, booking.TicketTypeID)
	if err != nil {
		return uuid.Nil, err
	}

	if booking.UserID != input.AttendeeID || event.OrganizerID != input.OrganizerID {
		return uuid.Nil, apperr.ErrForbidden
	}

	conv := entities.Conversation{
		ID:        uuid.New(),
		BookingID: booking.ID,
		CreatedAt: time.Now(),
	}


	err = s.messagesRepo.CreateConversation(ctx, conv, input.OrganizerID , input.AttendeeID)
	if err != nil {
		return uuid.Nil, err
	}

	return conv.ID, nil
}

func (s *MessageService) SendMessage(ctx context.Context, input CreateMessageInput) error {

	msg := entities.Message{
		ID:             uuid.New(),
		ConversationID: input.ConversationID,
		SenderID:       input.SenderID,
		Content:        strings.TrimSpace(input.Content),
		IsRead:         false,
		Status:         0,
		SentAt:         time.Now(),
	}

	return s.messagesRepo.Create(ctx, msg)
}

func (s *MessageService) EditMessage(ctx context.Context, id uuid.UUID, mu entities.MessageUpdate) error {
	return s.messagesRepo.UpdateMessage(ctx, id, mu)
}

func (s *MessageService) GetConversationMessages(ctx context.Context,conversationID uuid.UUID,userID uuid.UUID)([]Message, error) {

	participants, err := s.messagesRepo.GetParticipantsByConversationID(ctx, conversationID)
	if err != nil {
		return nil, err
	}

	isParticipant := false
	for _, p := range participants {
		if p.UserId == userID {
			isParticipant = true
			break
		}
	}

	if !isParticipant {
		return nil, apperr.ErrForbidden
	}

	messages, err := s.messagesRepo.GetByConversationID(ctx, conversationID)
	if err != nil {
		return nil, err
	}

	result := make([]Message, len(messages))
	for i, m := range messages {
		result[i] = toMessage(m, resolveMessageStatus(m, userID))
	}

	return result, nil
}

func (s *MessageService) GetConversation(ctx context.Context,conversationID uuid.UUID,userID uuid.UUID) (ConversationWithParticipants, error) {

	conv, err := s.messagesRepo.GetConversationByID(ctx, conversationID)
	if err != nil {
		return ConversationWithParticipants{}, err
	}

	participants, err := s.messagesRepo.GetParticipantsByConversationID(ctx, conversationID)
	if err != nil {
		return ConversationWithParticipants{}, err
	}

	if len(participants) != 2 {
		return ConversationWithParticipants{}, apperr.ErrInternal
	}

	unseen , err := s.messagesRepo.GetUnreadMessagesCountByUser(ctx ,userID)
	if err != nil {
		return ConversationWithParticipants{}, err
	}


	return ConversationWithParticipants{
		Conversation: toConversation(conv , unseen[conversationID]),
		Participants: []ConvParticipant{
			toConvParticipant(participants[0]),
			toConvParticipant(participants[1]),
		},
	}, nil
}

func (s *MessageService) ListUserConversations(ctx context.Context,userID uuid.UUID,) ([]ConversationWithParticipants, error) {

	convs, err := s.messagesRepo.GetUserConversations(ctx, userID)
	if err != nil {
		return nil, err
	}

	if len(convs) == 0 {
		return []ConversationWithParticipants{}, nil
	}

	ids := make([]uuid.UUID, len(convs))
	for i, c := range convs {
		ids[i] = c.ID
	}

	participantsMap, err := s.messagesRepo.GetParticipantsByConversationIDs(ctx, ids)
	if err != nil {
		return nil, err
	}

	unseenMap, err := s.messagesRepo.GetUnreadMessagesCountByUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	result := make([]ConversationWithParticipants, len(convs))

	for i, c := range convs {

		participants := participantsMap[c.ID]

		if len(participants) != 2 {
			return nil, apperr.ErrInternal
		}

		ps := make([]ConvParticipant, len(participants))
		for j, p := range participants {
			ps[j] = toConvParticipant(p)
		}

		result[i] = ConversationWithParticipants{
			Conversation: toConversation(c, unseenMap[c.ID]),
			Participants: ps,
		}
	}

	return result, nil
}

func (s *MessageService) MarkConversationAsRead(ctx context.Context,conversationID uuid.UUID,userID uuid.UUID) error {
	
	return s.messagesRepo.MarkAsReadUpToMessage(
		ctx,
		conversationID,
		userID,
		time.Now(),
	)
}

func toConversation(c entities.Conversation , unseen int) Conversation {
	return Conversation{
		ID:        c.ID,
		BookingID: c.BookingID,
		CreatedAt: c.CreatedAt,
		UnseenCount: unseen,
	}
}

func toConvParticipant(p entities.ConvParticipant) ConvParticipant {
	return ConvParticipant{
		Role:   p.Role,
		UserID: p.UserId,
	}
}

func toMessage(m entities.Message, status bool) Message {
	return Message{
		ID:             m.ID,
		ConversationID: m.ConversationID,
		SenderID:       m.SenderID,
		Content:        m.Content,
		IsRead:         m.IsRead,
		Deleted:         status,
		SentAt:         m.SentAt,
		UpdatedAt:      m.UpdatedAt,
	}
}

func resolveMessageStatus(m entities.Message, userID uuid.UUID) bool {
	switch m.Status {
	case 0:
		return false

	case 1:
		if m.SenderID == userID {
			return true
		}
		return false

	case 2:
		return true

	default:
		return false
	}
}


