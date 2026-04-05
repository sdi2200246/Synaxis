package controllers

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/sdi2200246/synaxis/internal/error"
	"github.com/sdi2200246/synaxis/internal/services"
	"github.com/sdi2200246/synaxis/internal/types"
)
type UserHandler struct {
    userService *services.UserService
}

func NewUserHandler(userService *services.UserService) *UserHandler {
    return &UserHandler{userService: userService}
}

func (h *UserHandler)Register(c *gin.Context) {
    var input services.CandidateUser
    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(400, gin.H{"error": "invalid input"})
        return
    }

    err := h.userService.RegisterUser(c.Request.Context(), input)
    if err != nil {
        h.handleError(c, err)
        return
    }

    c.JSON(201, gin.H{"message": "registration pending admin approval"})
}

func (h *UserHandler)GetUsers(c *gin.Context){
    var filter types.UserFilter
    if err := c.ShouldBindQuery(&filter); err != nil {
        c.JSON(400, gin.H{"error": "invalid query params"})
        return
    }
    users , err := h.userService.GetUsers(c.Request.Context() , filter)
    if err !=nil{
        h.handleError(c,err)
        return
    }
    c.JSON(200, gin.H{"count":len(users), "users":users})
}


func (h *UserHandler) handleError(c *gin.Context, err error){
    switch {
    case errors.Is(err, apperr.ErrUsernameConflict):
        c.JSON(409, gin.H{"error": err.Error(), "field": "username"})
    case errors.Is(err, apperr.ErrEmailConflict):
        c.JSON(409, gin.H{"error": err.Error(), "field": "email"})
    case errors.Is(err, apperr.ErrBadInput):
        c.JSON(400, gin.H{"error": err.Error()})
    default:
        apperr.Handle(c, err)
    }
}