package controllers

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	apperr "github.com/sdi2200246/synaxis/internal/error"
	"github.com/sdi2200246/synaxis/internal/services"
)

type VisitsHandler struct {
	baseHandler *BaseHandler
	visitsService *services.VisitService
}

func NewVisitsHandler(vs *services.VisitService , bh *BaseHandler) *VisitsHandler {
    return &VisitsHandler{visitsService: vs , baseHandler: bh}
}

func (h *VisitsHandler) Record(c *gin.Context){

	userID , err := h.baseHandler.getUserIDFromContext(c)
	if err != nil{
		apperr.Handle(c , err)
		return
	}

	eventID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid event id"})
		return
	}

	if err = h.visitsService.RecordVisit(c.Request.Context() , userID , eventID) ; err != nil{
		apperr.Handle(c , err)
		return
	}
	c.Status(http.StatusCreated)
}




