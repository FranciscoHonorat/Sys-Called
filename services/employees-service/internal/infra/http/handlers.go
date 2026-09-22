package httpapi

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) ListEmployees(c *gin.Context) {
	output, err := h.listEmployees.Execute(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	c.JSON(http.StatusOK, output)
}
