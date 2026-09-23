package httpapi

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/application/auth"
	"github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/domain/actor"
)

const actorKey = "actor"

func RequireAuthentication(authenticate *auth.AuthenticateUseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		a, err := authenticate.Execute(c.Request.Context(), bearerToken(c))
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		c.Set(actorKey, a)
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

func ActorFrom(c *gin.Context) actor.Actor {
	a, _ := c.Get(actorKey)
	return a.(actor.Actor)
}
