package controllers

import (
	"log/slog"
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	apperr "github.com/sdi2200246/synaxis/internal/error"
	"github.com/sdi2200246/synaxis/internal/services"
)

type CreateConversationRequest struct {
	BookingID   uuid.UUID `json:"booking_id" binding:"required"`
	OrganizerID uuid.UUID `json:"organizer_id" binding:"required"`
	AttendeeID  uuid.UUID `json:"attendee_id" binding:"required"`
}

type CreateMessageRequest struct {
	Content string `json:"content" binding:"required"`
}

type UpdateMessageRequest struct{
	Content *string `json:"content,omitempty"`	
	Delete *int 	`json:"delete,omitempty"`
}

type MessagesHandler struct {
	messagesService *services.MessageService
}

func NewMessagesHandler(messagesService *services.MessageService) *MessagesHandler {
	return &MessagesHandler{messagesService: messagesService}
}

func (h *MessagesHandler) CreateConversation(c *gin.Context) {
	
	var input CreateConversationRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		slog.Error("Invalid input", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input", "details": err.Error()})
		return
	}

	convID, err := h.messagesService.CreateConversation(c.Request.Context(), services.CreateConversationInput{
		BookingID:   input.BookingID,
		OrganizerID: input.OrganizerID,
		AttendeeID:  input.AttendeeID,
	})
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, CreateConversationResponse{
		ConversationID: convID,
	})
}

func (h *MessagesHandler) CreateMessage(c *gin.Context) {
	conversationID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid conversation id"})
		return
	}

	userID, ok := getUserIDFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var input CreateMessageRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		slog.Error("Invalid input", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input", "details": err.Error()})
		return
	}

	err = h.messagesService.SendMessage(c.Request.Context(), services.CreateMessageInput{
		ConversationID: conversationID,
		SenderID:       userID,
		Content:        input.Content,
	})
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.Status(http.StatusCreated)
}

func (h *MessagesHandler) UpdateMessage(c *gin.Context) {
	messageID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid message id"})
		return
	}

	userID, ok := getUserIDFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var input UpdateMessageRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input", "details": err.Error()})
		return
	}

	if input.Content == nil && input.Delete == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "nothing to update"})
		return
	}
	
	if err := h.messagesService.UpdateMessage(c.Request.Context(), messageID, userID, services.UpdateMessageInput{
		Content: input.Content,
		Delete: input.Delete,
	}); err != nil {
		h.handleError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *MessagesHandler) ListUserConversations(c *gin.Context) {
	userID, ok := getUserIDFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	convs, err := h.messagesService.ListUserConversations(c.Request.Context(), userID)
	if err != nil {
		h.handleError(c, err)
		return
	}

	result := make([]ConversationWithParticipantsResponse, 0, len(convs))

	for _, conv := range convs {

		ps := make([]ConvParticipantResponse, 0, len(conv.Participants))
		for _, p := range conv.Participants {
			ps = append(ps, ConvParticipantResponse{
				Role:   p.Role,
				UserID: p.UserID,
			})
		}

		result = append(result, ConversationWithParticipantsResponse{
			Conversation: ConversationResponse{
				ID:          conv.Conversation.ID,
				BookingID:   conv.Conversation.BookingID,
				CreatedAt:   conv.Conversation.CreatedAt,
				UnseenCount: conv.Conversation.UnseenCount,
			},
			Participants: ps,
			EventTitle: conv.EventName,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"conversations": result,
	})
}

func (h *MessagesHandler) GetConversationMessages(c *gin.Context) {
	conversationID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid conversation id"})
		return
	}

	userID, ok := getUserIDFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	messages, err := h.messagesService.GetConversationMessages(
		c.Request.Context(),
		conversationID,
		userID,
	)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"messages": ToMessageListResponse(messages),
	})
}

func (h *MessagesHandler) MarkConversationAsRead(c *gin.Context) {
	conversationID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid conversation id"})
		return
	}

	userID, ok := getUserIDFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	err = h.messagesService.MarkConversationAsRead(
		c.Request.Context(),
		conversationID,
		userID,
	)
	if err != nil {
		h.handleError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}


func (h *MessagesHandler) handleError(c *gin.Context, err error) {
	apperr.Handle(c, err)
}

func getUserIDFromContext(c *gin.Context) (uuid.UUID, bool) {
	val, exists := c.Get("userID")
	if !exists {
		return uuid.Nil, false
	}

	userID, ok := val.(uuid.UUID)
	if !ok {
		return uuid.Nil, false
	}

	return userID, true
}