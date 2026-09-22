package httpapi

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/application"
)

type openTicketRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	AssigneeID  string `json:"assignee_id"`
	Priority    string `json:"priority"`
}

func (h *Handler) OpenTicket(c *gin.Context) {
	var body openTicketRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	output, err := h.openTicket.Execute(c.Request.Context(), application.OpenTicketInput{
		Title:       body.Title,
		Description: body.Description,
		AssigneeID:  body.AssigneeID,
		Priority:    body.Priority,
	})
	if err != nil {
		respondError(c, err)
		return
	}

	c.JSON(http.StatusCreated, output)
}

func (h *Handler) ListTickets(c *gin.Context) {
	output, err := h.listTickets.Execute(c.Request.Context())
	if err != nil {
		respondError(c, err)
		return
	}

	c.JSON(http.StatusOK, output)
}

func (h *Handler) GetTicket(c *gin.Context) {
	output, err := h.getTicket.Execute(c.Request.Context(), application.GetTicketInput{
		TicketID: c.Param("id"),
	})
	if err != nil {
		respondError(c, err)
		return
	}

	c.JSON(http.StatusOK, output)
}

type editTicketRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

func (h *Handler) EditTicket(c *gin.Context) {
	var body editTicketRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	err := h.editTicket.Execute(c.Request.Context(), application.EditTicketInput{
		TicketID:    c.Param("id"),
		Title:       body.Title,
		Description: body.Description,
	})
	if err != nil {
		respondError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

type assignTicketRequest struct {
	AssigneeID string `json:"assignee_id"`
}

func (h *Handler) AssignTicket(c *gin.Context) {
	var body assignTicketRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	err := h.assignTicket.Execute(c.Request.Context(), application.AssignTicketInput{
		TicketID:   c.Param("id"),
		AssigneeID: body.AssigneeID,
	})
	if err != nil {
		respondError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *Handler) AutoAssignTicket(c *gin.Context) {
	output, err := h.autoAssignTicket.Execute(c.Request.Context(), application.AutoAssignTicketInput{
		TicketID: c.Param("id"),
	})
	if err != nil {
		respondError(c, err)
		return
	}

	c.JSON(http.StatusOK, output)
}

type changeTicketPriorityRequest struct {
	Priority string `json:"priority"`
}

func (h *Handler) ChangeTicketPriority(c *gin.Context) {
	var body changeTicketPriorityRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	err := h.changeTicketPriority.Execute(c.Request.Context(), application.ChangeTicketPriorityInput{
		TicketID: c.Param("id"),
		Priority: body.Priority,
	})
	if err != nil {
		respondError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *Handler) MoveTicketToInProgress(c *gin.Context) {
	err := h.moveTicketToInProgress.Execute(c.Request.Context(), application.MoveTicketToInProgressInput{
		TicketID: c.Param("id"),
	})
	if err != nil {
		respondError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *Handler) CloseTicket(c *gin.Context) {
	err := h.closeTicket.Execute(c.Request.Context(), application.CloseTicketInput{
		TicketID: c.Param("id"),
	})
	if err != nil {
		respondError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

type addTicketResponseRequest struct {
	AuthorID string `json:"author_id"`
	Content  string `json:"content"`
}

func (h *Handler) AddTicketResponse(c *gin.Context) {
	var body addTicketResponseRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	err := h.addTicketResponse.Execute(c.Request.Context(), application.AddTicketResponseInput{
		TicketID: c.Param("id"),
		AuthorID: body.AuthorID,
		Content:  body.Content,
	})
	if err != nil {
		respondError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}
