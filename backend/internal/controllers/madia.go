package controllers

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	apperr "github.com/sdi2200246/synaxis/internal/error"
	"github.com/sdi2200246/synaxis/internal/services"
)

const (
	mediaUploadDir = "uploads/events"
	mediaURLPrefix = "/media/events"
)

type MediaResponse struct {
	ID         uuid.UUID `json:"id"`
	EventID    uuid.UUID `json:"event_id"`
	Filename   string    `json:"filename"`
	URL        string    `json:"url"`
	UploadedAt time.Time `json:"uploaded_at"`
}

type MediaHandler struct {
	baseHandler  *BaseHandler
	mediaService *services.MediaService
}

func NewMediaHandler(ms *services.MediaService, bh *BaseHandler) *MediaHandler {
	return &MediaHandler{baseHandler: bh, mediaService: ms}
}


func (h*MediaHandler) Upload(c *gin.Context){
	callerID, err := h.baseHandler.getUserIDFromContext(c)
	if err != nil {
		apperr.Handle(c, err)
		return
	}

	eventID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid event id"})
		return
	}

	file , err := c.FormFile("photo")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "photo file is required"})
		return
	}

	ext := filepath.Ext(file.Filename)

	media, err := h.mediaService.Upload(c.Request.Context(), callerID, eventID,file.Size,ext)
	if err != nil {
		apperr.Handle(c, err)
		return
	}

	dst := filepath.Join(mediaUploadDir, eventID.String(), media.Filename)
	if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
		apperr.Handle(c, apperr.ErrInternal)
		return
	}
	if err := c.SaveUploadedFile(file, dst); err != nil {
		apperr.Handle(c, apperr.ErrInternal)
		return
	}

	c.JSON(http.StatusCreated, toMediaResponse(media, eventID))

}


func (h *MediaHandler) Delete(c *gin.Context) {
	callerID, err := h.baseHandler.getUserIDFromContext(c)
	if err != nil {
		apperr.Handle(c, err)
		return
	}

	eventID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid event id"})
		return
	}

	mediaID, err := uuid.Parse(c.Param("media_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid media id"})
		return
	}

	media, err := h.mediaService.Delete(c.Request.Context(), callerID, eventID, mediaID)
	if err != nil {
		apperr.Handle(c, err)
		return
	}

	dst := filepath.Join(mediaUploadDir, eventID.String(), media.Filename)
	_ = os.Remove(dst)

	c.Status(http.StatusNoContent)
}


func toMediaResponse(m services.Media, eventID uuid.UUID) MediaResponse {
	return MediaResponse{
		ID:         m.ID,
		EventID:    m.EventID,
		Filename:   m.Filename,
		URL:        fmt.Sprintf("%s/%s/%s", mediaURLPrefix, eventID.String(), m.Filename),
		UploadedAt: m.UploadedAt,
	}
}