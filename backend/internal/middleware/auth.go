package middleware

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sdi2200246/synaxis/internal/error"
	"github.com/sdi2200246/synaxis/internal/services"
)

type AuthHandler struct{
	authService *services.AuthService
}

func NewAuthHandler(s *services.AuthService)*AuthHandler{
	return &AuthHandler{authService: s}
}

func (h *AuthHandler)Login(c *gin.Context){

	var credentials services.UserCridentials
	 if err := c.ShouldBindJSON(&credentials); err != nil {
        c.JSON(400, gin.H{"error": "invalid cridentials given"})
        return
    }

	token , err := h.authService.Login(c.Request.Context() , credentials)

    if err != nil {
        h.handleError(c, err)
        return
    }
	c.JSON(200, gin.H{"jwt_token": token})
}

func (h *AuthHandler) AuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        header := c.GetHeader("Authorization")
        if header == "" || !strings.HasPrefix(header, "Bearer ") {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing token"})
            return
        }

        tokenString := strings.TrimPrefix(header, "Bearer ")

        claims, err := h.authService.ValidateToken(tokenString)
        if err != nil {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
            return
        }
        c.Set("userID", claims.UserID)
        c.Set("role", claims.Role)
        c.Next()
    }
}

func (h *AuthHandler)AdminOnly() gin.HandlerFunc {
    return func(c *gin.Context) {
        if c.GetString("role") != "admin" {
              role, exists := c.Get("role")
              fmt.Printf("role: %v, type: %T, exists: %v\n", role, role, exists)
              c.AbortWithStatusJSON(403, gin.H{"error": "forbidden"})
            return
        }
        c.Next()
    }
}

func (h *AuthHandler) OptionalAuth() gin.HandlerFunc {
    return func(c *gin.Context) {
        header := c.GetHeader("Authorization")
        if header == "" {
            slog.Warn("there is not header")
            c.Next()
            return
        }
        tokenString := strings.TrimPrefix(header, "Bearer ")
        claims, err := h.authService.ValidateToken(tokenString)
        if err != nil {
            c.Next() // invalid token on a semi-public endpoint → treat as guest
            return
        }
        c.Set("userID", claims.UserID)
        c.Set("role", claims.Role)
        c.Next()
    }
}

func (h *AuthHandler) handleError(c *gin.Context, err error) {
    switch {
    case errors.Is(err, apperr.ErrUnauthorized):
        c.JSON(401, gin.H{"error": "invalid username or password"})
    case errors.Is(err, apperr.ErrPendingApproval):
        c.JSON(403, gin.H{"error": "account pending admin approval"})
    case errors.Is(err, apperr.ErrRejected):
        c.JSON(403, gin.H{"error": "account has been rejected"})
    default:
        apperr.Handle(c, err)
    }
}