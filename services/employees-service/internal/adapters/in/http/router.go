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

	r.GET("/employees", h.ListEmployees)

	return r
}
