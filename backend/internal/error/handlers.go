// internal/apperrors/handler.go
package apperr

import (
    "errors"
    "net/http"
    "github.com/gin-gonic/gin"
)

func Handle(c *gin.Context, err error) {
    switch {
    case errors.Is(err, ErrNotFound):
        c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
    case errors.Is(err, ErrUnauthorized):
        c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
    case errors.Is(err, ErrForbidden):
        c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
    case errors.Is(err, ErrConflict):
        c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
    case errors.Is(err, ErrBadInput):
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
    default:
        c.JSON(http.StatusInternalServerError, gin.H{"error": ErrInternal.Error()})
    }
}