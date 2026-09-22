package httpapi

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/franciscoHonorat/Sys-Called/services/employees-service/internal/application"
)

type Handler struct {
	listEmployees *application.ListEmployeesUseCase
}

func NewHandler(listEmployees *application.ListEmployeesUseCase) *Handler {
	return &Handler{listEmployees: listEmployees}
}

func (h *Handler) ListEmployees(c *gin.Context) {
	output, err := h.listEmployees.Execute(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	c.JSON(http.StatusOK, output)
}
