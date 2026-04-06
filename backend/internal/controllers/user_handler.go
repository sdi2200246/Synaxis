package controllers

import (
	"errors"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sdi2200246/synaxis/internal/error"
	"github.com/sdi2200246/synaxis/internal/services"
)

type UserFilter struct {
    Country   string `form:"country"`
    Status    string `form:"status"`
    CreatedAt string `form:"created_at"`
}

type RegisterUserInput struct{
    Username  string `json:"username"   binding:"required"`
    Password  string `json:"password"   binding:"required"`
    FirstName string `json:"first_name" binding:"required"`
    LastName  string `json:"last_name"  binding:"required"`
    Email     string `json:"email"      binding:"required,email"`
    Phone     string `json:"phone"      binding:"required"`
    Address   string `json:"address"    binding:"required"`
    City      string `json:"city"       binding:"required"`
    Country   string `json:"country"    binding:"required"`
    TaxID     string `json:"tax_id"     binding:"required"`
}

type UserHandler struct {
    userService *services.UserService
}

func NewUserHandler(userService *services.UserService) *UserHandler {
    return &UserHandler{userService: userService}
}

func (h *UserHandler)Register(c *gin.Context) {
    var input RegisterUserInput
    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(400, gin.H{"error": "invalid input"})
        return
    }

    err := h.userService.RegisterUser(c.Request.Context(), services.CandidateUser{
        Username:  input.Username,
        Password:  input.Password,
        FirstName: input.FirstName,
        LastName:  input.LastName,
        Email:     input.Email,
        Phone:     input.Phone,
        Address:   input.Address,
        City:      input.City,
        Country:   input.Country,
        TaxID:     input.TaxID,
    })
    if err != nil {
        h.handleError(c, err)
        return
    }

    c.JSON(201, gin.H{"message": "registration pending admin approval"})
}

func (h *UserHandler)GetUsers(c *gin.Context) {
    var filter UserFilter
    if err := c.ShouldBindQuery(&filter); err != nil {
        c.JSON(400, gin.H{"error": "invalid query params"})
        return
    }
    svcFilter := services.UserFilter{}

    if filter.Country != "" {
        svcFilter.Country = &filter.Country
    }
    if filter.Status != "" {
        svcFilter.Status = &filter.Status
    }
    if filter.CreatedAt != "" {
        t, err := time.Parse(time.RFC3339, filter.CreatedAt)
        if err != nil {
            c.JSON(400, gin.H{"error": "invalid created_at format"})
            return
        }
        svcFilter.CreatedAt = &t
    }

    users, err := h.userService.GetUsers(c.Request.Context(), svcFilter)
    if err != nil {
        h.handleError(c, err)
        return
    }
   

    plain := make([]AdminUserResponse, len(users))
        for i, u := range users {
            plain[i] = AdminUserResponse{
                ID:        u.ID,
                Username:  u.Username,
                FirstName: u.FirstName,
                LastName:  u.LastName,
                Email:     u.Email,
                Phone:     u.Phone,
                Address:   u.Address,
                City:      u.City,
                Country:   u.Country,
                TaxID:     u.TaxID,
                Status:    u.Status,
                CreatedAt: u.CreatedAt,
                UpdatedAt: u.UpdatedAt,
            }
    }    
    c.JSON(200, gin.H{"count": len(users), "users": plain})
}

func (h *UserHandler) ApproveUser(c *gin.Context) {
    id, err := uuid.Parse(c.Param("id"))
    if err != nil {
        c.JSON(400, gin.H{"error": "invalid user id"})
        return
    }

    err = h.userService.ApproveUser(c.Request.Context(), id)
    if err != nil {
        h.handleError(c, err)
        return
    }
    c.JSON(200, gin.H{"message": "user approved"})
}

func (h *UserHandler) RejectUser(c *gin.Context) {
    id, err := uuid.Parse(c.Param("id"))
    if err != nil {
        c.JSON(400, gin.H{"error": "invalid user id"})
        return
    }

    err = h.userService.RejectUser(c.Request.Context(), id)
    if err != nil {
        h.handleError(c, err)
        return
    }
    c.JSON(200, gin.H{"message": "user rejected"})
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