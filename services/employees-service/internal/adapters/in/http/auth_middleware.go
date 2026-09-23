package httpapi

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/franciscoHonorat/Sys-Called/services/employees-service/internal/application"
	"github.com/franciscoHonorat/Sys-Called/services/employees-service/internal/application/port/out"
)

const callerKey = "caller"

func RequireAuthentication(authenticate *application.AuthenticateUseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		caller, err := authenticate.Execute(bearerToken(c))
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}
		c.Set(callerKey, caller)
		c.Next()
	}
}

func bearerToken(c *gin.Context) string {
	token, ok := strings.CutPrefix(c.GetHeader("Authorization"), "Bearer ")
	if !ok {
		return ""
	}
	return token
}

func CallerFrom(c *gin.Context) out.Caller {
	caller, _ := c.Get(callerKey)
	return caller.(out.Caller)
}
