package controllers

import(
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	apperr "github.com/sdi2200246/synaxis/internal/error"
)
type BaseHandler struct{}

func (b *BaseHandler) getUserIDFromContext(c *gin.Context) (uuid.UUID, error) {
    val, exists := c.Get("userID")
    if !exists {
        return uuid.Nil, fmt.Errorf("user is unauthorized: %w", apperr.ErrUnauthorized)
    }
    userID, ok := val.(uuid.UUID)
    if !ok {
        return uuid.Nil, apperr.ErrInternal
    }
    return userID, nil
}
func (h *BaseHandler) CallerIDExists(c *gin.Context) (*uuid.UUID, bool) {
    val, exists := c.Get("userID")
    if !exists {
        return &uuid.Nil, false
    }
    id, ok := val.(uuid.UUID)
    return &id, ok
}