package httpapi

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/franciscoHonorat/Sys-Called/services/employees-service/internal/application"
)

type signUpRequest struct {
	Name     string `json:"name"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func (h *Handler) SignUp(c *gin.Context) {
	var body signUpRequest
	if !bindJSON(c, &body) {
		return
	}
	err := h.useCases.SignUp.Execute(c.Request.Context(), application.SignUpInput{
		Name:     body.Name,
		Username: body.Username,
		Password: body.Password,
	})
	respondStatus(c, http.StatusCreated, err)
}

type passwordResetRequest struct {
	Username string `json:"username"`
}

func (h *Handler) RequestPasswordReset(c *gin.Context) {
	var body passwordResetRequest
	if !bindJSON(c, &body) {
		return
	}
	err := h.useCases.RequestPasswordReset.Execute(c.Request.Context(), body.Username)
	respondStatus(c, http.StatusAccepted, err)
}

func (h *Handler) ApproveEmployee(c *gin.Context) {
	err := h.useCases.ApproveEmployee.Execute(c.Request.Context(), CallerFrom(c), c.Param("id"))
	respondStatus(c, http.StatusNoContent, err)
}

func (h *Handler) IssueTemporaryPassword(c *gin.Context) {
	password, err := h.useCases.IssueTemporaryPassword.Execute(c.Request.Context(), CallerFrom(c), c.Param("id"))
	respondJSON(c, gin.H{"temporary_password": password}, err)
}

type changePasswordRequest struct {
	CurrentPassword string `json:"current_password"`
	NewPassword     string `json:"new_password"`
}

func (h *Handler) ChangePassword(c *gin.Context) {
	var body changePasswordRequest
	if !bindJSON(c, &body) {
		return
	}
	err := h.useCases.ChangePassword.Execute(c.Request.Context(), CallerFrom(c), application.ChangePasswordInput{
		Current: body.CurrentPassword,
		New:     body.NewPassword,
	})
	respondStatus(c, http.StatusNoContent, err)
}
