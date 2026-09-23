package httpapi

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	domainErr "github.com/franciscoHonorat/Sys-Called/services/employees-service/internal/domain/domain-errors"
)

func bindJSON(c *gin.Context, body any) bool {
	if err := c.ShouldBindJSON(body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return false
	}
	return true
}

func respondJSON(c *gin.Context, body any, err error) {
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, body)
}

func respondStatus(c *gin.Context, status int, err error) {
	if err != nil {
		respondError(c, err)
		return
	}
	c.Status(status)
}

func respondError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, domainErr.ErrInvalidCredentials), errors.Is(err, domainErr.ErrInvalidRefreshToken):
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
	case errors.Is(err, domainErr.ErrForbidden), errors.Is(err, domainErr.ErrPendingApproval):
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
	case errors.Is(err, domainErr.ErrEmployeeNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	case errors.Is(err, domainErr.ErrUsernameTaken):
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
	case errors.Is(err, domainErr.ErrWeakPassword), errors.Is(err, domainErr.ErrPasswordTooLong), errors.Is(err, domainErr.ErrInvalidEmployeeName), errors.Is(err, domainErr.ErrInvalidUsername):
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
	}
}
