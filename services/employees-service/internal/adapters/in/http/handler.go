package httpapi

import (
	"github.com/gin-gonic/gin"

	"github.com/franciscoHonorat/Sys-Called/services/employees-service/internal/application"
)

type UseCases struct {
	Authenticate           *application.AuthenticateUseCase
	ListEmployees          *application.ListEmployeesUseCase
	Login                  *application.LoginUseCase
	RefreshSession         *application.RefreshSessionUseCase
	Logout                 *application.LogoutUseCase
	PublicKeys             *application.GetPublicKeysUseCase
	SignUp                 *application.SignUpUseCase
	RequestPasswordReset   *application.RequestPasswordResetUseCase
	ApproveEmployee        *application.ApproveEmployeeUseCase
	IssueTemporaryPassword *application.IssueTemporaryPasswordUseCase
	ChangePassword         *application.ChangePasswordUseCase
}

type Handler struct {
	useCases UseCases
}

func NewHandler(useCases UseCases) *Handler {
	return &Handler{useCases: useCases}
}

func (h *Handler) ListEmployees(c *gin.Context) {
	output, err := h.useCases.ListEmployees.Execute(c.Request.Context(), CallerFrom(c))
	respondJSON(c, output, err)
}
