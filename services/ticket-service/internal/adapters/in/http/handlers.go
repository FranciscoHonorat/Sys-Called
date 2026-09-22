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
		Title:       body.Title,
		Description: body.Description,
		AssigneeID:  body.AssigneeID,
		Priority:    body.Priority,
	})
	respondJSON(c, http.StatusCreated, output, err)
}

func (h *Handler) ListTickets(c *gin.Context) {
	output, err := h.useCases.ListTickets.Execute(c.Request.Context())
	respondJSON(c, http.StatusOK, output, err)
}

func (h *Handler) GetTicket(c *gin.Context) {
	output, err := h.useCases.GetTicket.Execute(c.Request.Context(), query.GetTicketInput{
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
		TicketID:   c.Param("id"),
		AssigneeID: body.AssigneeID,
	})
	respondNoContent(c, err)
}

func (h *Handler) AutoAssignTicket(c *gin.Context) {
	output, err := h.useCases.AutoAssignTicket.Execute(c.Request.Context(), command.AutoAssignTicketInput{
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
		TicketID: c.Param("id"),
		Priority: body.Priority,
	})
	respondNoContent(c, err)
}

func (h *Handler) MoveTicketToInProgress(c *gin.Context) {
	err := h.useCases.MoveTicketToInProgress.Execute(c.Request.Context(), command.MoveTicketToInProgressInput{
		TicketID: c.Param("id"),
	})
	respondNoContent(c, err)
}

func (h *Handler) CloseTicket(c *gin.Context) {
	err := h.useCases.CloseTicket.Execute(c.Request.Context(), command.CloseTicketInput{
		TicketID: c.Param("id"),
	})
	respondNoContent(c, err)
}

type addTicketResponseRequest struct {
	AuthorID string `json:"author_id"`
	Content  string `json:"content"`
}

func (h *Handler) AddTicketResponse(c *gin.Context) {
	var body addTicketResponseRequest
	if !bindJSON(c, &body) {
		return
	}

	err := h.useCases.AddTicketResponse.Execute(c.Request.Context(), command.AddTicketResponseInput{
		TicketID: c.Param("id"),
		AuthorID: body.AuthorID,
		Content:  body.Content,
	})
	respondNoContent(c, err)
}
