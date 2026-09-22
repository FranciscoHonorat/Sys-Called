package httpapi

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	domainErr "github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/domain/domain-errors"
)

func respondError(c *gin.Context, err error) {
	c.JSON(statusForError(err), gin.H{"error": err.Error()})
}

func statusForError(err error) int {
	switch {
	case errors.Is(err, domainErr.ErrEventStreamNotFound):
		return http.StatusNotFound
	case errors.Is(err, domainErr.ErrConcurrencyConflict),
		errors.Is(err, domainErr.ErrTicketAlreadyClosed),
		errors.Is(err, domainErr.ErrInvalidStatusTransition),
		errors.Is(err, domainErr.ErrResponseTicketMismatch):
		return http.StatusConflict
	case errors.Is(err, domainErr.ErrInvalidTitle),
		errors.Is(err, domainErr.ErrInvalidDescription),
		errors.Is(err, domainErr.ErrInvalidStatus),
		errors.Is(err, domainErr.ErrInvalidPriority),
		errors.Is(err, domainErr.ErrInvalidAssignee),
		errors.Is(err, domainErr.ErrInvalidID),
		errors.Is(err, domainErr.ErrInvalidUUID),
		errors.Is(err, domainErr.ErrInvalidContent),
		errors.Is(err, domainErr.ErrInvalidAuthorID),
		errors.Is(err, domainErr.ErrInvalidTicketID),
		errors.Is(err, domainErr.ErrInvalidResponse):
		return http.StatusBadRequest
	case errors.Is(err, domainErr.ErrNoResponsiblesAvailable):
		return http.StatusUnprocessableEntity
	default:
		return http.StatusInternalServerError
	}
}
