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

	r.GET("/employees", authenticated, h.ListEmployees)
	r.POST("/employees/:id/approve", authenticated, h.ApproveEmployee)
	r.POST("/employees/:id/temporary-password", authenticated, h.IssueTemporaryPassword)
	r.POST("/auth/signup", h.SignUp)
	r.POST("/auth/password-reset-requests", h.RequestPasswordReset)
	r.POST("/auth/change-password", authenticated, h.ChangePassword)
	r.POST("/auth/login", h.Login)
	r.POST("/auth/refresh", h.RefreshSession)
	r.POST("/auth/logout", h.Logout)
	r.GET("/.well-known/jwks.json", h.JWKS)

	return r
}
