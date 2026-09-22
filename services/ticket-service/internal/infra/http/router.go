package httpapi

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func NewRouter(h *Handler) *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery())

	r.GET("/healthz", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	r.POST("/tickets", h.OpenTicket)
	r.GET("/tickets/:id", h.GetTicket)
	r.POST("/tickets/:id/assign", h.AssignTicket)
	r.POST("/tickets/:id/priority", h.ChangeTicketPriority)
	r.POST("/tickets/:id/start", h.MoveTicketToInProgress)
	r.POST("/tickets/:id/close", h.CloseTicket)
	r.POST("/tickets/:id/responses", h.AddTicketResponse)

	return r
}
