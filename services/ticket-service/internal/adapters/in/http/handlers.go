package httpapi

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/application/command"
	"github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/application/query"
)

type openTicketRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	AssigneeID  string `json:"assignee_id"`
	Priority    string `json:"priority"`
}

func (h *Handler) OpenTicket(c *gin.Context) {
	var body openTicketRequest
	if !bindJSON(c, &body) {
		return
	}

	output, err := h.useCases.OpenTicket.Execute(c.Request.Context(), command.OpenTicketInput{
		Actor:       ActorFrom(c),
		Title:       body.Title,
		Description: body.Description,
		AssigneeID:  body.AssigneeID,
		Priority:    body.Priority,
	})
	respondJSON(c, http.StatusCreated, output, err)
}

func (h *Handler) ListTickets(c *gin.Context) {
	output, err := h.useCases.ListTickets.Execute(c.Request.Context(), ActorFrom(c))
	respondJSON(c, http.StatusOK, output, err)
}

func (h *Handler) GetTicket(c *gin.Context) {
	output, err := h.useCases.GetTicket.Execute(c.Request.Context(), query.GetTicketInput{
		Viewer:   ActorFrom(c),
		TicketID: c.Param("id"),
	})
	respondJSON(c, http.StatusOK, output, err)
}

type editTicketRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

func (h *Handler) EditTicket(c *gin.Context) {
	var body editTicketRequest
	if !bindJSON(c, &body) {
		return
	}

	err := h.useCases.EditTicket.Execute(c.Request.Context(), command.EditTicketInput{
		Actor:       ActorFrom(c),
		TicketID:    c.Param("id"),
		Title:       body.Title,
		Description: body.Description,
	})
	respondNoContent(c, err)
}

type assignTicketRequest struct {
	AssigneeID string `json:"assignee_id"`
}

func (h *Handler) AssignTicket(c *gin.Context) {
	var body assignTicketRequest
	if !bindJSON(c, &body) {
		return
	}

	err := h.useCases.AssignTicket.Execute(c.Request.Context(), command.AssignTicketInput{
		Actor:      ActorFrom(c),
		TicketID:   c.Param("id"),
		AssigneeID: body.AssigneeID,
	})
	respondNoContent(c, err)
}

func (h *Handler) AutoAssignTicket(c *gin.Context) {
	output, err := h.useCases.AutoAssignTicket.Execute(c.Request.Context(), command.AutoAssignTicketInput{
		Actor:    ActorFrom(c),
		TicketID: c.Param("id"),
	})
	respondJSON(c, http.StatusOK, output, err)
}

type changeTicketPriorityRequest struct {
	Priority string `json:"priority"`
}

func (h *Handler) ChangeTicketPriority(c *gin.Context) {
	var body changeTicketPriorityRequest
	if !bindJSON(c, &body) {
		return
	}

	err := h.useCases.ChangeTicketPriority.Execute(c.Request.Context(), command.ChangeTicketPriorityInput{
		Actor:    ActorFrom(c),
		TicketID: c.Param("id"),
		Priority: body.Priority,
	})
	respondNoContent(c, err)
}

func (h *Handler) MoveTicketToInProgress(c *gin.Context) {
	err := h.useCases.MoveTicketToInProgress.Execute(c.Request.Context(), command.MoveTicketToInProgressInput{
		Actor:    ActorFrom(c),
		TicketID: c.Param("id"),
	})
	respondNoContent(c, err)
}

type closeTicketRequest struct {
	Resolution string `json:"resolution"`
}

func (h *Handler) CloseTicket(c *gin.Context) {
	var body closeTicketRequest
	if !bindJSON(c, &body) {
		return
	}

	err := h.useCases.CloseTicket.Execute(c.Request.Context(), command.CloseTicketInput{
		Actor:      ActorFrom(c),
		TicketID:   c.Param("id"),
		Resolution: body.Resolution,
	})
	respondNoContent(c, err)
}

type addTicketResponseRequest struct {
	Content string `json:"content"`
}

func (h *Handler) AddTicketResponse(c *gin.Context) {
	var body addTicketResponseRequest
	if !bindJSON(c, &body) {
		return
	}

	err := h.useCases.AddTicketResponse.Execute(c.Request.Context(), command.AddTicketResponseInput{
		Actor:    ActorFrom(c),
		TicketID: c.Param("id"),
		Content:  body.Content,
	})
	respondNoContent(c, err)
}

func (h *Handler) ListResponsibles(c *gin.Context) {
	output, err := h.useCases.ListResponsibles.Execute(c.Request.Context())
	respondJSON(c, http.StatusOK, output, err)
}

func (h *Handler) ListNotifications(c *gin.Context) {
	output, err := h.useCases.ListNotifications.Execute(c.Request.Context(), ActorFrom(c))
	respondJSON(c, http.StatusOK, output, err)
}

func (h *Handler) MarkNotificationsRead(c *gin.Context) {
	respondNoContent(c, h.useCases.MarkNotificationsRead.Execute(c.Request.Context(), ActorFrom(c)))
}

func (h *Handler) SupportWorkload(c *gin.Context) {
	output, err := h.useCases.SupportWorkload.Execute(c.Request.Context(), ActorFrom(c))
	respondJSON(c, http.StatusOK, output, err)
}
