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

	authenticated := RequireAuthentication(h.useCases.Authenticate)
	r.GET("/responsibles", authenticated, h.ListResponsibles)
	r.GET("/responsibles/workload", authenticated, h.SupportWorkload)
	r.GET("/notifications", authenticated, h.ListNotifications)
	r.POST("/notifications/read", authenticated, h.MarkNotificationsRead)

	tickets := r.Group("/tickets", authenticated)
	tickets.POST("", h.OpenTicket)
	tickets.GET("", h.ListTickets)
	tickets.GET("/:id", h.GetTicket)
	tickets.PUT("/:id", h.EditTicket)
	tickets.POST("/:id/assign", h.AssignTicket)
	tickets.POST("/:id/assign/auto", h.AutoAssignTicket)
	tickets.POST("/:id/priority", h.ChangeTicketPriority)
	tickets.POST("/:id/start", h.MoveTicketToInProgress)
	tickets.POST("/:id/close", h.CloseTicket)
	tickets.POST("/:id/responses", h.AddTicketResponse)

	return r
}
